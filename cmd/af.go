package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/kunaldawn/mail.af/pkg/db"
	"github.com/kunaldawn/mail.af/pkg/db/models"
	"github.com/kunaldawn/mail.af/pkg/runner"
	"github.com/kunaldawn/mail.af/pkg/utils"
	"github.com/kunaldawn/mail.af/web"
	"github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type MailAF struct {
	runningJobs sync.Map
}

func NewMailAF() *MailAF {
	af := &MailAF{}
	af.init()

	return af
}

func (af *MailAF) init() {
	go func() {
		for true {
			log.Println("POLLING PENDING JOBS...")
			pendingJobs := make([]*models.Job, 0)
			senders := make([]*models.Sender, 0)

			db.DB().Jobs().Find(bson.M{"done": false}).All(&pendingJobs)
			db.DB().Senders().Find(bson.M{}).All(&senders)

			for _, pendingJob := range pendingJobs {
				if pendingJob.StartTime <= time.Now().UnixNano() {
					if _, ok := af.runningJobs.Load(pendingJob.Name); !ok {
						log.Println("QUEUE JOB :", pendingJob.ID, pendingJob.Name)

						af.runningJobs.Store(pendingJob.Name, true)
						runner.NewRunner(pendingJob, senders).Run(func() {
							af.runningJobs.Delete(pendingJob.Name)
						})
					}
				}
			}

			time.Sleep(time.Second * 30)
		}
	}()
}

func (af *MailAF) Start() {
	server := gin.Default()
	dashboard := server.Group("/", gin.Logger(), gin.Recovery())
	dashboard.Any("/*path", af.route)

	err := http.ListenAndServe(":9988", server)
	if err != nil {
		panic(err)
	}
}

func (af *MailAF) route(context *gin.Context) {
	path := context.Param("path")
	if strings.HasPrefix(path, "/api/v1") {
		af.apiV1(context)
	} else if strings.HasPrefix(path, "/auth/login") {
		af.auth(context)
	} else {
		web.GetWeb().Handle(context)
	}
}

func (af *MailAF) auth(context *gin.Context) {
	context.GetHeader("Authorization")
	auth := struct {
		Signature string `json:"signature"`
	}{}
	err := context.BindJSON(&auth)

	if err == nil && len(auth.Signature) > 0 {
		if auth.Signature == utils.GetConfig().GetString("signature") {
			origin := context.Request.Header.Get("origin")
			context.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			context.JSON(http.StatusOK, af.getBearerToken())
			return
		}
	}

	context.JSON(http.StatusUnauthorized, "authentication failed")
}

func (af *MailAF) getBearerToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": uuid.NewV4().String(),
	})
	tokenString, _ := token.SignedString([]byte(utils.GetConfig().GetString("jwt_secret")))
	return tokenString
}

func (af *MailAF) apiV1(context *gin.Context) {
	data := struct {
		Command string          `json:"command"`
		Data    json.RawMessage `json:"data"`
	}{}

	if err := context.BindJSON(&data); err == nil {
		switch data.Command {
		case "get_pending_jobs":
			runningJobs := make([]*models.Job, 0)
			db.DB().Jobs().Find(bson.M{"done": false}).All(&runningJobs)
			jobs := make([]interface{}, 0)
			for _, job := range runningJobs {
				jobs = append(jobs, map[string]interface{}{
					"id":         job.ID,
					"name":       job.Name,
					"start_time": job.StartTime,
					"running":    job.Running,
					"receivers":  len(job.Receivers),
				})
			}
			respData, _ := json.Marshal(jobs)
			context.JSON(http.StatusOK, json.RawMessage(respData))
		case "get_completed_jobs":
			pendingJobs := make([]*models.Job, 0)
			db.DB().Jobs().Find(bson.M{"done": true}).All(&pendingJobs)
			jobs := make([]interface{}, 0)
			for _, job := range pendingJobs {
				jobs = append(jobs, map[string]interface{}{
					"id":         job.ID,
					"name":       job.Name,
					"start_time": job.StartTime,
					"end_time":   job.EndTime,
					"receivers":  len(job.Receivers),
					"success":    job.SendSuccess,
					"failure":    job.SendFailed,
				})
			}
			respData, _ := json.Marshal(jobs)
			context.JSON(http.StatusOK, json.RawMessage(respData))
		case "create_job":
			jobData := struct {
				Name    string   `json:"name"`
				Groups  []string `json:"groups"`
				Subject string   `json:"subject"`
				Start   int64    `json:"start"`
				Image   string   `json:"image"`
			}{}
			if err := json.Unmarshal(data.Data, &jobData); err == nil {
				receiversMap := map[string]*models.Receiver{}

				for _, groupID := range jobData.Groups {
					group := &models.Group{}
					db.DB().Groups().Find(bson.M{"id": groupID}).One(group)
					if group.ID != "" {
						for _, receiver := range group.Receivers {
							receiversMap[receiver.Email] = receiver
						}
					}
				}

				if len(receiversMap) > 0 {
					receivers := make([]*models.Receiver, 0)
					for _, receiver := range receiversMap {
						receivers = append(receivers, receiver)
					}
					newJob := models.NewJob(jobData.Name, jobData.Start, jobData.Subject, jobData.Image, receivers)
					db.DB().Jobs().Insert(newJob)
					context.Status(http.StatusOK)
				} else {
					context.JSON(http.StatusBadRequest, struct{}{})
				}
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "remove_job":
			jobData := struct {
				ID string `json:"id"`
			}{}
			if err := json.Unmarshal(data.Data, &jobData); err == nil {
				if _, ok := af.runningJobs.Load(jobData.ID); !ok {
					db.DB().Jobs().Remove(bson.M{"id": jobData.ID})
					db.DB().Logs().Remove(bson.M{"job_id": jobData.ID})
					context.Status(http.StatusOK)
				} else {
					context.JSON(http.StatusBadRequest, struct{}{})
				}
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "add_sender":
			senderData := struct {
				Email    string `json:"email"`
				Password string `json:"password"`
			}{}
			if err := json.Unmarshal(data.Data, &senderData); err == nil {
				sender := models.NewSender(senderData.Email, senderData.Password)
				db.DB().Senders().Insert(sender)
				context.Status(http.StatusOK)
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "get_senders":
			senders := make([]*models.Sender, 0)
			db.DB().Senders().Find(bson.M{}).All(&senders)
			for _, sender := range senders {
				sender.Password = "***"
			}
			jsonData, _ := json.Marshal(senders)
			context.JSON(http.StatusOK, json.RawMessage(jsonData))
		case "remove_sender":
			senderData := struct {
				ID string `json:"id"`
			}{}
			if err := json.Unmarshal(data.Data, &senderData); err == nil {
				db.DB().Senders().Remove(bson.M{"id": senderData.ID})
				context.JSON(http.StatusOK, struct{}{})
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "add_group":
			groupData := struct {
				Name   string   `json:"name"`
				Emails []string `json:"emails"`
			}{}
			if err := json.Unmarshal(data.Data, &groupData); err == nil {
				receivers := make([]*models.Receiver, 0)
				for _, email := range groupData.Emails {
					receivers = append(receivers, models.NewReceiver(email))
				}
				group := models.NewGroup(groupData.Name, receivers)
				db.DB().Groups().Insert(group)
				context.Status(http.StatusOK)
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "get_groups":
			groups := make([]*models.Group, 0)
			db.DB().Groups().Find(bson.M{}).All(&groups)

			viewData := make([]map[string]interface{}, 0)
			for _, group := range groups {
				viewData = append(viewData, map[string]interface{}{
					"id":     group.ID,
					"name":   group.Name,
					"emails": len(group.Receivers),
				})
			}
			jsonData, _ := json.Marshal(viewData)
			context.JSON(http.StatusOK, json.RawMessage(jsonData))
		case "remove_group":
			groupData := struct {
				ID string `json:"id"`
			}{}
			if err := json.Unmarshal(data.Data, &groupData); err == nil {
				db.DB().Groups().Remove(bson.M{"id": groupData.ID})
				context.JSON(http.StatusOK, struct{}{})
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		case "get_logs":
			logsData := struct {
				ID string `json:"id"`
			}{}
			if err := json.Unmarshal(data.Data, &logsData); err == nil {
				logs := make([]*models.Log, 0)
				db.DB().Logs().Find(bson.M{"job_id": logsData.ID}).All(&logs)
				jsonData, _ := json.Marshal(logs)
				context.JSON(http.StatusOK, json.RawMessage(jsonData))
			} else {
				context.JSON(http.StatusBadRequest, struct{}{})
			}
		default:
			context.JSON(http.StatusBadRequest, struct{}{})
		}
	}
}

func main() {
	utils.GetConfig().SetConfigName("af.config")
	utils.GetConfig().AddConfigPath("./")
	utils.GetConfig().SetEnvPrefix("af")
	utils.GetConfig().SetDefault("jwt_secret", "kajshdguabcvjhbaucbaweubcuwyebgfouweghbfuiwq")
	utils.GetConfig().SetDefault("signature", "admin9988")
	utils.GetConfig().AutomaticEnv()
	utils.Start()

	af := NewMailAF()
	af.Start()
}

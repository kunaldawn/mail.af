package runner

import (
	"github.com/kunaldawn/mail.af/pkg/db"
	"github.com/kunaldawn/mail.af/pkg/db/models"
	"github.com/kunaldawn/mail.af/pkg/email"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type Runner struct {
	job     *models.Job
	senders []*models.Sender
}

func NewRunner(job *models.Job, senders []*models.Sender) *Runner {
	return &Runner{job: job, senders: senders}
}

func (runner *Runner) Run(callback func()) {
	go func() {
		if !runner.job.Running {
			log.Println("STARTING JOB :", runner.job.ID, runner.job.Name)
			db.DB().Jobs().Update(bson.M{"id": runner.job.ID}, bson.M{"$set": bson.M{"running": true}})
		} else {
			log.Println("RESUMING JOB :", runner.job.ID, runner.job.Name)
		}

		senderCache := map[int]*email.Sender{}
		for index, sender := range runner.senders {
			if emailSender, err := email.NewSender(sender); err == nil {
				senderCache[index] = emailSender
			} else {
				log.Println("ERROR : UNABLE TO CONNECT SENDER :", runner.job.ID, runner.job.Name, sender.ID, sender.Email)
				callback()
				return
			}
		}
		log.Println("SENDERS INITIALIZED :", runner.job.ID, runner.job.Name, len(senderCache))

		senderIndex := 0
		successCount := 0
		failedCount := 0
		if len(senderCache) > 0 {
			for _, receiver := range runner.job.Receivers {
				sendLog := &models.Log{}
				db.DB().Logs().Find(bson.M{"job_id": runner.job.ID, "receiver.email": receiver.Email}).One(sendLog)
				if sendLog.JobID == runner.job.ID {
					if sendLog.Success {
						successCount++
					} else {
						failedCount++
					}
				} else {
					log.Println("SENDING EMAIL :", runner.job.ID, runner.job.Name, runner.senders[senderIndex].Email, receiver.Email)

					if sender, ok := senderCache[senderIndex]; ok {
						if err := email.Send(sender, receiver, runner.job.Subject, runner.job.Image); err == nil {
							log.Println("SENDING EMAIL SUCCESS :", runner.job.ID, runner.job.Name, runner.senders[senderIndex].Email, receiver.Email)
							sendLog := models.NewLog(runner.job.ID, receiver, true)
							db.DB().Logs().Insert(sendLog)
							successCount++
						} else {
							log.Println("SENDING EMAIL FAIL :", runner.job.ID, runner.job.Name, runner.senders[senderIndex].Email, receiver.Email, err)
							sendLog := models.NewLog(runner.job.ID, receiver, false)
							db.DB().Logs().Insert(sendLog)
							failedCount++
						}
					}
				}

				senderIndex++
				if senderIndex > len(runner.senders)-1 {
					senderIndex = 0
				}
			}
		}

		log.Println("ENDING JOB :", runner.job.ID, runner.job.Name, successCount, failedCount)
		db.DB().Jobs().Update(bson.M{"id": runner.job.ID}, bson.M{"$set": bson.M{"running": false, "end_time": time.Now().UnixNano(), "done": true, "send_success": successCount, "send_failed": failedCount}})
		callback()
	}()
}

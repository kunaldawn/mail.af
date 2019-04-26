package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"sync"
	"time"
)

var mongoInstance *Mongo
var once sync.Once

type Mongo struct {
	session *mgo.Session
}

func DB() *Mongo {
	once.Do(func() {
		mongoInstance = &Mongo{}

		sessionFound := false
		for !sessionFound {
			session, err := mgo.Dial("localhost:27017")
			if err == nil {
				mongoInstance.session = session
				sessionFound = true
				log.Println("CONNECTED TO MONGO DB")
				break
			}
			log.Println("ERROR : unable to connect to mongo db, retrying...")
			time.Sleep(time.Second)
		}
	})

	return mongoInstance
}

func (mongo *Mongo) Jobs() *mgo.Collection {
	return mongo.session.DB("mailaf").C("jobs")
}

func (mongo *Mongo) Groups() *mgo.Collection {
	return mongo.session.DB("mailaf").C("groups")
}

func (mongo *Mongo) Logs() *mgo.Collection {
	return mongo.session.DB("mailaf").C("logs")
}

func (mongo *Mongo) Receivers() *mgo.Collection {
	return mongo.session.DB("mailaf").C("receivers")
}

func (mongo *Mongo) Senders() *mgo.Collection {
	return mongo.session.DB("mailaf").C("senders")
}

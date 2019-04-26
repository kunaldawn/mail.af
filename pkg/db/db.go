/*
 __  __       _ _      _    _____
|  \/  | __ _(_) |    / \  |  ___|
| |\/| |/ _` | | |   / _ \ | |_
| |  | | (_| | | |_ / ___ \|  _|
|_|  |_|\__,_|_|_(_)_/   \_\_|

Send mails as fuck!
Author : Kunal Dawn (kunal.dawn@gmail.com)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>
*/
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

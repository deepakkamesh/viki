package viki

import (
	"database/sql"
	"flag"
	"log"
	"time"

	"github.com/deepakkamesh/viki/devicemanager"
	_ "github.com/mattn/go-sqlite3"
)

func (m *Viki) MyLogger(c chan devicemanager.DeviceData) {

	log.Printf("starting user routine logger...")
	logPath := flag.Lookup("log").Value.String()
	db, err := sql.Open("sqlite3", logPath+"/log.db")
	if err != nil {
		log.Printf("error opening sqlite db %s.", err)
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO logs(tmstmp,object,state) values(?,?,?)")
	if err != nil {
		log.Printf("error preparing sql %s", err)
		return
	}

	for {
		select {
		// Channel to recieve any events.
		case got := <-c:
			// TODO: need to deal with nonstring data.
			d, _ := got.Data.(string)
			name := m.GetObjectName(got.Object)
			if name != "" {
				log.Printf("Got data from %s %s\n", name, d)
				// Write to log db.
				if _, err := stmt.Exec(time.Now().UnixNano(), name, translateState(d)); err != nil {
					log.Printf("error executing sql %s", err)
				}
			} else {
				log.Printf("Got data from unknown object %s %s\n", got.Object, d)
			}
		}
	}
}

func translateState(st string) int {
	switch st {
	case "On":
		return 1
	case "Off":
		return 0
	case "Open":
		return 1
	case "Closed":
		return 0
	default:
		return 0
	}
}

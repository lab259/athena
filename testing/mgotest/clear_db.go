package mgotest

import (
	"strings"

	"github.com/globalsign/mgo"
	mgorscsrv "github.com/lab259/athena/rscsrv/mgo"
	mgosrv "github.com/lab259/go-rscsrv-mgo"
	"github.com/onsi/gomega"
)

func ClearDB(db *mgo.Database, ignoreCollection ...string) {
	gomega.ExpectWithOffset(1, clearDB(db, ignoreCollection...)).To(gomega.Succeed())
}

func ClearDefaultMgoService(db string, ignoreCollection ...string) {
	gomega.ExpectWithOffset(1, mgorscsrv.DefaultMgoService.RunWithSession(func(session *mgo.Session) error {
		return clearDB(session.DB(db), ignoreCollection...)
	})).To(gomega.Succeed())
}

func ClearMgoService(srv *mgosrv.MgoService, db string, ignoreCollection ...string) {
	gomega.ExpectWithOffset(1, srv.RunWithSession(func(session *mgo.Session) error {
		return clearDB(session.DB(db), ignoreCollection...)
	})).To(gomega.Succeed())
}

func clearDB(db *mgo.Database, ignoreCollection ...string) error {
	collectionsToIgnore := strings.Join(ignoreCollection, ",")
	if collections, err := db.CollectionNames(); err == nil {
		for _, collection := range collections {
			if strings.HasPrefix(collection, "system.") {
				continue
			}
			if strings.Contains(collectionsToIgnore, collection) {
				continue
			}
			err := db.C(collection).DropCollection()
			if err != nil {
				return err
			}
		}
	} else {
		return err
	}
	return nil
}

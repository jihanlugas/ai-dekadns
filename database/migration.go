package database

import (
	"ai-dekadns/model"
	"log"
)

func RunMigration() {
	var err error
	db := GetCorePostsqlConn()

	err = db.Exec("CREATE SCHEMA IF NOT EXISTS dekadns").Error
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().AutoMigrate(&model.Zone{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().AutoMigrate(&model.Type{})
	if err != nil {
		log.Fatal(err)
	}
}

func RunRollback() {
	var err error
	db := GetCorePostsqlConn()

	err = db.Migrator().DropTable(&model.Zone{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Migrator().DropTable(&model.Type{})
	if err != nil {
		log.Fatal(err)
	}
}

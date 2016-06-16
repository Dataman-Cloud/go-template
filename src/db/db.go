package db

import (
	"fmt"
	"sync"

	"github.com/Dataman-Cloud/go-template/src/config"
	log "github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/driver/mysql"
	"github.com/mattes/migrate/migrate"
)

func MysqlInit() {
	DB()
	//	upgradeDB()
}

var db *sqlx.DB

func DB() *sqlx.DB {
	if db != nil {
		return db
	}
	mutex := sync.Mutex{}
	mutex.Lock()
	db, _ = InitDB()
	defer mutex.Unlock()
	return db
}

func upgradeDB() {
	uri := fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.GetConfig().UserName,
		config.GetConfig().PassWord,
		config.GetConfig().Host,
		config.GetConfig().Port,
		config.GetConfig().DataBase)
	log.Info("upgrade db mysql drive: ", uri)
	errors, ok := migrate.UpSync(uri, "./sql")
	if errors != nil && len(errors) > 0 {
		for _, err := range errors {
			log.Error("db err", err)
		}
		log.Error("can't upgrade db", errors)

		panic(-1)
	}
	if !ok {
		log.Error("can't upgrade db")
		panic(-1)
	}
	log.Info("DB upgraded")

}

func InitDB() (*sqlx.DB, error) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		config.GetConfig().UserName,
		config.GetConfig().PassWord,
		config.GetConfig().Host,
		config.GetConfig().Port,
		config.GetConfig().DataBase)
	db, err := sqlx.Open("mysql", uri)
	if err != nil {
		log.Errorf("cat not connection mysql error: %v, uri:%s", err, uri)
		return db, err
	}
	err = db.Ping()
	if err != nil {
		log.Error("can not ping mysql error: ", err)
		return db, err
	}
	//	db.SetMaxIdleConns(int(GetConfig().Mc.MaxIdleConns))
	//	db.SetMaxOpenConns(int(GetConfig().Mc.MaxOpenConns))
	return db, err
}

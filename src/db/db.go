package db

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type SQL struct {
	instance *gorm.DB
}

func ConnectionMySQLDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil
	}

	fmt.Println("\033[32m- Successful connection to database\033[0m")

	return db
}

func (d *SQL) GetInstance() *gorm.DB {
	return d.instance
}

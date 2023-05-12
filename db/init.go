package db

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"telegramInsiderBot/insiders"
)

var datetimePrecision = 2

var DB *gorm.DB

func Init() {
	dbPwd := os.Getenv("MYSQL_PASSWORD")
	dsn := fmt.Sprintf("root:%s@tcp(localhost:3306)/insider?charset=utf8&parseTime=True&loc=Local", dbPwd)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,                // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use DB type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}))

	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&insiders.InsiderTableHeader{})
	if err != nil {
		return
	}
	DB = db
}

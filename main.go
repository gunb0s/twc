package main

import (
	"telegramInsiderBot/db"
	"telegramInsiderBot/telegrams"
)

func main() {
	db.Init()
	telegrams.ExecuteBot()

	//err = insiders.Crawl(db)
	//if err != nil {
	//	return
	//}
}

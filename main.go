package main

import (
	"ginEssential/common"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	db := common.InitDB()
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Println("sql close err" + err.Error())
		}
		sqlDB.Close()
	}()

	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run()) // 监听并在 0.0.0.0:8080 上启动服务
}

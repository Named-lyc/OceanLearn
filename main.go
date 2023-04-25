package main

import (
	"OceanLearn/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db := common.InitDB()
	defer db.Close()
	r := gin.Default()
	r = CollectRoute(r)
	panic(r.Run(":8080"))
}

// InitDB 初始化连接数据库

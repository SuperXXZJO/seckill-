package main

import (
	"summerCourse/controller"
	"summerCourse/model"
	"summerCourse/service"
	"github.com/gin-gonic/gin"
)

func main() {
	model.InitDB()  //数据库初始化
	service.InitService()   //开始秒杀

	r := gin.Default()
	r.GET("/getGoods", controller.SelectGoods) //查看所有商品
	r.POST("/order", controller.MakeOrder) //下订单

	r.Run(":8080")
}


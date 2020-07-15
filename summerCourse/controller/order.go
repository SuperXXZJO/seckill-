package controller

import (
	"summerCourse/service"
	"github.com/gin-gonic/gin"
	"strconv"
)

//下订单
func MakeOrder(ctx *gin.Context) {
	userId := ctx.PostForm("userId")
	goodsId := ctx.PostForm("goodsId")
	itemId,_ := strconv.Atoi(goodsId)
	//将用户下订单的信息传入订单管道中
	service.OrderChan <- service.User{
		UserId:  userId,
		GoodsId: uint(itemId),
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"info": "success",
	})
}



package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang-elk/client"
	"golang-elk/model"
)

func Add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var value model.Evalue
		e := ctx.BindJSON(&value)
		if e != nil {
			ctx.JSON(200, gin.H{
				"code": "99999",
				"msg":  "解析数据失败",
			})
		} else {
			bytes, e := json.Marshal(value.Attr)
			if e != nil {
				ctx.JSON(200, gin.H{
					"code": "99998",
					"msg":  "结构体转json失败",
				})
			}
			client.Put(value.Key, string(bytes))
		}
		ctx.JSON(200, gin.H{
			"code": "10000",
			"msg":  "成功",
		})
	}

}

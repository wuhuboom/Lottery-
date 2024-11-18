package controller

import (
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
)

func GetResult(c *gin.Context) {
	db := mysql.DB
	if kind, err := c.GetQuery("kind"); err {
		db = db.Where("kind=?", kind)
	}
	if kind, err := c.GetQuery("type"); err {
		db = db.Where("type=?", kind)
	}
	if kind, err := c.GetQuery("qi_shu"); err {
		db = db.Where("qi_shu=?", kind)
	}
	re := model.CaiJi{}
	db.First(&re)
	tools.ReturnError200Data(c, re, "OK")
	return
}

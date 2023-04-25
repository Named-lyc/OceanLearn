package controller

import (
	"OceanLearn/common"
	"OceanLearn/model"
	"OceanLearn/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
)

func Register(context *gin.Context) {
	DB := common.GetDB()
	//获取参数
	name := context.PostForm("name")
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")
	//数据验证
	if len(telephone) != 11 {
		//gin.h == type H map[string]any
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位置"})
		return
	}
	if len(password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	//如果没有名称，则给一个10位随机字符串
	if len(name) == 0 {
		name = util.RandomString(10)
	}

	log.Println(name, telephone, password)

	//查询手机号
	if isTelephoneExist(DB, telephone) {
		context.JSON(422, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}
	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)

	//返回结果
	context.JSON(200, gin.H{
		"message": "注册成功",
	})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

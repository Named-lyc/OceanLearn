package controller

import (
	"OceanLearn/common"
	"OceanLearn/model"
	"OceanLearn/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	//同时加密用户密码,给给定Cost返回密码的bcrypt哈希，默认cost为10，cost为哈希次数
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "加密错误"})
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassword),
	}
	DB.Create(&newUser)

	//返回结果
	context.JSON(200, gin.H{
		"code":    200,
		"message": "注册成功",
	})
}
func Login(context *gin.Context) {
	DB := common.GetDB()
	telephone := context.PostForm("telephone")
	password := context.PostForm("password")

	if len(telephone) != 11 {
		//gin.h == type H map[string]any
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须为11位置"})
		return
	}
	if len(password) < 6 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
		return
	}

	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ? ", telephone).First(&user)
	if user.ID == 0 {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码
	//bcrypt加密：加盐单向加密，不可逆加密算法，每次加密后密文都不一样
	//CompareHashAndPassword：对比bcrypt哈希字符串和密码明文是否匹配
	//bcrypt.CompareHashAndPassword(密文，明文密码)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	//发放token
	token, err := common.ReleaseToken(user)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"code": "500", "msg": "系统异常"})
		log.Printf("token generate err： %v", err)
		return
	}
	//返回结果
	context.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登录成功",
	})
}

func Info(context *gin.Context) {
	user, _ := context.Get("user")
	context.JSON(http.StatusOK, gin.H{"code": 200, "data": gin.H{"user": user}})
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

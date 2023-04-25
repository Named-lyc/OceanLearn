package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type User struct {
	gorm.Model
	Name      string `gorm:"column:name;type:varchar(20);not null"`
	Telephone string `gorm:"column:telephone;type:varchar(11);not null;unique"`
	Password  string `gorm:"column:password;size:255;not null"`
}

func main() {
	db := InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(context *gin.Context) {
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
			name = RandomString(10)
		}

		log.Println(name, telephone, password)

		//查询手机号
		if isTelephoneExist(db, telephone) {
			context.JSON(422, gin.H{"code": 422, "msg": "用户已经存在"})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)

		//返回结果
		context.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	panic(r.Run(":8080"))
}

func isTelephoneExist(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
func RandomString(i int) string {
	var letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	result := make([]byte, i)
	rand.Seed(time.Now().Unix())
	for k := range result {
		result[k] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// InitDB 初始化连接数据库
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "123456"
	charset := "utf8"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open(driverName, dsn)

	if err != nil {
		panic("连接数据库失败" + err.Error())
	}
	//自动创建数据表
	db.AutoMigrate(&User{})
	return db
}

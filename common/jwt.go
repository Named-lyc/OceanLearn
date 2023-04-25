package common

import (
	"OceanLearn/model"
	"github.com/golang-jwt/jwt"
	"time"
)

// token: header.Payload.Signature
// header: {token类型,算法名}
// Payload:包含声明(要求),声明是关于实体(通常是用户)和其他数据的声明
// Signature:签名用户验证消息在传递过程中是否被更改
// token身份验证流程：1.客户端首次登录输入账号密码，设备信息
//
//	2.服务端校验用户密码，账号与设备绑定，生成token，token与账号id登关联
//	3.服务端返回token给客户端
//	4。客户端再次携带token，设备信息访问api
//	5，服务端校验token
//	6。服务端响应api请求
//
// 定义jwt密钥
var jwtKey = []byte("jwt_secret")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

// 发放token
func ReleaseToken(user model.User) (string, error) {
	//设置token过期时间
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	//payload
	claims := &Claims{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			//token过期时间
			ExpiresAt: expirationTime.Unix(),
			//token发放时间
			IssuedAt: time.Now().Unix(),
			//谁发放的token
			Issuer: "OceanLearn",
			//主题
			Subject: "user token",
		},
	}
	//生成token，通过HS256算法返回签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//将密钥转jwtKey换成[]byte传入tokenString，生成签名
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, err
}

// 解析token,从tokenString里面解析出claims 然后返回
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return token, claims, err
}

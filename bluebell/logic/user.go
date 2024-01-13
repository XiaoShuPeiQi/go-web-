package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
	"fmt"
	// "errors"
)

var UserToken string

func Signup(p *models.ParamSignup) (err error) {
	// 1.比较数据库
	if err = mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	//2.生成userID
	userID := snowflake.GetID()
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 3.写入数据库
	if err := mysql.InsertData(user); err != nil {
		return err
	}
	// 返回信息
	return nil
}
func Login(p *models.ParamLogin) (err error) {
	// 查询是否存在
	if exist := mysql.CheckUserExist(p.Username); exist == nil {
		return mysql.ErrorUserNotExist
	}
	// 绑定对象
	user := &models.User{
		UserID:   1,
		Username: p.Username,
		Password: p.Password,
	}
	// 校验密码
	if err := mysql.CheckPassword(user); err != nil {
		return err
	}
	UserToken, _ = jwt.GenToken(user.UserID, user.Username)
	fmt.Println(UserToken)
	return nil

}

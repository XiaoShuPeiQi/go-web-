package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]*models.Community,error) {
	// 查询dao层返回数据
	return mysql.GetCommunityList()
}
func GetCommunityDetail(id int64)(*models.CommunityDetail , error){
	// 查询dao层返回数据
	return mysql.GetCommunityDetailByID(id)

}

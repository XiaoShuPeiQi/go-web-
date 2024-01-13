package mysql

import (
	"bluebell/models"
	"bluebell/settings"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	ErrorUserExist     = errors.New("用户已存在")
	ErrorWrongPassword = errors.New("密码错误")
	ErrorUserNotExist  = errors.New("用户不存在")
)

var db *sqlx.DB

const secret = "xiaoshupeiqi"

func Init(msf *settings.MysqlConfig) (err error) {
	// dsn := fmt.Sprintf(
	// 	"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
	// 	viper.GetString("mysql.user"),
	// 	viper.GetString("mysql.password"),
	// 	viper.GetString("mysql.host"),
	// 	viper.GetInt("mysql.port"),
	// 	viper.GetString("mysql.dbname"),
	// )
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		msf.User,
		msf.Password,
		msf.Host,
		msf.Port,
		msf.Dbname,
	)
	if db, err = sqlx.Connect("mysql", dsn); err != nil {
		zap.L().Error("连接mysql错误", zap.Error(err))
		return err
	}
	db.SetMaxOpenConns(msf.MaxOpenConns)
	db.SetMaxIdleConns(msf.MaxIdleConns)
	return err
}

func Close() {
	db.Close()
}

func InsertData(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)
	sqlStr := "insert into user(user_id,username,password) values(?,?,?)"
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}
func CheckUserExist(name string) (err error) {
	sqlStr := "select count(username) from user where username = ?"
	var count int
	if err = db.Get(&count, sqlStr, name); err != nil {
		return
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}
func CheckPassword(user *models.User) (err error) {
	user.Password = encryptPassword(user.Password)
	sqlStr := "select user_id ,username , password  from user where username = ? and password = ?"
	if err = db.Get(user, sqlStr, user.Username, user.Password); err != nil {
		if err == sql.ErrNoRows {
			return ErrorWrongPassword
		}
		return err
	}
	return
}

// encryptPassword 加密
func encryptPassword(password string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(password)))
}

// GetCommunityList 查询community表的所有数据
func GetCommunityList() ([]*models.Community, error) {
	sqlStr := "select community_id , community_name , create_time from community"
	var list []*models.Community
	if err := db.Select(&list, sqlStr); err != nil {
		zap.L().Error("GetCommunityList 查询数据出错", zap.Error(err))
		return nil, err
	}
	return list, nil
}

func GetCommunityDetailByID(id int64) (*models.CommunityDetail, error) {
	sqlstr := `select 
				community_id , community_name , introduction 
				from community 
				where community_id = ?`

	data := new(models.CommunityDetail)
	if err := db.Get(data, sqlstr, id); err != nil {
		zap.L().Error("GetCommunityDetailByID查询数据库出错", zap.Error(err))
		return nil, err
	}
	return data, nil
}

func InsertPost(p *models.Post) (err error) {
	sqlStr := `insert into post(post_id,title,content,author_id,community_id) values(?,?,?,?,?)`
	_, err = db.Exec(sqlStr, p.PostID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	if err != nil {
		zap.L().Error("InsertPost出错", zap.Error(err))
	}
	return
}

func GetUserNameByID(id int64) (name string, err error) {
	sqlStr := "select username from user where user_id = ?"
	post := new(models.PostDetail)
	if err = db.Get(post, sqlStr, id); err != nil {
		zap.L().Error("SearchUserByID wrong", zap.Error(err))
		return "", err
	}
	return post.AuthorName, nil
}

func GetCommunityByID(id int64) (name string, intro string, err error) {
	sqlStr := `select 
				community_name , introduction 
				from community 
				where community_id = ?`
	post := new(models.PostDetail)
	if err = db.Get(post, sqlStr, id); err != nil {
		zap.L().Error("SearchCommunityByID wrong", zap.Error(err))
		return "", "", err
	}
	return post.CommuName, post.Introduction, nil

}

func GetPostByID(id int64) (post *models.PostDetail, err error) {
	sqlStr := `select title , content , author_id , community_id , create_time from post where post_id = ?`
	post = new(models.PostDetail)
	if err := db.Get(post, sqlStr, id); err != nil {
		zap.L().Error("db.Get(post, sqlStr, id) wrong", zap.Error(err))
	}
	return
}

func GetPostList(size int64, page int64) (list []*models.Post, err error) {
	sqlStr := `select  
    			post_id , title , content , author_id , community_id , create_time
				from post
				limit ? offset ?
				`
	err = db.Select(&list, sqlStr, size, (page-1)*size)
	return
}

func GetPostListByIDs(ids []string) ([]*models.Post, error) {
	sqlStr := `select  
    			post_id , title , content , author_id , community_id ,create_time
				from post
				where post_id in (?)
				order by FIND_IN_SET(post_id, ?)
				`
	//两个问号，ids和join各占一个
	//query为构建好的语句(就是把 in(?)->in(?,?,?,?......))
	query, args, err := sqlx.In(sqlStr, ids, strings.Join(ids, ","))
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query) //将占位符根据数据库驱动重新绑定，保证执行顺利
	postList := make([]*models.Post, 0)
	err = db.Select(&postList, query, args...)
	return postList, err
}

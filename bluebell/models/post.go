package models

import "time"

type Post struct {
	PostID      int64     `json:"id,string" db:"post_id"`               //帖子id
	Title       string    `json:"title" db:"title"`                     //帖子的标题
	Content     string    `json:"content" db:"content"`                 //帖子内容
	AuthorID    int64     `json:"authorID,string" db:"author_id"`       //作者id
	CommunityID int64     `json:"communityID,string" db:"community_id"` //社区id
	Status      int64     `json:"status" db:"status"`                   //状态码
	CreateTIme  time.Time `json:"createTIme" db:"create_time"`          //创建时间
}
type PostSwagger struct {
	Title       string `json:"title" db:"title"`                     //帖子的标题
	Content     string `json:"content" db:"content"`                 //帖子内容
	CommunityID int64  `json:"communityID,string" db:"community_id"` //社区id
}

type PostDetail struct {
	*Post        //邮件详情
	VoteScore    int64
	AuthorName   string `db:"username"`
	CommuName    string `db:"community_name"`
	Introduction string `db:"introduction"`
}

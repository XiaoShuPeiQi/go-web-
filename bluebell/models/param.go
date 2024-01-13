package models

type ParamSignup struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	AcPassword string `json:"ac_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteForPost struct {
	PostID    int64 `json:"post_id,string" binding:"required"`       //给哪个帖子投票
	Direction int8  `json:"direction,string" binding:"oneof=-1 0 1"` //投什么票
}

type ParamPostList struct {
	Page  int64  `form:"page"`
	Size  int64  `form:"size"`
	Order string `form:"order"`
}

package redis

const (
	KeyTimeZset  = "bluebell:post:time"    //投票时间为分数
	KeyScoreZset = "bluebell:post:score"   //投票时间＋赞成反对为分数
	KeyPostZsetP = "bluebell:post:postID:" //每个帖子都有谁投什么票，
)

package redis

import (
	"bluebell/models"
	"github.com/go-redis/redis/v8"
)

func GetPostIDsByOrder(p *models.ParamPostList) (ids []string, err error) {
	var key string = KeyScoreZset //决定去哪里查
	//order=时间或者分数
	if p.Order == "time" {
		key = KeyTimeZset
	}
	start := (p.Page - 1) * p.Size
	end := start + p.Size - 1
	return rdb.ZRevRange(ctx, key, start, end).Result()
}

// GetPostVoteData 查询每篇帖子有多少用户投了赞成票
func GetPostVoteData(ids []string) (data []int64, err error) {
	data = make([]int64, 0, len(ids))
	pipeline := rdb.Pipeline()
	//遍历拿到篇帖子的情况 使用pipeline，减少RTT
	for _, id := range ids {
		key := KeyPostZsetP + id
		pipeline.ZCount(ctx, key, "1", "1")
	}
	exec, err := pipeline.Exec(ctx)
	if err != nil {
		return nil, err
	}
	//拿到结果并追加到data中
	for _, cmder := range exec {
		value := cmder.(*redis.IntCmd).Val()
		data = append(data, value)
	}
	return
}

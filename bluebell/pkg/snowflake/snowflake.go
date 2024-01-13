package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	// 按照2006-01-02的格式解析时间字符串,这是go中固定的一个模板
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}
	// 将时间的单位设置成ms
	sf.Epoch = st.UnixNano() / 1000000
	// 创建node
	node, err = sf.NewNode(machineID)
	return
}
// 得到id
func GetID() int64 {
	return node.Generate().Int64()
}

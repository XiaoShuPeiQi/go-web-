package controllers

type _ResponseData struct {
	Code ResCode     `json:"code"` // 业务响应状态码
	Msg  interface{} `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

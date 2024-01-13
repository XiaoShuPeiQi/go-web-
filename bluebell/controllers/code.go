package controllers

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeServerBusy
	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:      "操作成功",
	CodeInvalidParam: "请求参数错误",
	CodeServerBusy:   "服务繁忙",
	CodeNeedLogin:    "请先登录",
	CodeInvalidToken: "无效的Token",
}

func (c ResCode) getMsg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}

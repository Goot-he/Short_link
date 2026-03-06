package errormsg

const (
	SUCCESS = 200
	ERROR   = 500

	Moved_Permanently = 301
	Found             = 302
	Not_Modified      = 304

	ERROR_LONGURL_AVAILABLE      = 2001
	ERROR_SHORTURL_AVAILABLE     = 2002
	ERROR_LONGURL_NOT_AVAILABLE  = 2003
	ERROR_SHORTURL_NOT_AVAILABLE = 2004
	ERROR_SHORTURL_RUNTIME       = 2005

	ERROR_BLOOM_NOT_FOUND        = 3001
	ERROR_GENERATE_SHORTURL_FAIL = 3002
)

var ErrorMap = map[int]string{
	SUCCESS: "操作成功",
	ERROR:   "操作失败",

	Moved_Permanently: "永久重定向",
	Found:             "临时重定向",
	Not_Modified:      "未修改",

	ERROR_LONGURL_AVAILABLE:      "传入的长链接已经存在",
	ERROR_SHORTURL_AVAILABLE:     "传入的短链接已经存在",
	ERROR_LONGURL_NOT_AVAILABLE:  "长链接不存在",
	ERROR_SHORTURL_NOT_AVAILABLE: "短链接不存在",
	ERROR_SHORTURL_RUNTIME:       "短链接已经超时",

	ERROR_BLOOM_NOT_FOUND:        "通过布隆过滤器得到的结果是不存在",
	ERROR_GENERATE_SHORTURL_FAIL: "随机生成短码失败 请稍后再试",
}

func GetErrorMsg(code int) string {
	return ErrorMap[code]
}

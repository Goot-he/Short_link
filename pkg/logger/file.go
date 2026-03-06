package logger

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

const PATH = "D:\\Go_Code\\reGo_code\\new_url\\pkg\\logger\\FileLog"

// 这个函数的作用就是返回一个log对象
var fileOut = logrus.New()

func init() {
	//一般来说的普通操作是这样，大多数情况下日志都是需要切割的，否则会影响文件的查询速度已经其他问题的出现
	//file, err := os.OpenFile(PATH, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	fileRotate, err := getWriter()
	if err == nil {
		//	设置文件打印的位置
		fileOut.SetOutput(fileRotate)
	} else {
		fileOut.Error(err)
	}
	registerAdapter("file", fileOut)
}

// 文件切割操作：
func getWriter() (io.Writer, error) {
	path := PATH
	// 设置一个文件写入对象
	writer, err := rotatelogs.New(path + ".%Y%m%d%H%M")
	// rotatelogs.WithLinkName(path)//将最新的文件软连接到指定的path下，windows环境不支持
	rotatelogs.WithMaxAge(time.Second * 1200)     //日志最长的保存时间
	rotatelogs.WithRotationTime(time.Second * 60) //日志分割的时间
	return writer, err
}

package logger

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"new_url/config"
	"path"
	"strings"
)

const (
	red    = 31 //红色
	yellow = 33 //黄色
	blue   = 35 //紫色
	gray   = 37 //灰白色
)

type LogFormatter struct{}

// Format 改变提示的颜色来区分数据
func (t *LogFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//根据不同的等级设置颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}
	// 如果 entry 已经有一个缓冲区（entry.Buffer 不为 nil），则使用它；否则，创建一个新的缓冲区
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	//自定义日期格式
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	//构造日志消息，包括额外字段
	fields := make([]string, 0, len(entry.Data))
	for k, v := range entry.Data {
		fields = append(fields, fmt.Sprintf("[%s = %v]", k, v))
	}
	fieldsStr := strings.Join(fields, " ") //拼接字符串函数 将全部的字符使用" "拼接

	if entry.HasCaller() {
		// 自定义文件路径
		funcVal := entry.Caller.Function
		fileVal := fmt.Sprintf("%s:%d", path.Base(entry.Caller.File), entry.Caller.Line)
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s %s %s\n", timestamp, levelColor, entry.Level, fileVal, funcVal, entry.Message, fieldsStr)
	} else {
		fmt.Fprintf(b, "[%s] \x1b[%dm[%s]\x1b[0m %s %s\n", timestamp, levelColor, entry.Level, entry.Message, fieldsStr)
	}
	return b.Bytes(), nil
}

//=======================================================================

// 用一个map管理全部的不同类型的日志
var adapters = make(map[string]*logrus.Logger)

// 这里就是一个注册处 全部的日志类型都是调用这个函数注册
func registerAdapter(adapterName string, log *logrus.Logger) {
	adapters[adapterName] = log
}

//在使用的时候：log := adapters["adaptersName"] 直接从结构体里面拿log对象

// ParseLogLevel 判断传入的是什么等级
func ParseLogLevel(s string) (logrus.Level, error) {
	s = strings.ToLower(s) //全部转化为小写
	switch s {
	case "debug":
		return logrus.DebugLevel, nil
	case "trace":
		return logrus.TraceLevel, nil
	case "info":
		return logrus.InfoLevel, nil
	case "warning":
		return logrus.WarnLevel, nil
	case "error":
		return logrus.ErrorLevel, nil
	case "fatal":
		return logrus.FatalLevel, nil
	default:
		err := errors.New("failed to parse logmiddle level,use the default level")
		return logrus.DebugLevel, err
	}
}

// DefaultLogger 统一初始化log配置的函数 处理传入的log对象
func DefaultLogger(mylog *logrus.Logger, Level string) {
	mylog.SetReportCaller(true)         //让日志能打印出：文件名 + 行号 + 函数名
	mylog.SetFormatter(&LogFormatter{}) //设置日志格式 就上面自定义的格式
	LogLevel, _ := ParseLogLevel(Level) //解析日志级别 然后传入下面这个函数
	mylog.SetLevel(LogLevel)            //设置最低日志打印级别 Trace < Debug < Info < Warn < Error < Fatal < Panic
}

func InitLogger() *logrus.Logger {
	logger, ok := adapters[config.GlobalCfg.Logger.LogType]
	if !ok {
		fmt.Printf("logger adapter not found: %s", config.GlobalCfg.Logger.LogType)
		return nil
	}

	DefaultLogger(logger, config.GlobalCfg.Logger.Level)

	return logger
}

//func InitLogger(Log **logrus.Logger) error {
//	NewLog, ok := adapters[config.GlobalCfg.Logger.LogType]
//	if !ok {
//		return errors.New("logmiddle not exist")
//	}
//	DefaultLogger(NewLog, config.GlobalCfg.Logger.Level)
//	*Log = NewLog
//	return nil
//}

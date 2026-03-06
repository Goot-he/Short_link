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
	red    = 31
	yellow = 33
	blue   = 35
	gray   = 37
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
	fieldsStr := strings.Join(fields, " ")

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

var adapters = make(map[string]*logrus.Logger)

func registerAdapter(adapterName string, log *logrus.Logger) {
	adapters[adapterName] = log
}

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
	case "warming":
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

// DefaultLogger 处理传入的log对象
func DefaultLogger(mylog *logrus.Logger, Level string) {
	mylog.SetReportCaller(true)
	mylog.SetFormatter(&LogFormatter{})
	LogLevel, _ := ParseLogLevel(Level)
	mylog.SetLevel(LogLevel)
}

func InitLogger(Log **logrus.Logger) error {
	NewLog, ok := adapters[config.GlobalCfg.Logger.LogType]
	if !ok {
		return errors.New("logmiddle not exist")
	}
	DefaultLogger(NewLog, config.GlobalCfg.Logger.Level)
	*Log = NewLog
	return nil
}

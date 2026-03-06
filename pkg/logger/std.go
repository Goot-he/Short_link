package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var std = logrus.New()

func init() {
	registerAdapter("std", std)
	std.SetOutput(os.Stdout)
	//设置打印位置，重定向到终端。
}

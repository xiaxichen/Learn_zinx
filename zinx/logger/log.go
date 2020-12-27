package logger

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.Level = logrus.DebugLevel
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	Log.SetReportCaller(true)
	Log.Info("start Log！")
}

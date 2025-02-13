// 日志配置
package main

import (
	"fmt"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
	"github.com/gookit/slog/rotatefile"
)

func InitLogConfig() {
	//defer slog.MustClose()
	// DangerLevels 包含： slog.PanicLevel, slog.ErrorLevel, slog.WarnLevel
	//h1 := handler.MustFileHandler("./logs/errorXzhi.log", handler.WithLogLevels(slog.DangerLevels),
	//	handler.WithRotateTime(3),
	//	handler.WithCompress(true), handler.WithBackupNum(5))
	// NormalLevels 包含： slog.InfoLevel, slog.NoticeLevel, slog.DebugLevel, slog.TraceLevel
	//h2 := handler.MustFileHandler("./logs/infoXzhi.log", handler.WithLogLevels(slog.NormalLevels))
	// 注册 handler 到 logger(调度器)

	h, err := handler.NewTimeRotateFileHandler(
		fmt.Sprintf("%s/logs/run.log", GetProgramDir()),
		rotatefile.EveryDay,
		handler.WithBuffSize(0),
		handler.WithBackupNum(0),
		handler.WithCompress(true),
	)

	if err != nil {
		panic(err)
	}

	slog.AddHandler(h)
	//slog.PushHandler(h1)
	//slog.PushHandler(h2)
}

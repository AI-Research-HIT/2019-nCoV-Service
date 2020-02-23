package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AI-Research-HIT/2019-nCoV-Service/db"

	"github.com/AI-Research-HIT/2019-nCoV-Service/config"
	"github.com/AI-Research-HIT/2019-nCoV-Service/service"
	"github.com/ender-wan/ewlog"
	"github.com/ender-wan/goutility/pidf"
)

func main() {
	err := config.ParseConfig()
	if err != nil {
		ewlog.Fatal(err)
	}

	file, err := os.OpenFile(config.Config.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		ewlog.Fatal(err)
	}
	defer file.Close()

	ewlog.SetLogLevel(config.Config.LogLevel)

	ewlog.AddLogOutput(file)

	pidF := pidf.New(config.Config.PidFile)
	defer pidF.Close()

	qC := make(chan os.Signal, 1)
	signal.Notify(qC, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancelFn := context.WithCancel(context.Background())
	db.ConnectToMongo()
	service.StartService(ctx)

	ewlog.Info("Server started")

	select {
	case s := <-qC:
		ewlog.Infof("main$qC: %s", s)
		cancelFn()
	}

	ewlog.Info("Server stopping")

	time.Sleep(time.Second * 3)

	ewlog.Info("Server stopped")
}

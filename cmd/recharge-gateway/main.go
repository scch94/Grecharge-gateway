package main

import (
	"context"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/scch94/Grecharge-gateway/client"
	"github.com/scch94/Grecharge-gateway/config"
	"github.com/scch94/Grecharge-gateway/server"
	"github.com/scch94/ins_log"
)

func main() {

	//creamos el contexto para esta ejecucion
	ctx := context.Background()

	//load configuration
	if err := config.Upconfig(ctx); err != nil {
		ins_log.Errorf(ctx, "error loading configuration: %v", err)
		return
	}

	//crear directoruio de log y cambio cada hora
	go initializeAndWatchLogger(ctx)
	//vijilamos el archivo de config
	go config.WatchConfig(ctx)
	ins_log.SetService("recharge-gateway")
	ins_log.SetLevel(config.Config.LogLevel)
	ctx = ins_log.SetPackageNameInContext(ctx, "main")

	ins_log.Infof(ctx, "starting recharge gateway version %s", getVersion())
	ins_log.Infof(ctx, "configuration: %+v", config.Config)

	//inicamos el client
	client.InitHtppClient()

	//iniciamos el servidor
	go startServer(ctx)

	// Capture exit signals for graceful shutdown
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)

	// Keep the program running
	sig := <-exitSignal
	ins_log.Infof(ctx, "closing micropagos message gateway version %v, received shutdown signal: %v,", getVersion(), sig)

}
func initializeAndWatchLogger(ctx context.Context) {
	var file *os.File
	var logFileName string
	var err error
	for {
		select {
		case <-ctx.Done():
			return
		default:
			logDir := "../../log"

			// Create the log directory if it doesn't exist
			if err = os.MkdirAll(logDir, 0755); err != nil {
				ins_log.Errorf(ctx, "error creating log directory: %v", err)
				return
			}

			// Define the log file name
			today := time.Now().Format("2006010215")
			replacer := strings.NewReplacer(" ", "_")
			today = replacer.Replace(today)
			logFileName = filepath.Join(logDir, today+config.Config.Log_name+".log")

			// Open the log file
			file, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				ins_log.Errorf(ctx, "error opening log file: %v", err)
				return
			}

			// Create a writer that writes to both file and console
			multiWriter := io.MultiWriter(os.Stdout, file)
			ins_log.StartLoggerWithWriter(multiWriter)

			// Esperar hasta el inicio de la próxima hora
			nextHour := time.Now().Truncate(time.Hour).Add(time.Hour)
			time.Sleep(time.Until(nextHour))

			// Close the previous log file
			file.Close()
		}
	}
}

func startServer(ctx context.Context) {
	if err := server.StartServer(ctx); err != nil {
		ins_log.Errorf(ctx, "error starting server: %s", err.Error())
	}
}

func getVersion() string {
	return "1.0.0"
}

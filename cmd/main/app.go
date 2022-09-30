package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"rest-api-tutorial/config"
	"rest-api-tutorial/internal/user"
	"rest-api-tutorial/pkg/logging"
	"time"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("Create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	logger.Info("Register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Trace("Start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Infof("Args %v", os.Args[0])
		logger.Infof("Dir %v", filepath.Dir(os.Args[0]))
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("Create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("Listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("Server is listening socket %s", socketPath)
	} else {
		logger.Info("Listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIp, cfg.Listen.Port))
		logger.Infof("Server is listening port %s:%s", cfg.Listen.BindIp, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}

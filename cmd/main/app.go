package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"

	v1 "github.com/taaanechka/rest-api-go/internal/api-server/api/http/v1"
	"github.com/taaanechka/rest-api-go/internal/api-server/repositories/user/mongodb"
	userservice "github.com/taaanechka/rest-api-go/internal/api-server/services"
	"github.com/taaanechka/rest-api-go/internal/api-server/services/ports/userstorage"
	"github.com/taaanechka/rest-api-go/internal/config"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

func main() {
	lg := logging.GetLogger()
	lg.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	storage, err := mongodb.NewStorage(lg, userstorage.Config(cfg.Users))
	if err != nil {
		lg.Errorf("failed to init storage: %v", err)
		return
	}

	service := userservice.NewService(lg, storage)

	lg.Info("register user handler")
	handler := v1.NewHandler(lg, service)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	lg := logging.GetLogger()
	lg.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		lg.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			lg.Fatal(err)
		}
		lg.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		lg.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		lg.Infof("server is listening unix socket: %s", socketPath)
	} else {
		lg.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		lg.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		lg.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	lg.Fatal(server.Serve(listener))
}

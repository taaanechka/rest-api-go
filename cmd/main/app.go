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
	"github.com/taaanechka/rest-api-go/internal/config"
	"github.com/taaanechka/rest-api-go/internal/user"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	// cfgMongo := cfg.MongoDB
	// mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port,
	// 	cfgMongo.Username, cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	// if err != nil {
	// 	panic(err)
	// }
	// storage := user.NewStorage(mongoDBClient, cfgMongo.Collection, logger)

	// // Create
	// user1 := user.User{
	// 	ID: "",
	// 	Email: "dev.test@mail.ru",
	// 	Username: "dev",
	// 	PasswordHash: "12345",
	// }
	// user1ID, err := storage.Create(context.Background(), user1)
	// if err != nil {
	// 	panic(err)
	// }
	// logger.Info(user1ID)

	// // FindOne
	// user1Found, err := storage.FindOne(context.Background(), user1ID)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(user1Found)

	// // Update
	// user1Found.Email = "newEmail@mail.ru"
	// err = storage.Update(context.Background(), user1Found)
	// if err != nil {
	// 	panic(err)
	// }

	// // Delete
	// err = storage.Delete(context.Background(), user1ID)
	// if err != nil {
	// 	panic(err)
	// }

	// // FindAll
	// users, _ := storage.FindAll(context.Background())
	// fmt.Printf("\nusers:\n")
	// for _, u := range users {
	// 	fmt.Printf("%v\n", u)
	// }
	// fmt.Println()

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
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

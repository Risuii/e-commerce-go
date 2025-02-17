package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	Config "e-commerce/config"
	Constants "e-commerce/constants"
	Library "e-commerce/library"
	CustomErrorPackage "e-commerce/pkg/custom_error"
	LoggerPackage "e-commerce/pkg/logger"
	Wire "e-commerce/wire"
)

func main() {
	path := "main"
	// INIT LIBRARY
	library := Library.New()
	// INIT CONFIGURATION
	config := Config.New(library)
	// INIT LOGGER
	LoggerPackage.New(library)

	// SETUP CONFIGURATION
	err := config.Setup()
	// WHEN SETUP CONFIGURATION RETURNS AN ERROR
	if err != nil {
		err = err.(*CustomErrorPackage.CustomError).UnshiftPath(path)
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:  err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Title: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}
	// SET PORT THAT WE GET FROM .env FILE
	port := fmt.Sprintf(":%s", config.GetConfig().App.AppPort)
	// INIT ROUTER
	router := Wire.InjectRoute(config, library)
	// SETUP ROUTER
	router.Setup()

	// INIT QUEUE

	// SETUP QUEUE

	srv := &http.Server{
		Addr:    port,
		Handler: router.GetEngine(),
	}

	// HTTP Server
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			err = CustomErrorPackage.New(Constants.ErrServeFailed, err, path, library)
			LoggerPackage.WriteLog(logrus.Fields{
				Constants.Path:  err.(*CustomErrorPackage.CustomError).GetPath(),
				Constants.Title: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
			}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// wait for exit signal
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// shutdown server
	if err := srv.Shutdown(ctx); err != nil {
		// Todo: ErrServeFailed -> ErrStopFailed
		err = CustomErrorPackage.New(Constants.ErrServeFailed, err, path, library)
		LoggerPackage.WriteLog(logrus.Fields{
			Constants.Path:  err.(*CustomErrorPackage.CustomError).GetPath(),
			Constants.Title: err.(*CustomErrorPackage.CustomError).GetDisplay().Error(),
		}).Panic(err.(*CustomErrorPackage.CustomError).GetPlain())
	}

	// wait until all queue processed

	fmt.Println(Constants.ShuttingDown)
}

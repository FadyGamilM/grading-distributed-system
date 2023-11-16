package main

import (
	"context"
	"flag"
	"godistributed/logger"
	"godistributed/service"
	"log"
)

var (
	port = flag.String("port", "5000", "port of the server")
	host = flag.String("host", "127.0.0.1", "host of the server")
)

func main() {
	// parse the flags
	flag.Parse()

	// initialize the logger instance by callung Run method and pass the file destination
	logger.Run("./app.log")

	ctx, err := service.Start(
		// we dont have context so i will pass a background context
		context.Background(),
		*host,
		*port,
		"logger-service",
		logger.RegisterHandlers,
	)
	// if there is an error instantiating the logger-service, we will log using the standard logger
	if err != nil {
		log.Fatalf("error trying to start the logger-service : %v\n", err)
	}

	// now wait untill the context is canceled
	<-ctx.Done()

	// log for info thar we are shutting down the service
	log.Println("logger-service is shutdown !")
}

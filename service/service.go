// this file will responsible for registering a web handler to the server

package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// each service we will create will have a registerHandlers method which register the required routes and their handlers
func Start(
	ctx context.Context,
	host, port, serviceName string,
	registerHandlers func() *gin.Engine,
) (context.Context, error) {
	// register the handlers of this service
	handler := registerHandlers()

	// call the backend logic of starting a service
	ctx = startService(ctx, host, port, serviceName, handler)

	// return the response
	return ctx, nil
}

// this backend logic is used to handle a web server
// 1. define a new server instance
// 2. set the port
// 3. spin a concurrent go routine to give the client the ability to shutdown the server
// 4. spin a concurrent go routine to listen for running the server and keep it up and listen for any unexpected shutdowns
func startService(
	ctx context.Context,
	host, port, serviceName string,
	h *gin.Engine,
) context.Context {
	// 1. and 2.
	server := http.Server{
		Addr:    fmt.Sprintf("%v:%v", host, port),
		Handler: h,
	}

	// define a context with cancelation option for the 2 concurrent  go routines logic
	ctx, cancel := context.WithCancel(ctx)

	// 3.
	go func() {
		fmt.Printf("[%v service] started, press any key to shutdown the server \n", serviceName)
		var shutdown_key string
		fmt.Scanln(&shutdown_key)
		// shutdown and cancel the context to trigger the Done channel
		server.Shutdown(ctx)
		cancel()
	}()

	// 4.
	go func() {
		// start the service and this will hang here until the service is shutdowns so we will cancel the context which will trigger the ctx.Done channel
		log.Println(server.ListenAndServe())
		cancel()
	}()

	// finally return the context
	return ctx
}

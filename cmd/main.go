package main

import (
	appServer "chatsystem/internal"
	exp "chatsystem/internal/exceptions"
	"fmt"
	"os"
	"os/signal"
	"time"
)

/*
|********************************
| FitSwift Fashion Seekers Systems
*********************************
|
|
|
|
*/

func main() {

	defer func() {
		if err := recover(); err != nil {
			exp.Loggers.System.Warn(fmt.Errorf("the system almost crashed due to: %v", err))
		}
	}()
	e := appServer.Start()
	//Wait for interrupt signal to gracefully shutdown the server with
	//a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	// healthCheck = "unhealthy"
	time.Sleep(5 * time.Second)
	appServer.Stop(e)
}

package main

import (
	"fourth-exam/user-service-evrone/internal/app"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	config := config.New()

	app, err := app.NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	//running
	go func() {
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	app.Logger.Info("User service stops !")

	// app stops
	app.Stop()
}

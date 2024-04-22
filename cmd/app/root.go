package app

import (
	"fmt"
	"fourth-exam/user-service-evrone/internal/app"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var rootCmd = &cobra.Command{
	Use:   "grpc-server",
	Short: "This command to run grpc server",
	Run: func(cmd *cobra.Command, args []string) {
		config := config.New()

		app, err := app.NewApp(config)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			if err := app.Run(); err != nil {
				app.Logger.Error("error while app run", zap.Error(err))
			}
		}()

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs

		app.Logger.Info("user service stops")

		// stop app
		app.Stop()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error while executing CLI '%s", err)
		os.Exit(1)
	}
}

package app

import (
	"fourth-exam/user-service-evrone/internal/app"
	"fourth-exam/user-service-evrone/internal/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	USER_CREATE_CONSUMER = "user_create_consumer"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer", // command name that we will use to invoke this command 'go run cmd/main.go consumer ...'
	Short: "To run consumer give the name followed by arguments consumer", // example usage of this command: 'go run cmd/main.go help'
	Long: `Example : 
		go run cmd/main.go consumer name_of_consumer`, // example usage of this command: 'go run cmd/main.go help consumer'
	Args: cobra.ExactArgs(1), // number of arguments the command expects

	Run: func(cmd *cobra.Command, args []string) { // this function will be executed when this command is invoked
		consumerName := args[0] // extracts the first argument from the 'args' slice

		switch consumerName { // switch statement based on consumerNames 
		case USER_CREATE_CONSUMER:
			UserCreateConsumerRun()
		default:
			log.Fatalf("No consumer with name '%s'", consumerName)
		}
	},
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}

func UserCreateConsumerRun() {
	config := config.New()

	app, err := app.NewUserConsumer(config)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		app.Run()
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	app.Logger.Info("user consumer stops")

	// stop app
	app.Close()
}

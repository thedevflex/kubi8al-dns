package main

import (
	"os"

	"github.com/thedevflex/kubi8al-dns/server"
	logs "github.com/thedevflex/kubi8al-webhook/utils/logger"
)

func main() {
	logs.InitLogger()
	if os.Getenv("WEBHOOK_SECRET") == "" {
		logs.Fatal("WEBHOOK_SECRET environment variable is not set")
		os.Exit(1)
	}
	if os.Getenv("EMMITER_API_ADDRESS") == "" {
		logs.Fatal("EMMITER_API_ADDRESS environment variable is not set")
		os.Exit(1)
	}

	logs.Info("Starting kubi8al-webhook server")
	logs.Infof("Listening on port: %s", os.Getenv("PORT"))
	logs.Infof("EMMITER_API_ADDRESS: %s", os.Getenv("EMMITER_API_ADDRESS"))

	app := server.New()
	server.Setup(app)

	if err := server.Start(app); err != nil {
		logs.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}

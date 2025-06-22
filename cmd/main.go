package main

import (
	"os"

	"github.com/thedevflex/kubi8al-dns/server"
	logs "github.com/thedevflex/kubi8al-webhook/utils/logger"
)

func main() {
	logs.InitLogger()
	app := server.New()
	server.Setup(app)

	if err := server.Start(app); err != nil {
		logs.Fatalf("Error starting server: %v", err)
		os.Exit(1)
	}
}

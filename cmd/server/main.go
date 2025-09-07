package main

import (
	"final-project/internal/metrics"
	"final-project/internal/static"
	"final-project/pkg/api"
	"final-project/pkg/db"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := db.Init("scheduler.db"); err != nil {
		fmt.Println(err)
	}
	api.Init()
	metrics.Init()
	if err := static.Static(); err != nil {
		fmt.Println(err)
	}
	log.Fatal(http.ListenAndServe(":7540", nil))
}

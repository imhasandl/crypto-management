package main

import (
	"log"
	"net/http"
	"os"

	"github.com/imhasandl/crypto-management/database"
	_ "github.com/imhasandl/crypto-management/docs"
	"github.com/imhasandl/crypto-management/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("can't get .env data: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Set port in .env file")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("Set database url in .env file")
	}

	db, err := database.InitDatabase(dbURL)
	if err != nil {
		log.Fatalf("can't connect to database: %v", err)
	}
	defer db.Close()

	apiConfig := handlers.NewConfig(db)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /currency/add", apiConfig.AddCurrency)
	mux.HandleFunc("POST /currency/remove", apiConfig.RemoveCurrency)
	mux.HandleFunc("POST /currency/price", apiConfig.GetCurrencyPrice)

	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("can't start server: %v", err)
	}
}

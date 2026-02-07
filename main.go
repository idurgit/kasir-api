package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// EndpointDetail represents a single endpoint's metadata
type EndpointDetail struct {
	Path        string `json:"path"`
	Description string `json:"description"`
}

// EndpointGroup groups endpoints by HTTP method
type EndpointGroup map[string]EndpointDetail

// APIInfoResponse contains the API metadata
type APIInfoResponse struct {
	Endpoint    map[string]EndpointGroup `json:"endpoint"`
	Environment string                   `json:"environment"`
	Message     string                   `json:"message"`
	Version     string                   `json:"version"`
}

type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// handleAPIInfo returns API metadata including endpoints, environment, version
func handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Build endpoint metadata for current implementation
	endpoints := map[string]EndpointGroup{
		"GET": {
			"list_products": {
				Path:        "/api/product",
				Description: "get all products",
			},
			"get_product": {
				Path:        "/api/product/{id}",
				Description: "get a single product",
			},
			"health": {
				Path:        "/health",
				Description: "health check endpoint",
			},
		},
		"POST": {
			"create_product": {
				Path:        "/api/product",
				Description: "create a new product",
			},
		},
		"PUT": {
			"update_product": {
				Path:        "/api/product/{id}",
				Description: "update all fields",
			},
		},
		"DELETE": {
			"delete_product": {
				Path:        "/api/product/{id}",
				Description: "delete a product",
			},
		},
	}

	// Get configuration from environment or use defaults
	environment := getEnv("API_ENV", "development")
	version := getEnv("API_VERSION", "1.0.0")
	message := getEnv("API_MESSAGE", "simple API")

	response := APIInfoResponse{
		Endpoint:    endpoints,
		Environment: environment,
		Message:     message,
		Version:     version,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
	}

	// setup database connection
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize Database:", err)
	}
	defer db.Close()
	
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// setup routes
	http.HandleFunc("/api/product", productHandler.HandleProducts)
	http.HandleFunc("/api/product/", productHandler.HandleProductByID)
	http.HandleFunc("/api/info", handleAPIInfo)
		
	// localhost:8080 / health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "API is running"})
	})
	

	fmt.Println("Starting server on localhost:" + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Println("Error starting server")
	}
}

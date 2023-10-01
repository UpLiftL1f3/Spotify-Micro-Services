package data

import (
	"fmt"
	"os"
	"strings"
)

var (
	ConnectionString string
	Database         string
)

func LoadEnvVariables() {
	// Load environment variables from the .env file
	// err := godotenv.Load("/.dockerenv")
	// if err != nil {
	// 	log.Fatal("Error loading .env file", err)
	// }

	// Access environment variables
	username := os.Getenv("MONGO_USERNAME")
	username = strings.ToLower(username)
	password := os.Getenv("MONGO_PASSWORD")
	databaseName := os.Getenv("MONGO_DATABASE1")
	clusterURI := os.Getenv("MONGO_CLUSTER_URI")
	options := os.Getenv("MONGO_OPTIONS")

	// Construct the MongoDB connection string
	ConnectionString = fmt.Sprintf("mongodb+srv://%s:%s@%s/%s", username, password, clusterURI, options)
	Database = databaseName
}

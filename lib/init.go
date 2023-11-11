package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
)

func Init() *authn.AuthN {
	// ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancelFn()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Get token
	token := os.Getenv("PANGEA_AUTHN_TOKEN")
	if token == "" {
		log.Fatal("Unauthorized: No token present")
	}

	// Create config and client
	client := authn.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	return client
}

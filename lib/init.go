package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
)

func Init() *authn.AuthN {
	// ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancelFn()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("\033[31m", "Error loading .env file", "\033[0m")
	}
	// Get token
	token := os.Getenv("PANGEA_AUTHN_TOKEN")
	if token == "" {
		log.Fatal("\033[31m", "Unauthorized: No token present", "\033[0m")
	}

	// Create config and client
	client := authn.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	return client
}

func InitVault() vault.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("\033[31m", "Error loading .env file", "\033[0m")
	}
	// Get token
	token := os.Getenv("PANGEA_VAULT_TOKEN")
	if token == "" {
		log.Fatal("\033[31m", "Unauthorized: No token present", "\033[0m")
	}

	// Create config and client
	client := vault.New(&pangea.Config{
		Token:  token,
		Domain: os.Getenv("PANGEA_DOMAIN"),
	})

	return client
}

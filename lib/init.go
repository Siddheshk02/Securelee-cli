package lib

import (
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
)

func Init() *authn.AuthN {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("\033[31m", "Error loading .env file", "\033[0m")
	// }
	// // Get token
	// token := os.Getenv("PANGEA_AUTHN_TOKEN")
	// if token == "" {
	// 	log.Fatal("\033[31m", "Unauthorized: No token present", "\033[0m")
	// }

	token := "pts_xajlrac4we4mufoebqgejbrh2ieq72c4"
	domain := "aws.us.pangea.cloud"
	// Create config and client
	client := authn.New(&pangea.Config{
		Token:  token,
		Domain: domain,
	})

	return client
}

func InitVault() vault.Client {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal("\033[31m", "Error loading .env file", "\033[0m")
	// }
	// // Get token
	// token := os.Getenv("PANGEA_VAULT_TOKEN")
	// if token == "" {
	// 	log.Fatal("\033[31m", "Unauthorized: No token present", "\033[0m")
	// }

	token := "pts_xajlrac4we4mufoebqgejbrh2ieq72c4"
	domain := "aws.us.pangea.cloud"
	// Create config and client
	client := vault.New(&pangea.Config{
		Token:  token,
		Domain: domain,
	})

	return client
}

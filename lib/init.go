package lib

import (
	"log"
	"os"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/service/authn"
)

func Init() *authn.AuthN {
	// ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
	// defer cancelFn()

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

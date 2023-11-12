package lib

import (
	"context"
	"regexp"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
)

func LoginWithPass(Email string, Password string) (*authn.UserLoginResult, error) {

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	input := authn.UserLoginPasswordRequest{
		Email:    Email,
		Password: Password,
	}
	resp, err := client.User.Login.Password(ctx, input)
	if err != nil {
		return nil, err
	}

	return resp.Result, nil
}

func IsValidPassword(password string) bool {
	// Check if password has at least 8 characters
	if len(password) < 8 {
		return false
	}

	// Check if password has at least 1 numeric character
	hasNumeric, _ := regexp.MatchString(`[0-9]`, password)
	if !hasNumeric {
		return false
	}

	// Check if password has at least 1 special character
	hasSpecialChar, _ := regexp.MatchString(`[!@#$%^&*(),.?":{}|<>]`, password)
	if !hasSpecialChar {
		return false
	}

	return true
}

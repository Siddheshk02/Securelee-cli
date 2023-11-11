package lib

import (
	"context"
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

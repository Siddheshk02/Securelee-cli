package lib

import (
	"context"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
)

func PassReset(Id string, Password string) error {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	input := authn.UserPasswordResetRequest{
		UserID:      Id,
		NewPassword: Password,
	}
	_, err := client.User.Password.Reset(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

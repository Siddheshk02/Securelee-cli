package lib

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v3/pangea"
	"github.com/skratchdot/open-golang/open"
)

func Login() {
	open.Start("https://pdn-vehpksfu665ae7k5jewmycb4fxqircam.login.aws.us.pangea.cloud/authorize?state=xxxxxxxxxxxxx")

	return
}

type TokenInfo struct {
	Token   string    `json:"token"`
	Email   string    `json:"email"`
	Name    string    `json:"name"`
	User_ID string    `json:"userID"`
	Expiry  time.Time `json:"expiry"`
}

// Check Token for Current User
func Check() bool {
	user, err := user.Current()
	if err != nil {
		log.Fatal("\033[31m", err, "\033[0m")
	}
	filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

	if !FileExists(filePath) {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("\033[31m", "Error while fetching the Token!!", "\033[0m")
			os.Exit(0)
		}

		var token TokenInfo
		err = json.Unmarshal(file, &token)
		if err != nil {
			log.Fatal("\033[31m", err, "\033[0m")
		}

		//Check Token Expiry
		if time.Now().After(token.Expiry) {
			os.Remove(filePath)
			return false
		}

		//Check Token is Valid or not
		_, ch := CheckToken(token.Token)
		if ch != "" {
			return false
		}

		return true
	}

	return false
}

// Check Token Validity
func CheckToken(token string) (*authn.ClientTokenCheckResult, string) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	input := authn.ClientTokenCheckRequest{
		Token: token,
	}

	resp, err := client.Client.Token.Check(ctx, input)

	if err != nil && resp == nil {
		return nil, "No User"
	}

	if *resp.Status == "Success" {
		return resp.Result, ""
	}

	return nil, "Invalid Token"

}

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func LoginWithEmail(Email string, Password string) (string, string, error) {

	//check if user is new or no
	var usertype string
	resp, err := NewUser(Email)
	if resp == "" && err != nil {
		return "", "", err
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	if resp == "InvalidUser" {
		usertype = "New User"

		var first, last string
		fmt.Print("\033[33m", "\n > Enter your First Name: ", "\033[0m")
		fmt.Scan(&first)
		fmt.Println("")
		fmt.Print("\033[33m", " > Enter your Last Name: ", "\033[0m")
		fmt.Scan(&last)

		profile := &authn.ProfileData{
			"first_name": first,
			"last_name":  last,
		}

		input := authn.UserCreateRequest{
			Email:         Email,
			Profile:       profile,
			Authenticator: Password,
			IDProvider:    authn.IDPPassword,
		}

		//creating the user profile
		out, err := client.User.Create(ctx, input)
		if err != nil || out == nil {
			fmt.Println("\033[31m", "\nFailed to create a new user", "\033[0m")
			os.Exit(0)
		}

		//adding password for the user profile
		err = PassReset(out.Result.ID, Password)
		if err != nil {
			return "", "", err
		}

	} else if resp == "Success" {
		usertype = "Old User"
	}

	//login using password
	result, str, err := LoginWithPass(Email, Password)
	if err != nil {
		return "", "", err
	} else if result == nil && str != "" && err == nil {
		return "", str, nil
	}

	return result.ActiveToken.Token, usertype, nil

}

// Logout the Current User's Session
func Logout() error {
	res := Check()
	if !res {
		fmt.Print("\033[31m", "\n > No User logged in, You must Login to use Securelee Vault Services.\n", "\033[0m")
		fmt.Print("\033[36m", "\n > Use 'Securelee-cli login' command to complete the Authentication.\n", "\033[0m")
		fmt.Println("")
		os.Exit(0)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	user, err := user.Current()
	if err != nil {
		log.Fatal("\033[31m", err, "\033[0m")
	}
	filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	var TokenData struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(jsonData, &TokenData)
	if err != nil {
		return err
	}

	token := TokenData.Token

	input := authn.ClientSessionLogoutRequest{
		Token: token,
	}
	_, err = client.Client.Session.Logout(ctx, input)
	if err != nil {
		// return err
		fmt.Print("\033[31m", "\n > No User logged in, You must Login to use Securelee Vault Services.\n", "\033[0m")
		fmt.Print("\033[36m", "\n > Use 'Securelee-cli login' command to complete the Authentication.\n", "\033[0m")
		fmt.Println("")
		os.Exit(0)
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	return nil

}

type APIError struct {
	RequestID    string `json:"request_id"`
	RequestTime  string `json:"request_time"`
	ResponseTime string `json:"response_time"`
	Status       string `json:"status"`
	Summary      string `json:"summary"`
}

func NewUser(Email string) (string, error) {
	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	input := authn.UserProfileGetRequest{
		Email: pangea.String(Email),
	}
	resp, err := client.User.Profile.Get(ctx, input)

	if resp == nil {
		re := regexp.MustCompile(`\{[^{}]*\}`)
		match := re.Find([]byte(err.Error()))

		if match == nil {
			return "", errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err = json.Unmarshal(match, &apiError)
		if err != nil {
			return "", err

		}

		return apiError.Status, nil
	}

	return *resp.Status, nil
}

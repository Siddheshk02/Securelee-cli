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
		log.Fatal(err)
	}
	filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

	if !FileExists(filePath) {
		file, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal("Error while fetching the Token!!")
		}

		var token TokenInfo
		err = json.Unmarshal(file, &token)
		if err != nil {
			log.Fatal(err)
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
		// log.Fatal(err, "Error, Please Try Again.")
		// fmt.Println("\n > No User logged in, You must Login to use Securelee Vault Services.")
		// os.Exit(0)
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
		fmt.Print("> Enter your First Name: ")
		fmt.Scan(&first)
		fmt.Println("")
		fmt.Print("> Enter your Last Name: ")
		fmt.Scan(&last)
		fmt.Println("")

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
			fmt.Println("Failed to create a new user")
			log.Fatal(err)
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
	result, err := LoginWithPass(Email, Password)
	if err != nil {
		return "", "", err
	}

	token := fmt.Sprintf("%s", result.ActiveToken.Token)

	return token, usertype, nil

}

// Logout the Current User's Session
func Logout() error {
	res := Check()
	if !res {
		fmt.Println("\n > No User logged in, You must Login to use Securelee Vault Services.")
		fmt.Println("")
		os.Exit(0)
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFn()
	client := Init()

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
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
		fmt.Println("\n > No User logged in, You must Login to use Securelee Vault Services.")
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
			// log.Fatal("No JSON data found in the error message")
			return "", errors.New("No JSON data found in the error message")
		}

		var apiError APIError
		err = json.Unmarshal(match, &apiError)
		if err != nil {
			// log.Fatal("Error unmarshalling JSON:", err)
			return "", err

		}

		return apiError.Status, nil
	}

	// fmt.Println(apiError.Status)

	return *resp.Status, nil
}

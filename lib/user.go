package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
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

type ResponseData struct {
	RequestID    string `json:"request_id"`
	RequestTime  string `json:"request_time"`
	ResponseTime string `json:"response_time"`
	Status       string `json:"status"`
	Result       ResultData
}

type ResultData struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Life      int    `json:"life"`
	Expire    string `json:"expire"`
	Identity  string `json:"identity"`
	Email     string `json:"email"`
	Profile   ProfileData
	CreatedAt string `json:"created_at"`
}

type ProfileData struct {
	LastLoginCity    string `json:"Last-Login-City"`
	LastLoginCountry string `json:"Last-Login-Country"`
	FirstName        string `json:"first_name"`
	LastName         string `json:"last_name"`
}

// Check Token Validity
func CheckToken(token string) (ResponseData, string) {

	url := "https://authn.aws.us.pangea.cloud/v2/client/token/check"
	method := "POST"

	payload, _ := json.Marshal(map[string]string{
		"token": token,
	})

	AuthToken := "Bearer " + "pts_xajlrac4we4mufoebqgejbrh2ieq72c4"

	client1 := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", AuthToken)

	res, _ := client1.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var responseData ResponseData
	_ = json.Unmarshal(body, &responseData)

	if responseData.Status == "Success" {
		return responseData, ""
	}

	return ResponseData{}, "Invalid Token"
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
	ref := Check()
	if !ref {
		fmt.Print("\033[31m", "\n > No User logged in, You must Login to use Securelee Vault Services.\n", "\033[0m")
		fmt.Print("\033[36m", "\n > Use 'Securelee-cli login' command to complete the Authentication.\n", "\033[0m")
		fmt.Println("")
		os.Exit(0)
	}

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

	url := "https://authn.aws.us.pangea.cloud/v2/client/session/logout"
	method := "POST"
	payload, _ := json.Marshal(map[string]string{
		"token": TokenData.Token,
	})

	AuthToken := "Bearer " + "pts_xajlrac4we4mufoebqgejbrh2ieq72c4"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))

	if err != nil {
		fmt.Println("Error !!! Please Try Again.")
		os.Exit(0)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", AuthToken)

	res, _ := client.Do(req)

	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var responseData ResponseData
	_ = json.Unmarshal(body, &responseData)

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

	url := "https://authn.aws.us.pangea.cloud/v2/user/profile/get"
	method := "POST"

	payload, _ := json.Marshal(map[string]string{
		"email": Email,
	})

	AuthToken := "Bearer " + "pts_xajlrac4we4mufoebqgejbrh2ieq72c4"

	client1 := &http.Client{}
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", AuthToken)

	res, _ := client1.Do(req)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	var responseData ResponseData
	_ = json.Unmarshal(body, &responseData)

	return responseData.Status, nil
}

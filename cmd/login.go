/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"time"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/Siddheshk02/Securelee-cli/mailing"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/spf13/cobra"
)

// SignUpCmd represents the SignUp command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Securelee Vault.",
	Long:  `Login to Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {

		res := lib.Check()
		if res == true {
			fmt.Println("\n> Already logged in!")
			fmt.Println("")
			os.Exit(0)
		}

		var choice int
		var token string
		var err error
		var result *authn.ClientTokenCheckResult
		var ch, usertype string
		fmt.Println("\n Select any one option: ")
		fmt.Println("\n> 1. Login using Socials through Browser")
		fmt.Println("> 2. Login using Email and Password on the Terminal")
		fmt.Println("")
		fmt.Print("> Enter your choice (for e.g. 1): ")
		fmt.Scanf("%d", &choice)
		fmt.Println("")
		if choice == 1 {
			fmt.Print("\n> Press Enter to Login using Browser (You can get back to CLI after Successful Authentication) : ")
			fmt.Scanf(" ")
			lib.Login()
			time.Sleep(10 * time.Second)

			fmt.Println("")
			fmt.Print("> Please Enter the Token from the Securelee Authentication Tab: ")
			fmt.Scan(&token)

			result, ch = lib.CheckToken(token)
			if ch != "" {
				log.Fatal(ch)
			}

		} else if choice == 2 {
			var email, password string
			fmt.Print("> Enter your Email Address : ")
			fmt.Scan(&email)
			fmt.Println("")
			fmt.Print("> Enter your Password ")
			fmt.Print("\n { Password must have \n   - at least 8 characters,\n   - at least 1 number characters,\n   - at least 1 special characters } : ")
			fmt.Scan(&password)
			fmt.Println("")
			check := lib.IsValidPassword(password)
			if !check {
				fmt.Println("> Invalid Password.")
				os.Exit(0)
			}

			token, usertype, err = lib.LoginWithEmail(email, password)

			if token == "" && usertype == "" && err != nil {
				log.Fatalln(err.Error())
			}

			if token != "" {
				result, ch = lib.CheckToken(token)
				if ch != "" {
					log.Fatal(ch)
				}
			}

		} else {
			fmt.Println("> Invalid Choice Entered!!, Please try again")
			fmt.Println("")
			os.Exit(0)
		}

		parsedTime, err := time.Parse(time.RFC3339, result.Expire)
		if err != nil {
			log.Fatal(err.Error())
		}

		info := struct {
			Token   string    `json:"token"`
			Email   string    `json:"email"`
			Name    string    `json:"name"`
			User_ID string    `json:"userID"`
			Expiry  time.Time `json:"expiry"`
		}{
			Token:   token,
			Email:   result.Email,
			Name:    result.Profile["first_name"] + " " + result.Profile["last_name"],
			User_ID: result.ID,
			Expiry:  parsedTime,
		}

		err = mailing.SendMail(info.Name, info.Email)
		if err != nil {
			log.Fatalln(err.Error())
		}

		if choice == 1 {
			fmt.Println("\n> Successfully Logged in as ", info.Name, " (", info.Email, ")")
			fmt.Println("")
		} else if choice == 2 {
			if usertype == "Old User" {
				fmt.Println("\n> Successfully Logged in as ", info.Name, " (", info.Email, ")")
				fmt.Println("")
			} else if usertype == "New User" {
				fmt.Println("\n> Successfully Created and Logged in as ", info.Name, " (", info.Email, ")")
				fmt.Println("")
			}
		}

		currentUser, err := user.Current()
		if err != nil {
			log.Fatal("Error occured!, try again.")
		}

		path := currentUser.HomeDir + "/Securelee"
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal("Error occured!, try again.")
		}

		tokenPath := currentUser.HomeDir + "/Securelee/token.json"

		data, err := json.Marshal(info)
		if err != nil {
			log.Fatal(err.Error())
		}

		file, err := os.Create(tokenPath)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer file.Close()

		err = ioutil.WriteFile(file.Name(), []byte(data), 0644)
		if err != nil {
			log.Fatal(err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// SignUpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// SignUpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

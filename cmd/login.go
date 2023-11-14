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
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/authn"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// SignUpCmd represents the SignUp command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Securelee Vault.",
	Long:  `Login to Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {

		res := lib.Check()
		if res == true {
			fmt.Print("\033[36m", "\n > Already logged in! Use", "\033[0m")
			fmt.Print("\033[33m", " Securelee-cli whoami", "\033[0m")
			fmt.Println("\033[36m", "command to get the current logged in user.", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}

		var choice int
		var token string
		var err error
		var result *authn.ClientTokenCheckResult
		var ch, usertype string
		fmt.Print("\033[33m", "\n Select any one option: \n", "\033[0m")
		fmt.Print("\033[33m", "\n > 1. Login using Socials through Browser\n", "\033[0m")
		fmt.Print("\033[33m", "> 2. Login using Email and Password on the Terminal\n", "\033[0m")
		fmt.Println("")
		fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
		fmt.Scanf("%d", &choice)
		fmt.Println("")
		if choice == 1 {
			// c := color.New(color.FgCyan, color.Bold)
			// var str any
			fmt.Print("\033[33m", " > Press Enter to Login using Browser", "\033[0m")
			fmt.Print("\033[35m", "  (You can get back to CLI after Successful Authentication) : ", "\033[0m")
			// fmt.Print("")

			_, _ = term.ReadPassword(int(os.Stdin.Fd()))

			// _, key, _ := keyboard.GetSingleKey()
			fmt.Println("")

			// Check if the key is Enter
			// if key == keyboard.KeyEnter {
			lib.Login()
			time.Sleep(25 * time.Second)

			fmt.Print("\033[33m", "\n > Please Enter the Token from the Securelee Authentication Tab: ", "\033[0m")
			fmt.Scan(&token)
			// time.Sleep()
			// }

			if token == "" {
				fmt.Println("\033[31m", " > Invalid Token. Please try again.", "\033[0m")
			}

			result, ch = lib.CheckToken(token)
			if ch != "" {
				fmt.Println("\033[31m", "\n", ch, "\033[0m")
			}

		} else if choice == 2 {
			var email string
			fmt.Print("\033[33m", " > Enter your Email Address : ", "\033[0m")
			fmt.Scan(&email)
			fmt.Println("")
			fmt.Print("\033[33m", " > Enter your Password ", "\033[0m")
			fmt.Print("\033[35m", "\n  { Password must have \n   - at least 8 characters,\n   - at least 1 number characters,\n   - at least 1 special characters } : ", "\033[0m")
			// fmt.Scan(&password)
			password, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println("")
			check := lib.IsValidPassword(string(password))
			if !check {
				fmt.Println("\033[31m", " > Invalid Password.", "\033[0m")
				return
			}

			token, usertype, err = lib.LoginWithEmail(email, string(password))

			if token == "" && usertype == "" && err != nil {
				log.Fatalln("\033[31m", err.Error(), "\033[0m")
			} else if token == "" && usertype != "" && err == nil {
				fmt.Println("\033[31m", "\n > ", usertype, "\033[0m")
				fmt.Println("")
				return
			}

			if token != "" {
				result, ch = lib.CheckToken(token)
				if ch != "" {
					log.Fatal("\033[31m", ch, "\033[0m")
				}
			}

		} else {
			fmt.Println("\033[31m", "\n > Invalid Choice Entered!!, Please try again", "\033[0m")
			fmt.Println("")
			return
		}

		parsedTime, err := time.Parse(time.RFC3339, result.Expire)
		if err != nil {
			log.Fatal("\033[31m", err.Error(), "\033[0m")
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

		// err = mailing.SendMail(info.Name, info.Email)
		// if err != nil {
		// 	log.Fatalln("\033[31m", err.Error(), "\033[0m")
		// }

		if choice == 1 {
			fmt.Println("\033[36m", "\n > Successfully Logged in as ", info.Name, " (", info.Email, ")", "\033[0m")
			fmt.Println("")
		} else if choice == 2 {
			if usertype == "Old User" {
				fmt.Println("\033[36m", "\n > Successfully Logged in as ", info.Name, " (", info.Email, ")", "\033[0m")
				fmt.Println("")
			} else if usertype == "New User" {
				fmt.Println("\033[36m", "\n > Successfully Created and Logged in as ", info.Name, " (", info.Email, ")", "\033[0m")
				fmt.Println("")
			}
		}

		currentUser, err := user.Current()
		if err != nil {
			log.Fatal("\033[31m", " > Error occured!, try again.", "\033[0m")
		}

		path := currentUser.HomeDir + "/Securelee"
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Fatal("\033[31m", "Error occured!, try again.", "\033[0m")
		}

		tokenPath := currentUser.HomeDir + "/Securelee/token.json"

		data, err := json.Marshal(info)
		if err != nil {
			log.Fatal("\033[31m", err.Error(), "\033[0m")
		}

		file, err := os.Create(tokenPath)
		if err != nil {
			fmt.Println("\033[31m", err.Error(), "\033[0m")
			return
		}
		defer file.Close()

		err = ioutil.WriteFile(file.Name(), []byte(data), 0644)
		if err != nil {
			log.Fatal("\033[31m", err.Error(), "\033[0m")
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

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
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// SignUpCmd represents the SignUp command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Securelee Vault.",
	Long:  `Login to Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {

		out := colorable.NewColorableStdout()

		red := "\033[31m"
		cyan := "\033[36m"
		yellow := "\033[33m"
		magenta := "\033[35m"
		reset := "\033[0m"

		res := lib.Check()
		if res == true {
			fmt.Fprintf(out, "\n%s > Already logged in! Use %s", cyan, reset)
			fmt.Fprintf(out, "%s Securelee-cli whoami%s", yellow, reset)
			fmt.Fprintf(out, "%s command to get the current logged in user.%s\n", cyan, reset)
			fmt.Println("")
			os.Exit(0)
		}

		var choice int
		var token string
		var err error
		var result lib.ResponseData
		var ch, usertype string
		fmt.Fprintf(out, "\n%s Select any one option: %s\n", yellow, reset)
		fmt.Fprintf(out, "\n%s > 1. Login using Socials through Browser%s\n", yellow, reset)
		fmt.Fprintf(out, "%s > 2. Login using Email and Password on the Terminal%s\n", yellow, reset)
		fmt.Println("")
		fmt.Fprintf(out, "%s > Enter your choice (for e.g. 1): %s", cyan, reset)
		fmt.Scanf("%d", &choice)
		fmt.Println("")
		if choice == 1 {
			fmt.Fprintf(out, "%s > Press Enter to Login using Browser%s", yellow, reset)
			fmt.Fprintf(out, "%s  (You can get back to CLI after Successful Authentication) : %s", magenta, reset)

			_, _ = term.ReadPassword(int(os.Stdin.Fd()))

			fmt.Println("")

			lib.Login()
			time.Sleep(25 * time.Second)

			fmt.Fprintf(out, "%s\n > Please Enter the Token from the Securelee Authentication Tab: %s", yellow, reset)
			fmt.Scan(&token)

			if token == "" {
				fmt.Fprintf(out, "%s > Invalid Token. Please try again.%s", red, reset)
				return
			}

			result, ch = lib.CheckToken(token)
			if ch != "" {
				fmt.Println("\n", ch)
				return
			}

		} else if choice == 2 {
			var email string
			fmt.Fprintf(out, "%s > Enter your Email Address : %s", yellow, reset)
			fmt.Scan(&email)
			fmt.Println("")
			fmt.Fprintf(out, "%s > Enter your Password %s", yellow, reset)
			fmt.Fprintf(out, "\n%s  { Password must have \n   - at least 8 characters,\n   - at least 1 number characters,\n   - at least 1 special characters } : %s", magenta, reset)

			password, _ := term.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println("")
			check := lib.IsValidPassword(string(password))
			if !check {
				fmt.Fprintf(out, "%s > Invalid Password.%s", red, reset)
				return
			}

			// token, usertype, err = lib.LoginWithEmail(email, string(password))

			if token == "" && usertype == "" && err == nil {
				fmt.Println("\033[31m", "\n > Failed to Login! Please try another login method. \033[0m")
				os.Exit(0)

			} else if token == "" && usertype != "" && err == nil {
				fmt.Println("\033[31m", "\n > ", usertype, "\033[0m")
				fmt.Println("")
				return
			}

			if token != "" {
				result, ch = lib.CheckToken(token)
				if ch != "" {
					log.Fatal(ch)
				}
			}

		} else {
			fmt.Fprintf(out, "%s\n > Invalid Choice Entered!!, Please try again%s", red, reset)
			fmt.Println("")
			return
		}

		parsedTime, err := time.Parse(time.RFC3339, result.Result.Expire)
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
			Email:   result.Result.Email,
			Name:    result.Result.Profile.FirstName + " " + result.Result.Profile.LastName,
			User_ID: result.Result.ID,
			Expiry:  parsedTime,
		}

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
			fmt.Fprintf(out, "%s > Error occured!, try again.%s", red, reset)
			return
		}

		path := currentUser.HomeDir + "/Securelee"
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Fprintf(out, "%s > Error occured!, try again.%s", red, reset)
			return
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

/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/pangeacyber/pangea-go/pangea-sdk/v2/service/vault"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a Particular Secret or Key.",
	Long:  `Get a Particular Secret or Key.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("get called")
		res := lib.Check()
		if !res {
			fmt.Print("\033[31m", "\n > No User logged in, You must Login to use Securelee Vault Services.\n", "\033[0m")
			fmt.Print("\033[36m", "\n > Use 'Securelee-cli login' command to complete the Authentication.\n", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}

		user, err := user.Current()
		if err != nil {
			log.Fatalln("\033[31m", err.Error(), "\033[0m")
		}
		filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatalln("\033[31m", err.Error(), "\033[0m")
		}

		var userData struct {
			UserID string `json:"userID"`
		}
		err = json.Unmarshal(jsonData, &userData)
		if err != nil {
			log.Fatalln("\033[31m", err.Error(), "\033[0m")
		}

		var id string
		fmt.Print("\033[33m", "\n > Enter the ID of the Secret or Key you want : ", "\033[0m")
		fmt.Scan(&id)
		fmt.Println("")

		ctx, cancelFn := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancelFn()
		client := lib.InitVault()

		resp, err := client.Get(ctx,
			&vault.GetRequest{
				ID: id,
			})
		if err != nil {

			re := regexp.MustCompile(`\{[^{}]*\}`)
			match := re.Find([]byte(err.Error()))

			if match == nil {
				fmt.Println("\033[31m", "\n > No JSON data found in the error message.\n", "\033[0m")
				return
			}

			var apiError lib.APIError
			err = json.Unmarshal(match, &apiError)
			if err != nil {
				fmt.Println("\033[31m", err.Error(), "\033[0m")
				return
			}

			if apiError.Status == "ValidationError" {
				fmt.Println("\033[31m", "\n > Invalid Item ID. Please try again.\n", "\033[0m")
				return
			}

		}
		folderKey := "/" + userData.UserID + "/keys/"
		folderSecret := "/" + userData.UserID + "/secrets/"

		if resp != nil && *resp.Status == "Success" {
			if resp.Result.Folder != folderKey {
				if resp.Result.Folder != folderSecret {
					fmt.Println("\033[31m", "\n > Invalid Item ID. Please try again.\n", "\033[0m")
					return
				}
			}
			if *resp.Summary == "Key retrieved" || *resp.Summary == "Key pair retrieved" {
				if resp.Result.Type == "symmetric_key" {
					fmt.Println("\033[36m", "\n > Item details : ", "\033[0m")
					fmt.Print("\033[33m", "\n > ID : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.ID, "\033[0m")
					fmt.Print("\033[33m", " > Name : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Name, "\033[0m")
					fmt.Print("\033[33m", " > Type : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Type, "\033[0m")
					fmt.Print("\033[33m", " > Item Version State : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.CurrentVersion.State, "\033[0m")
					fmt.Print("\033[33m", " > Algoritm : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Algorithm, "\033[0m")
					fmt.Println("")
					return
				} else if resp.Result.Type == "asymmetric_key" {
					fmt.Println("\033[36m", "\n > Item details : ", "\033[0m")
					fmt.Print("\033[33m", "\n > ID : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.ID, "\033[0m")
					fmt.Print("\033[33m", " > Name : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Name, "\033[0m")
					fmt.Print("\033[33m", " > Type : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Type, "\033[0m")
					fmt.Print("\033[33m", " > Item Version State : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.CurrentVersion.State, "\033[0m")
					fmt.Print("\033[33m", " > Algoritm : ", "\033[0m")
					fmt.Println("\033[36m", resp.Result.Algorithm, "\033[0m")
					if resp.Result.CurrentVersion.State == "active" {
						fmt.Print("\033[33m", " > The Public Key : ", "\033[0m")
						fmt.Println("\033[36m", *resp.Result.CurrentVersion.PublicKey, "\033[0m")
					}
					fmt.Println("")
					return
				}

			} else if *resp.Summary == "Secret retrieved" {
				fmt.Println("\033[36m", "\n > Item details : ", "\033[0m")
				fmt.Print("\033[33m", "\n > ID : ", "\033[0m")
				fmt.Println("\033[36m", resp.Result.ID, "\033[0m")
				fmt.Print("\033[33m", " > Name : ", "\033[0m")
				fmt.Println("\033[36m", resp.Result.Name, "\033[0m")
				fmt.Print("\033[33m", " > Type : ", "\033[0m")
				fmt.Println("\033[36m", resp.Result.Type, "\033[0m")
				fmt.Print("\033[33m", " > Item Version State : ", "\033[0m")
				fmt.Println("\033[36m", resp.Result.CurrentVersion.State, "\033[0m")
				if resp.Result.CurrentVersion.State == "deactivated" {
					fmt.Print("\033[33m", " > Secret : ", "\033[0m")
					fmt.Println("\033[36m", *resp.Result.CurrentVersion.Secret, "\033[0m")
				}
				fmt.Println("")
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

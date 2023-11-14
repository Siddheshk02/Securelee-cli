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
	"path/filepath"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update or Modify a Secret Message or Key.",
	Long:  `Update or Modify a Secret Message or Key.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("update called")
		res := lib.Check()
		if !res {
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
			log.Fatal("\033[31m", err, "\033[0m")
		}

		var userData struct {
			UserID string `json:"userID"`
		}
		err = json.Unmarshal(jsonData, &userData)
		if err != nil {
			log.Fatal("\033[31m", err, "\033[0m")
		}

		var choice uint
		fmt.Print("\033[33m", "\n Select any one option: \n", "\033[0m")
		fmt.Println("\033[33m", "\n > 1. Update a Secret", "\033[0m")
		fmt.Println("\033[33m", " > 2. Update a Key", "\033[0m")
		fmt.Println("")
		fmt.Print("\033[36m", " > Enter your choice (for e.g. 1): ", "\033[0m")
		fmt.Scanf("%d", &choice)
		fmt.Println("")

		if choice == 1 {
			err := lib.ListSecrets(userData.UserID)
			if err != nil {
				log.Fatal("\033[31m", err, "\033[0m")
			}
			var id string
			fmt.Print("\033[33m", " > Enter the ID of the Secret to be Updated : ", "\033[0m")
			fmt.Scan(&id)
			if id == "" {
				fmt.Println("\033[31m", "\n > ID cannot be empty. Please try again.", "\033[0m")
				os.Exit(0)
			}

			err = lib.Update(id, "Secret")
			if err != nil {
				log.Fatal("\033[31m", err, "\033[0m")
			}

		} else if choice == 2 {
			err := lib.ListKeys(userData.UserID)
			if err != nil {
				log.Fatal("\033[31m", err, "\033[0m")
			}
			var id string
			fmt.Print("\033[33m", " > Enter the ID of the Key to be Updated : ", "\033[0m")
			fmt.Scan(&id)
			if id == "" {
				fmt.Println("\033[31m", "\n > ID cannot be empty. Please try again.", "\033[0m")
				os.Exit(0)
			}

			err = lib.Update(id, "Key")
			if err != nil {
				log.Fatal("\033[31m", err, "\033[0m")
			}

		} else {
			fmt.Println("\033[31m", " > Invalid Choice Entered!!, Please try again", "\033[0m")
			fmt.Println("")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

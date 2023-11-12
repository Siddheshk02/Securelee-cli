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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a Secret Message or Key.",
	Long:  `Delete a Secret Message or Key.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("delete called")
		res := lib.Check()
		if !res {
			fmt.Println("\n > No User logged in, You must Login to use Securelee Vault Services.")
			fmt.Println("")
			os.Exit(0)
		}

		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		filePath := filepath.Join(user.HomeDir, "/securelee/token.json")

		jsonData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		var userData struct {
			UserID string `json:"userID"`
		}
		err = json.Unmarshal(jsonData, &userData)
		if err != nil {
			log.Fatal(err)
		}

		var choice uint
		fmt.Println("\n Select any one option: ")
		fmt.Println("\n> 1. Delete a Secret")
		fmt.Println("> 2. Delete a Key")
		fmt.Println("")
		fmt.Print("> Enter your choice (for e.g. 1): ")
		fmt.Scanf("%d", &choice)
		fmt.Println("")

		if choice == 1 {
			err := lib.ListSecrets(userData.UserID)
			if err != nil {
				log.Fatal(err)
			}
			var id string
			fmt.Print("> Enter the ID of the Secret to be Deleted : ")
			fmt.Scan(&id)
			if id == "" {
				fmt.Println("\n> ID cannot be empty. Please try again.")
				os.Exit(0)
			}

			err = lib.Delete(id)
			if err != nil {
				log.Fatal(err)
			}

		} else if choice == 2 {
			err := lib.ListKeys(userData.UserID)
			if err != nil {
				log.Fatal(err)
			}
			var id string
			fmt.Print("> Enter the ID of the Key to be Deleted : ")
			fmt.Scan(&id)
			if id == "" {
				fmt.Println("\n> ID cannot be empty. Please try again.")
				os.Exit(0)
			}

			err = lib.Delete(id)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("> Invalid Choice Entered!!, Please try again")
			fmt.Println("")
			os.Exit(0)
		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

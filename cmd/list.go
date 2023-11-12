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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Get all Secret Messages or Keys.",
	Long:  `Get all Secret Messages or Keys.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("list called")
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

		var choice int
		fmt.Println("\n Select any one option: ")
		fmt.Println("\n> 1. List all Secrets.")
		fmt.Println("> 2. List all Keys")
		fmt.Println("")
		fmt.Print("> Enter your choice (for e.g. 1): ")
		fmt.Scanf("%d", &choice)
		fmt.Println("")

		if choice == 1 {
			err := lib.ListSecrets(userData.UserID)
			if err != nil {
				log.Fatal(err)
			}

		} else if choice == 2 {
			err := lib.ListKeys(userData.UserID)
			if err != nil {
				log.Fatal(err)
			}

			// for i := 0; i < count; i++ {
			// 	fmt.Println("\n> ", i+1, " id : ", lists[i].ID)
			// 	fmt.Println("    Name : ", lists[i].Name)
			// 	fmt.Println("    Type : ", lists[i].Type)
			// 	fmt.Println("    Purpose : ", lists[i].Purpose)
			// 	fmt.Println("    Algorithm : ", lists[i].Algorithm)
			// }

		} else {
			fmt.Println("> Invalid Choice Entered!!, Please try again")
			fmt.Println("")
			os.Exit(0)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

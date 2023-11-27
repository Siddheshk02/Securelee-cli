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
	"strings"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create and Store a Secret Message or Key.",
	Long:  `Create and Store a Secret Messages or Keys.`,
	Run: func(cmd *cobra.Command, args []string) {
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
			Email  string `json:"email"`
		}
		err = json.Unmarshal(jsonData, &userData)
		if err != nil {
			log.Fatalln("\033[31m", err.Error(), "\033[0m")
		}
		email := strings.Replace(userData.Email, "@", "-", -1)
		Newemail := strings.Replace(email, ".", "-", -1)
		fmt.Println(Newemail)

		err = lib.Create(Newemail)
		if err != nil {
			log.Fatalln("\033[31m", err.Error(), "\033[0m")
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

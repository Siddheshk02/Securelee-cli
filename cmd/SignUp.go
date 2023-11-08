/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"

	"github.com/Siddheshk02/Securelee-cli/lib"
	"github.com/spf13/cobra"
)

// SignUpCmd represents the SignUp command
var SignUpCmd = &cobra.Command{
	Use:   "SignUp",
	Short: "Sign-up to Securelee Vault.",
	Long:  `Sign-up to Securelee Vault.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("SignUp called")
		// var email string
		// var password string
		// reader := bufio.NewReader(os.Stdin)

		// fmt.Print("\n> Enter the Full Name : \n>")
		// fullName, _ := reader.ReadString('\n')

		// names := strings.Split(fullName, " ")
		// var firstName, lastName string

		// if len(names) >= 2 {
		// 	firstName = names[0]
		// 	lastName = names[1]
		// } else {
		// 	fmt.Println("Invalid full name format.")
		// }

		// fmt.Print("\n> Enter the Email : \n> ")
		// fmt.Scanf("%s\n", &email)
		// if len(strings.TrimSpace(email)) == 0 {
		// 	err := fmt.Errorf("Your Email can't be empty %v", email)
		// 	fmt.Println(err.Error())
		// 	os.Exit(1)
		// }

		// //email verification
		// err := mailing.EmailVerify(email)
		// if err != nil {
		// 	log.Fatal(err)
		// }

		// //if err==nil enter password
		// fmt.Print("\n> Please set a Password for this account : \n> ")
		// fmt.Scanf("%s\n", &password)

		//call for Createandlogin user
		// result := lib.CreateAndLogin(firstName, lastName, email, password)
		fmt.Print("\n> Press Enter to Sign up/Login using Browser (You can get back to CLI Terminal after Successful Authentication.) : ")
		fmt.Scanln(" ")
		lib.SignUp()

		var Token string
		fmt.Println("> Please Enter the Token from the Securelee Authentication Tab: ")
		fmt.Scanf("%s", &Token)

		user, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		filePath := filepath.Join(user.HomeDir, "Securelee/token.json")

		// data, err := json.Marshal(fcresp.Result)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		err = ioutil.WriteFile(filePath, []byte(Token), 0644)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println("Create user success. Result: " + pangea.Stringify(result))
	},
}

func init() {
	rootCmd.AddCommand(SignUpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// SignUpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// SignUpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

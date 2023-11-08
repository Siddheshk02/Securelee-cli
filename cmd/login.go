/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/Siddheshk02/Securelee-cli/controller"
	"github.com/gorilla/mux"
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("login called")
		r := mux.NewRouter()

		r.HandleFunc("/login-success", controller.Callback).Methods("GET").Host("pdn-vehpksfu665ae7k5jewmycb4fxqircam.login.aws.us.pangea.cloud").Queries("state", "{state:[a-zA-Z0-9]+}")
		l, err := net.Listen("tcp", ":80")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print("\n> Press Enter to Sign up/Login : ")
		fmt.Scanln(" ")
		open.Start("https://pdn-vehpksfu665ae7k5jewmycb4fxqircam.login.aws.us.pangea.cloud/authorize?redirect_uri=%2Flogin-success&state=xxxxxxxxxxxxx")

		http.Serve(l, r)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

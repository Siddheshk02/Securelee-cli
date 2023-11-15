/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/mbndr/figlet4go"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "Securelee-cli",
	Version: "v1.1.2",
	Short:   "\nA CLI based Vault App for storing your Secret Messages or Keys Securely.",
	Long:    `A CLI based Vault App for storing your Secret Messages or Keys Securely.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		ascii := figlet4go.NewAsciiRender()
		options := figlet4go.NewRenderOptions()
		options.FontColor = []figlet4go.Color{
			figlet4go.ColorCyan,
		}

		renderStr, _ := ascii.RenderOpts("Securelee.", options)
		fmt.Print(renderStr)
		fmt.Print(" > v1.1.2")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.Securelee-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

package cmd

import (
	"fmt"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/config"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/logger"
	"os"

	"github.com/spf13/cobra"
	_ "github.com/spf13/viper/remote"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-starter",
	Short: "A brief description of your application",
}

func Execute() {
	cfg := config.LoadConfig()
	fmt.Println(cfg)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	logger.Info("about to start the application")
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(migrationCmd)
}

package cmd

import (
	"github.com/monstar-lab-bd/golang-starter-rest-api/domain"
	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/config"
	"log"

	"github.com/monstar-lab-bd/golang-starter-rest-api/internal/conn"
	"github.com/spf13/cobra"
)

var migrationCmd = &cobra.Command{
	Use:   "migrate",
	Short: "migration command apply the db migration",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if err := conn.ConnectDb(config.Db()); err != nil {
			log.Println(err)
			return err
		}
		return nil
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := conn.Db().AutoMigrate(
			domain.User{},
		); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(migrationCmd)
}

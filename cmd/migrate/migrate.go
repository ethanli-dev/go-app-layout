/*
Copyright © 2025 lixw
*/
package migrate

import (
	"fmt"

	"github.com/ethanli-dev/go-app-layout/internal/model"
	"github.com/ethanli-dev/go-app-layout/pkg/config"
	"github.com/ethanli-dev/go-app-layout/pkg/database"
	"github.com/spf13/cobra"
)

// StartCmd represents the server command
var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "migrate",
		Aliases: []string{"migrate"},
		Short:   "Run database migration",
		Example: "go-app-layout migrate -c config/dev.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(fmt.Sprintf("starting migrate with config: %s", configYml))
			cfg, err := config.New(configYml)
			if err != nil {
				return err
			}
			db, err := database.New(database.WithUrl(cfg.Database.Url))
			if err != nil {
				return err
			}
			return db.AutoMigrate(
				&model.Tenant{},
			)
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/dev.yml", "Start server with configuration file")
}

/*
Copyright © 2025 lixw
*/
package server

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

// StartCmd represents the server command
var (
	configYml string
	StartCmd  = &cobra.Command{
		Use:     "server",
		Aliases: []string{"start"},
		Short:   "Start the server",
		Example: "go-app-layout server -c config/dev.yml",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(fmt.Sprintf("starting server with config: %s", configYml))
			app, err := CreateApp(configYml)
			if err != nil {
				return err
			}
			return app.Run(context.Background())
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/dev.yml", "Start server with configuration file")
}

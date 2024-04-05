package cmd

import (
	"github.com/spf13/cobra"

	"github.com/ignite/gex/pkg/xurl"
	"github.com/ignite/gex/services/explorer"
)

const defaultHost = "http://localhost:26657"

// NewExplorer creates a new explorer command.
func NewExplorer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "explorer [host]",
		Short: "Gex is a cosmos explorer for terminals",
		Long:  "Gex is a tool for generate block explorer for blockchains built with Cosmos SDK.", Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			host := defaultHost
			if len(args) > 0 && args[0] != "" {
				host = args[0]
			}

			hostURL, err := xurl.Parse(host)
			if err != nil {
				return err
			}

			return explorer.Run(cmd.Context(), hostURL.String())
		},
	}

	return cmd
}

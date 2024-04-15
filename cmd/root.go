package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/ignite/gex/version"
)

const checkVersionTimeout = time.Millisecond * 600

// NewRootCmd creates a new root command for `Gex` with its sub commands.
// Returns the cobra.Command, a cleanup function and an error. The cleanup
// function must be invoked by the caller to clean eventual Ignite App instances.
func NewRootCmd() *cobra.Command {
	c := &cobra.Command{
		Use:           "gex",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
			if cmd.Use != "completion" && cmd.Use != "explorer" {
				checkNewVersion(cmd.Context())
			}

			return nil
		},
	}

	c.AddCommand(
		NewExplorer(),
		NewVersion(),
	)

	return c
}

func checkNewVersion(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, checkVersionTimeout)
	defer cancel()

	isAvailable, next, err := version.CheckNext(ctx)
	if err != nil || !isAvailable {
		return
	}

	fmt.Printf("⬆️ Gex %[1]v is available! To upgrade: https://github.com/ignite/gex/releases/%[1]v", next)
}

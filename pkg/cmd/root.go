package cmd

import (
	"github.com/jzelinskie/cobrautil/v2/cobraotel"
	"github.com/jzelinskie/cobrautil/v2/cobrazerolog"
	"github.com/spf13/cobra"

	"github.com/zapravila/spicedb/pkg/cmd/server"
	"github.com/zapravila/spicedb/pkg/cmd/termination"
	"github.com/zapravila/spicedb/pkg/releases"
	"github.com/zapravila/spicedb/pkg/runtime"
)

func RegisterRootFlags(cmd *cobra.Command) {
	cobrazerolog.New().RegisterFlags(cmd.PersistentFlags())
	cobraotel.New(cmd.Use).RegisterFlags(cmd.PersistentFlags())
	releases.RegisterFlags(cmd.PersistentFlags())
	termination.RegisterFlags(cmd.PersistentFlags())
	runtime.RegisterFlags(cmd.PersistentFlags())
}

func NewRootCommand(programName string) *cobra.Command {
	return &cobra.Command{
		Use:           programName,
		Short:         "A modern permissions database",
		Long:          "A database that stores, computes, and validates application permissions",
		Example:       server.ServeExample(programName),
		SilenceErrors: true,
		SilenceUsage:  true,
	}
}

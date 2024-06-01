package cmd

import (
	"fmt"

	"github.com/jzelinskie/cobrautil/v2"
	"github.com/jzelinskie/cobrautil/v2/cobraotel"
	"github.com/jzelinskie/cobrautil/v2/cobrazerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/zapravila/spicedb/pkg/cmd/server"
	"github.com/zapravila/spicedb/pkg/cmd/termination"
	"github.com/zapravila/spicedb/pkg/releases"
	"github.com/zapravila/spicedb/pkg/runtime"
)

func RegisterRootFlags(cmd *cobra.Command) error {
	cobrazerolog.New().RegisterFlags(cmd.PersistentFlags())
	cobraotel.New(cmd.Use).RegisterFlags(cmd.PersistentFlags())
	releases.RegisterFlags(cmd.PersistentFlags())
	termination.RegisterFlags(cmd.PersistentFlags())
	runtime.RegisterFlags(cmd.PersistentFlags())


	zl := cobrazerolog.New()
	zl.RegisterFlags(cmd.PersistentFlags())
	if err := zl.RegisterFlagCompletion(cmd); err != nil {
		return fmt.Errorf("failed to register zerolog flag completion: %w", err)
	}

	ot := cobraotel.New(cmd.Use)
	ot.RegisterFlags(cmd.PersistentFlags())
	if err := ot.RegisterFlagCompletion(cmd); err != nil {
		return fmt.Errorf("failed to register otel flag completion: %w", err)
	}

	releases.RegisterFlags(cmd.PersistentFlags())
	termination.RegisterFlags(cmd.PersistentFlags())
	runtime.RegisterFlags(cmd.PersistentFlags())



	return nil
}

// DeprecatedRunE wraps the RunFunc with a warning log statement.
func DeprecatedRunE(fn cobrautil.CobraRunFunc, newCmd string) cobrautil.CobraRunFunc {
	return func(cmd *cobra.Command, args []string) error {
		log.Warn().Str("newCommand", newCmd).Msg("use of deprecated command")
		return fn(cmd, args)
	}
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

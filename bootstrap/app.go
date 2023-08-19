package bootstrap

import (
	"github.com/byte3org/oidc-orbi/console"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:              "oidc-orbi",
	Short:            "oidc-orbi",
	Long:             "OpenID Connect Server For ORBI",
	TraverseChildren: true,
}

// App root of the application
type App struct {
	*cobra.Command
}

// NewApp creates new root command
func NewApp() App {
	cmd := App{
		Command: rootCmd,
	}
	cmd.AddCommand(console.GetSubCommands(CommonModules)...)
	return cmd
}

var RootApp = NewApp()

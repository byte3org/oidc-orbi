package console

import (
	"context"
	"log"

	"github.com/byte3org/oidc-orbi/internal/lib"
	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

var cmds = map[string]lib.Command{
	"app:serve": NewServeCommand(),
}

// GetSubCommands gives a list of sub commands
func GetSubCommands(opt fx.Option) []*cobra.Command {
	subCommands := make([]*cobra.Command, 0)
	for name, cmd := range cmds {
		subCommands = append(subCommands, WrapSubCommand(name, cmd, opt))
	}
	return subCommands
}

func WrapSubCommand(name string, cmd lib.Command, opt fx.Option) *cobra.Command {
	wrappedCmd := &cobra.Command{
		Use:   name,
		Short: cmd.Short(),
		Run: func(c *cobra.Command, args []string) {
			opts := fx.Options(
				fx.Invoke(cmd.Run()),
			)
			ctx := context.Background()
			app := fx.New(opt, opts)
			err := app.Start(ctx)

			defer func() {
				err = app.Stop(ctx)
				if err != nil {
					log.Fatal(err)
				}
			}()

			if err != nil {
				panic(err)
			}
		},
	}
	cmd.Setup(wrappedCmd)
	return wrappedCmd
}

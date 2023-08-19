package console

import (
	"fmt"
	"github.com/byte3org/oidc-orbi/internal/infrastructure"
	"github.com/byte3org/oidc-orbi/internal/lib"
	"github.com/byte3org/oidc-orbi/internal/repository"
	"github.com/byte3org/oidc-orbi/internal/storage"
	"github.com/byte3org/oidc-orbi/internal/ui"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

// ServeCommand test command
type ServeCommand struct{}

func (s *ServeCommand) Short() string {
	return "serve application"
}

func (s *ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		env *lib.Env,
		database infrastructure.Database,
		userRepo repository.UsersRepository,
	) {
		port := "9998"
		issuer := fmt.Sprintf("http://localhost:%s/", port)

		storage := storage.NewStorage(storage.NewUserStore(issuer, userRepo))

		router := ui.SetupServer(issuer, storage)

		server := &http.Server{
			Addr:    ":" + port,
			Handler: router,
		}

		log.Printf("server listening on http://localhost:%s/", port)
		log.Println("press ctrl+c to stop")

		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}

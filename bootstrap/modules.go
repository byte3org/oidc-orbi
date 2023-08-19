package bootstrap

import (
	"github.com/byte3org/oidc-orbi/internal/infrastructure"
	"github.com/byte3org/oidc-orbi/internal/lib"
	"github.com/byte3org/oidc-orbi/internal/repository"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	infrastructure.Module,
	repository.Module,
    lib.Module, 
)

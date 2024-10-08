package cmd

import (
	"github.com/rafaalrazzak/e-campus-be/internal/http"
	"github.com/rafaalrazzak/e-campus-be/internal/routes"
	"github.com/rafaalrazzak/e-campus-be/pkg/framework/config"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/database"
	"github.com/rafaalrazzak/e-campus-be/pkg/services/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Providers() fx.Option {
	return fx.Provide(
		zap.NewDevelopment,
		config.NewConfig,
		database.NewDatabaseConn,
		database.NewECampusDBImpl,
		redis.NewRedisImpl,
		http.NewFiberApp,
		routes.UserRoutes,
		routes.AuthRoutes,
	)
}

func Entrypoint() fx.Option {
	return fx.Invoke(
		database.Migrator,
		http.ServeHTTP,
	)
}

func Run() {
	fx.New(Providers(), Entrypoint()).Run()
}

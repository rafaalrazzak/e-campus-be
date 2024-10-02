package bunapp

import (
	"context"
	"database/sql"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bunrouter"
	"github.com/urfave/cli/v2"
)

type appCtxKey struct{}

// AppFromContext retrieves the App instance from the context.
func AppFromContext(ctx context.Context) *App {
	return ctx.Value(appCtxKey{}).(*App)
}

// ContextWithApp sets the App instance in the context.
func ContextWithApp(ctx context.Context, app *App) context.Context {
	ctx = context.WithValue(ctx, appCtxKey{}, app)
	return ctx
}

// App represents the application context.
type App struct {
	ctx context.Context
	cfg *AppConfig

	stopping uint32
	stopCh   chan struct{}

	onStop      appHooks
	onAfterStop appHooks

	router    *bunrouter.Router
	apiRouter *bunrouter.Group

	// lazy init
	dbOnce sync.Once
	db     *bun.DB
}

// New creates a new App instance.
func New(ctx context.Context, cfg *AppConfig) *App {
	app := &App{
		cfg:    cfg,
		stopCh: make(chan struct{}),
	}
	app.ctx = ContextWithApp(ctx, app)
	app.initRouter()
	return app
}

// StartCLI initializes the app context for CLI commands.
func StartCLI(c *cli.Context) (context.Context, *App, error) {
	return Start(c.Context, c.Command.Name, c.String("env"))
}

// Start initializes the app with the given service and environment.
func Start(ctx context.Context, service, envName string) (context.Context, *App, error) {
	cfg, err := LoadConfig(service, envName)
	if err != nil {
		return nil, nil, err
	}
	return StartConfig(ctx, cfg)
}

// StartConfig initializes the app with the provided configuration.
func StartConfig(ctx context.Context, cfg *AppConfig) (context.Context, *App, error) {
	rand.Seed(time.Now().UnixNano())

	app := New(ctx, cfg)
	if err := onStart.Run(ctx, app); err != nil {
		return nil, nil, err
	}
	return app.ctx, app, nil
}

// Stop stops the application.
func (app *App) Stop() {
	_ = app.onStop.Run(app.ctx, app)
	_ = app.onAfterStop.Run(app.ctx, app)
}

// OnStop adds a hook to be executed on stopping.
func (app *App) OnStop(name string, fn HookFunc) {
	app.onStop.Add(newHook(name, fn))
}

// OnAfterStop adds a hook to be executed after stopping.
func (app *App) OnAfterStop(name string, fn HookFunc) {
	app.onAfterStop.Add(newHook(name, fn))
}

// Context returns the app's context.
func (app *App) Context() context.Context {
	return app.ctx
}

// Config returns the app's configuration.
func (app *App) Config() *AppConfig {
	return app.cfg
}

// Running checks if the application is running.
func (app *App) Running() bool {
	return !app.Stopping()
}

// Stopping checks if the application is stopping.
func (app *App) Stopping() bool {
	return atomic.LoadUint32(&app.stopping) == 1
}

// IsDebug checks if the application is in debug mode.
func (app *App) IsDebug() bool {
	return app.cfg.Debug
}

// Router returns the app's router.
func (app *App) Router() *bunrouter.Router {
	return app.router
}

// APIRouter returns the app's API router.
func (app *App) APIRouter() *bunrouter.Group {
	return app.apiRouter
}

// DB initializes and returns the database connection.
func (app *App) DB() *bun.DB {
	app.dbOnce.Do(func() {
		dsn := app.cfg.DB.DSN
		sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

		db := bun.NewDB(sqldb, pgdialect.New())
		db.AddQueryHook(bundebug.NewQueryHook(
			bundebug.WithEnabled(app.IsDebug()),
			bundebug.FromEnv(""),
		))

		app.OnStop("db.Close", func(ctx context.Context, _ *App) error {
			return db.Close()
		})

		app.db = db
	})
	return app.db
}

// WaitExitSignal listens for termination signals and returns the first received signal.
func WaitExitSignal() os.Signal {
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)
	return <-ch
}

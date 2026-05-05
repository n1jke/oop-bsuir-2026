package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/config"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/application"
	bookhttp "github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/codegen"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/http/server"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/internal/infrastructure/repository"
	"github.com/n1jke/oop-bsuir-2026/laboratory_work-7/pkg"
)

func NewApp() fx.Option {
	return fx.Options(
		fx.Provide(
			ProvideConfig,
			ProvideLogger,
			ProvideDBPool,
			ProvideTransactor,
			ProvideRepositories,
			ProvideServices,
			ProvideServer,
		),
		fx.WithLogger(func(log *slog.Logger) fxevent.Logger {
			return &fxevent.SlogLogger{Logger: log}
		}),
		fx.Invoke(RegisterDBLifecycle),
		fx.Invoke(RegisterServiceLifecycle),
	)
}

func ProvideConfig() (*config.AppConfig, error) {
	if err := config.LoadEnv(); err != nil {
		return nil, err
	}

	return config.LoadConfig()
}

func ProvideLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func ProvideDBPool(cfg *config.AppConfig) (*pgxpool.Pool, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.DB.ConnectionString())
	if err != nil {
		return nil, err
	}

	connConfig.MaxConns = 3
	connConfig.MinConns = 1
	connConfig.MaxConnIdleTime = 100 * time.Millisecond
	connConfig.MaxConnLifetime = time.Second

	pool, err := pgxpool.NewWithConfig(context.Background(), connConfig)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func ProvideTransactor(pool *pgxpool.Pool) application.Transactor {
	return repository.NewTxChain(pool)
}

func ProvideRepositories(pool *pgxpool.Pool) (
	application.BookRepository, application.ExchangeRepository, application.OwnedBookRepository,
	application.ReviewRepository, application.UserRepository,
) {
	return repository.NewBookRepoSQL(pool), repository.NewExchangeRepoSQL(pool),
		repository.NewOwnedBookRepoSQL(pool), repository.NewReviewRepoSQL(pool),
		repository.NewUserRepoSQL(pool)
}

//nolint:gocritic // di container
func ProvideServices(logger *slog.Logger, tx application.Transactor, cfg *config.AppConfig,
	bookRepo application.BookRepository, exchangeRepo application.ExchangeRepository,
	ownedBookRepo application.OwnedBookRepository, reviewRepo application.ReviewRepository,
	userRepo application.UserRepository,
) (
	*application.AuthService, *application.UserService, *application.BookService,
	*application.LibraryService, *application.ReviewService, *application.ExchangeService,
) {
	return application.NewAuthService(logger, userRepo, tx, cfg.SVC.Key, 24*time.Hour),
		application.NewUserService(logger, userRepo),
		application.NewBookService(logger, bookRepo, tx),
		application.NewLibraryService(logger, ownedBookRepo, bookRepo, userRepo, tx),
		application.NewReviewService(logger, reviewRepo, bookRepo, userRepo, tx),
		application.NewExchangeService(logger, exchangeRepo, ownedBookRepo, userRepo, tx)
}

func ProvideServer(cfg *config.AppConfig, logger *slog.Logger, authService *application.AuthService,
	userService *application.UserService, bookService *application.BookService, libraryService *application.LibraryService,
	reviewService *application.ReviewService, exchangeService *application.ExchangeService,
) (*pkg.HTTPServer, error) {
	addr, err := pkg.ResolveAddr(cfg.SVC.Addr)
	if err != nil {
		return nil, err
	}

	svr := server.NewHandler(authService, userService, bookService, libraryService, reviewService, exchangeService)
	handler := bookhttp.HandlerWithOptions(svr, bookhttp.StdHTTPServerOptions{
		Middlewares: []bookhttp.MiddlewareFunc{
			bookhttp.MiddlewareFunc(server.AuthMiddleware(logger, []byte(cfg.SVC.Key))),
		},
	})

	return pkg.NewHTTPServer(addr, handler,
		cfg.SVC.HTTPRead, cfg.SVC.HTTPWrite, cfg.SVC.HTTPIdle), nil
}

func RegisterDBLifecycle(lc fx.Lifecycle, pool *pgxpool.Pool) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return nil
		},
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})
}

func RegisterServiceLifecycle(lc fx.Lifecycle, logger *slog.Logger, cfg *config.AppConfig, svr *pkg.HTTPServer) {
	var cancel context.CancelFunc

	done := make(chan struct{})

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			var ctx context.Context

			ctx, cancel = context.WithCancel(context.Background())

			go func() {
				defer close(done)

				if err := svr.Start(ctx); err != nil {
					logger.Error("book svr crashed", slog.Any("err", err))
					cancel()
				}
			}()

			logger.Info("book svr started")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("book stopping")

			ctxShutdown, cancelShutdown := context.WithTimeout(ctx, cfg.SVC.Shutdown)
			defer cancelShutdown()

			if err := svr.Stop(ctxShutdown); err != nil {
				logger.Error("graceful stop svr", slog.Any("err", err))
			}

			cancel()

			select {
			case <-done:
				logger.Info("book stopped")
			case <-ctxShutdown.Done():
				logger.Warn("book stop deadline exceeded", slog.Any("err", ctx.Err()))
			}

			return nil
		},
	})
}

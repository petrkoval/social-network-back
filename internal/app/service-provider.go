package app

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/petrkoval/social-network-back/internal/config"
	"github.com/petrkoval/social-network-back/internal/logger"
	"github.com/petrkoval/social-network-back/internal/services"
	"github.com/petrkoval/social-network-back/internal/storage"
	"github.com/petrkoval/social-network-back/internal/transport/http"
	"github.com/petrkoval/social-network-back/internal/transport/http/handlers/auth"
	"github.com/petrkoval/social-network-back/pkg/db/postgres"
	"github.com/rs/zerolog"
)

type ServiceProvider struct {
	cfg         *config.Config
	logger      *zerolog.Logger
	dbClient    *pgxpool.Pool
	router      *http.Router
	authHandler *auth.Handler
}

func NewServiceProvider() *ServiceProvider {
	return &ServiceProvider{}
}

func (sp *ServiceProvider) Init() {
	sp.initLogger()
	sp.initConfig()
	sp.initDbClient()
	sp.initRouter()
	sp.initHandlers()
}

func (sp *ServiceProvider) StartServer() {
	sp.logger.Debug().Msg("starting server")

	sp.logger.Fatal().Err(sp.router.Start()).Msg("failed to start server")
}

func (sp *ServiceProvider) initLogger() {
	if sp.logger == nil {
		sp.logger = logger.NewLogger()
		sp.logger.Debug().Msg("initialized logger")
	}
}

func (sp *ServiceProvider) initConfig() {
	var err error
	sp.logger.Debug().Msg("initializing config")

	if sp.cfg == nil {
		sp.cfg, err = config.MustLoad()
		if err != nil {
			sp.logger.Fatal().Err(err).Msg("failed to init config")
		}
	}
}

func (sp *ServiceProvider) initDbClient() {
	var err error
	sp.logger.Debug().Msg("initializing db client")

	if sp.dbClient == nil {
		sp.dbClient, err = postgres.NewPostgreSQLClient(sp.cfg.Database, sp.logger)
		if err != nil {
			sp.logger.Fatal().Err(err).Msg("failed to init database client")
		}
	}
}

func (sp *ServiceProvider) initRouter() {
	sp.logger.Debug().Msg("initializing router")

	if sp.router == nil {
		sp.router = http.NewRouter(sp.cfg.Server)
		sp.router.InitMiddlewares()
	}
}

func (sp *ServiceProvider) initHandlers() {
	sp.logger.Debug().Msg("initializing handlers")

	s := sp.newAuthService()

	authHandler := auth.NewAuthHandler(s, sp.logger)

	authHandler.MountOn(sp.router)
}

func (sp *ServiceProvider) newAuthService() *services.AuthService {
	sp.logger.Debug().Msg("creating auth service")

	userStorage := storage.NewUserStorage(sp.dbClient)
	tokenStorage := storage.NewTokenStorage(sp.dbClient)

	userService := services.NewUserService(userStorage, sp.logger)
	tokenService := services.NewTokenService(tokenStorage, sp.logger, sp.cfg.Tokens)

	return services.NewAuthService(tokenService, userService)
}

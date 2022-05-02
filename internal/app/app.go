package app

import (
	"context"
	conf "github.com/da-semenov/gophermart/internal/app/config"
	"github.com/da-semenov/gophermart/internal/app/handlers"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/client"
	"github.com/da-semenov/gophermart/internal/app/infrastructure/datastore"
	"github.com/da-semenov/gophermart/internal/app/repository"
	"github.com/da-semenov/gophermart/internal/app/service"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func RunApp() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()

	config := conf.NewConfig()
	err = config.Init()
	if err != nil {
		logger.Fatal("can't init configuration", zap.Error(err))
	}

	postgresHandlerTx, err := datastore.NewPostgresHandlerTX(context.Background(), config.DatabaseDSN, logger)
	if err != nil {
		logger.Fatal("can't init postgres handler", zap.Error(err))
	}

	if config.ReInit {
		err = repository.ClearDatabase(context.Background(), postgresHandlerTx)
		if err != nil {
			logger.Fatal("can't clear dbqueries structure", zap.Error(err))
			return
		}
	}
	err = repository.InitDatabase(context.Background(), postgresHandlerTx)
	if err != nil {
		logger.Fatal("can't init dbqueries structure", zap.Error(err))
		return
	}

	userRepository, err := repository.NewUserRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init user repository", zap.Error(err))
		return
	}

	orderRepository, err := repository.NewOrderRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init order repository", zap.Error(err))
		return
	}

	balanceRepository, err := repository.NewBalanceRepository(postgresHandlerTx, logger)
	if err != nil {
		logger.Fatal("can't init balance repository", zap.Error(err))
		return
	}

	authService := service.NewAuthService(userRepository, logger)
	orderService := service.NewOrderService(orderRepository, logger, config.ValidateOrderNum)
	balanceService := service.NewBalanceService(balanceRepository, logger)
	auth := handlers.NewAuth("secret")
	authHandler := handlers.NewAuthHandler(authService, auth, logger)
	orderHandler := handlers.NewOrderHandler(orderService, auth, logger)
	balanceHandler := handlers.NewBalanceHandler(balanceService, auth, logger)

	accrualClient := client.NewAccrualClient(config.AccrualSystemAddress, logger)
	gophermartClient := client.NewGophermartClient(config.ServerAddress, logger)
	accrualService := service.NewAccrualService(orderRepository, balanceRepository, accrualClient, gophermartClient, logger, config.EnableAccrual)
	accrualHandler := handlers.NewAccrualHandler(accrualService, logger)

	router := chi.NewRouter()
	publicRoutes(router, authHandler, accrualHandler, postgresHandlerTx, logger)
	protectedOrderRoutes(router, auth.GetJWTAuth(), postgresHandlerTx, orderHandler, logger)
	protectedBalanceRoutes(router, auth.GetJWTAuth(), postgresHandlerTx, balanceHandler, logger)

	go accrualService.StartProcessJob(1)
	log.Println("starting server on 8080...")
	log.Fatal(http.ListenAndServe(config.ServerAddress, router))
}

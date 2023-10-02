package main

import (
	"gitlab.com/lokalpay-dev/digital-goods/config"
	"gitlab.com/lokalpay-dev/digital-goods/internal/adapter"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/notifier"
	"gitlab.com/lokalpay-dev/digital-goods/internal/pkg/slog"
	"gitlab.com/lokalpay-dev/digital-goods/internal/repository"
	"gitlab.com/lokalpay-dev/digital-goods/internal/server/http"
	"gitlab.com/lokalpay-dev/digital-goods/internal/server/http/controller"
	"gitlab.com/lokalpay-dev/digital-goods/internal/service"
	"go.uber.org/zap"
)

func main() {
	slog.NewLogger(slog.Info)
	// read from env
	envConfig, err := config.Reader()
	if err != nil {
		slog.Fatalw("failed to read config file", zap.Error(err))
	}

	// bind env to schema
	cfg := config.BindConfig(envConfig)

	// adapter injector
	adptr := adapter.New(cfg.Provider)

	// init notifier
	notifier.New()

	// init repository
	repo := repository.New(cfg.Storage)

	// init service/use-case/business logic
	svc := service.New(
		repo,
		cfg.App,
		adptr.LFI,
		adptr.Prismalink,
	)

	// http server will be used only for callback operation
	httpController := controller.NewController(
		cfg,
		svc.Product,
		svc.Payment,
		svc.Order,
		svc.Fulfillment,
		svc.Customer,
	)
	httpServer := http.NewHttpServer(cfg.HTTPServer, httpController)
	httpServer.ListenAndServe()
}

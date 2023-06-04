package app

import (
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"synonym/pkg/config"
	"synonym/pkg/http/router"
	"synonym/pkg/http/server"
	"synonym/pkg/reporters"
	"synonym/pkg/thesaurus"
)

func initHTTPServer(configFile string) {
	cfg := config.NewConfig(configFile)
	logger := initLogger(cfg)
	rt := initRouter(cfg, logger)

	server.NewServer(cfg, logger, rt).Start()
}

func initRouter(cfg config.Config, logger *zap.Logger) http.Handler {
	thesaurusSvc := initService()

	return router.NewRouter(logger, thesaurusSvc)
}

func initService() thesaurus.Service {
	return thesaurus.NewThesaurus()
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.GetLogConfig().GetLevel(),
		getWriters(cfg.GetLogFileConfig())...,
	)
}

func getWriters(cfg config.LogFileConfig) []io.Writer {
	return []io.Writer{
		os.Stdout,
		reporters.NewExternalLogFile(cfg),
	}
}

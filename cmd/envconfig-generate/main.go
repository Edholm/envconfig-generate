package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"edholm.dev/envconfig-generate/internal/setup"
	"edholm.dev/envconfig-generate/internal/tagparser"
	"edholm.dev/go-logging"
)

func main() {
	ctx, cancelFunc := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	logger := logging.NewLoggerFromEnv()
	ctx = logging.WithLogger(ctx, logger)
	defer cancelFunc()

	providedFiles := os.Args[1:]
	if len(providedFiles) == 0 {
		logger.Info("you need to supply files to parse")
		os.Exit(1)
	}

	if err := realMain(ctx, providedFiles); err != nil {
		logger.Warnw("envconfig-generate failed", "err", err)
	}
}

func realMain(ctx context.Context, providedFiles []string) error {
	logger := logging.FromContext(ctx)

	asts := setup.ParseAst(ctx, providedFiles)
	availableConfigs := tagparser.Analyze(ctx, asts)

	for _, config := range availableConfigs {
		logger.Infof("%s/%s", config.Package, config.Name)
		for _, opt := range config.Options {
			logger.Infof("%s", opt.String())
		}
	}

	logger.Infow("parsed files", "fileCount", len(asts))

	return nil
}

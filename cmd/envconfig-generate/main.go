package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"edholm.dev/envconfig-generate/internal/output"
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
	asts := setup.ParseAst(ctx, providedFiles)
	availableConfigs := tagparser.Analyze(ctx, asts)

	md, err := output.ToMarkdown(ctx, availableConfigs)
	if err != nil {
		return err
	}

	fmt.Printf("%s", md)
	return nil
}

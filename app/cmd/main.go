package main

import (
	"context"
	"os"

	"github.com/theartofdevel/logging"
	"github.com/theartofdevel/tele-openai-template/internal/api/bot"
	"github.com/theartofdevel/tele-openai-template/internal/config"
	"github.com/theartofdevel/tele-openai-template/internal/service/openai"
	"github.com/theartofdevel/tele-openai-template/internal/service/replicate"
)

func main() {
	ctx := context.Background()

	cfg := config.GetConfig()

	logging.WithAttrs(
		ctx,
		cfg.LoggingAttrs()...,
	).Info("application initializing with configuration")

	logger := logging.NewLogger(
		logging.WithLevel(cfg.App.LogLevel),
		logging.WithIsJSON(cfg.App.IsLogJSON),
	)

	ctx = logging.ContextWithLogger(ctx, logger)

	openaiCfg := openai.NewConfig(cfg.OpenAI.ApiKey)

	openaiService, err := openai.NewService(openaiCfg)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to create openai service")
		os.Exit(1)
	}

	replicateCfg := replicate.NewConfig(cfg.Replicate.Token)

	replicateService, err := replicate.NewService(replicateCfg)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to create replicate service")
		os.Exit(1)
	}

	botCfg := bot.NewConfig(cfg.Bot.Token, cfg.Bot.Timeout, cfg.Bot.Whitelist)

	botWrapper, err := bot.NewWrapper(
		botCfg,
		openaiService,
		replicateService,
	)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("failed to create bot wrapper")
		os.Exit(1)
	}

	logging.L(ctx).Info("application started!")

	err = botWrapper.Start(ctx)
	if err != nil {
		logging.WithAttrs(ctx, logging.ErrAttr(err)).Error("bot stopped")
	}

	logging.L(ctx).Info("application stopped")

	os.Exit(0)
}

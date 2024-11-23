package bot

import (
	"context"
	"strings"
	"time"

	"github.com/theartofdevel/logging"
	"github.com/theartofdevel/tele-openai-template/internal/service/replicate"
	tb "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

var (
	// Universal markup builders.
	menu = &tb.ReplyMarkup{ResizeKeyboard: true}

	// Reply buttons.
	btnNewConversation = menu.Text("⚙ New conversation")
	btnGenerateImage   = menu.Text("🎨 Generate image")
)

type openAiService interface {
	ChatCompletion(context.Context, int64, string) (string, error)
	GenerateImagePrompt(ctx context.Context, prompt string) (string, error)

	NewConversation(_ context.Context, id int64)
}

type replicateService interface {
	GenerateImage(ctx context.Context, reqGen *replicate.Request) (replicate.Response, error)
}

type Wrapper struct {
	bot    *tb.Bot
	config *Config

	states UserStates

	openaiSvc    openAiService
	replicateSvc replicateService
}

func NewWrapper(config *Config, openaiSvc openAiService, replicateSvc replicateService) (*Wrapper, error) {
	err := config.Validate()
	if err != nil {
		return nil, err
	}

	settings := tb.Settings{
		Token:  config.Token,
		Poller: &tb.LongPoller{Timeout: config.Timeout},
	}

	bot, err := tb.NewBot(settings)
	if err != nil {
		return nil, err
	}

	bot.Use(middleware.Logger())
	bot.Use(middleware.AutoRespond())
	bot.Use(middleware.Whitelist(config.Whitelist...))

	w := &Wrapper{
		bot:    bot,
		config: config,

		states: NewUserStates(),

		openaiSvc:    openaiSvc,
		replicateSvc: replicateSvc,
	}

	w.prepare()

	return w, nil
}

func (w *Wrapper) Start(_ context.Context) error {
	w.bot.Start()

	return nil
}

func (w *Wrapper) prepare() {
	menu.Reply(
		menu.Row(btnNewConversation),
		menu.Row(btnGenerateImage),
	)

	w.bot.Handle("/start", w.startHandler)
	w.bot.Handle(&btnNewConversation, w.newConversationHandler)
	w.bot.Handle(&btnGenerateImage, w.generateImageHandler)
	w.bot.Handle(tb.OnText, w.onTextHandler)
}

func (w *Wrapper) startHandler(c tb.Context) error {
	ctx := context.TODO()

	id := c.Sender().ID

	logging.L(ctx).Info("user start bot", logging.Int64Attr("user_id", id))

	w.openaiSvc.NewConversation(ctx, id)

	return c.Send("You are ready to start new conversation, just enter what do you want!", menu)
}

func (w *Wrapper) newConversationHandler(c tb.Context) error {
	ctx := context.TODO()

	w.states.Set(c.Sender().ID, ChatDoing)

	w.openaiSvc.NewConversation(ctx, c.Sender().ID)

	return c.Send("Done. Start a new conversation.", menu)
}

func (w *Wrapper) generateImageHandler(c tb.Context) error {
	_ = context.TODO()

	w.states.Set(c.Sender().ID, ImageDoing)

	return c.Send("Enter request in this format: ratio (3:2, 4:3, 16:9)\nAnd here description what "+
		"do you want to generate. Try to be more specific and precise in your request.", menu)
}

func (w *Wrapper) onTextHandler(c tb.Context) error {
	ctx := context.TODO()
	id := c.Sender().ID

	txt := c.Text()

	get, exist := w.states.Get(id)
	if !exist {
		return c.Send("Choose what do you want", menu)
	}

	switch get.State {
	case ChatDoing:
		answer, err := w.openaiSvc.ChatCompletion(ctx, c.Sender().ID, txt)
		if err != nil {
			return c.Send("Oops. There is a problem with your request.", menu)
		}

		return c.Send(answer, menu)
	case ImageDoing:
		err := c.Send("Preparing data...")
		if err != nil {
			return err
		}

		split := strings.Split(txt, "\n")
		if len(split) != 2 {
			return c.Send("Oops. There is a problem with your request. "+
				"You have to enter ratio, next row, description", menu)
		}

		ratio := split[0]
		prompt := split[1]

		aspectRatio, err := replicate.NewAspectRatio(ratio)
		if err != nil {
			logging.L(ctx).Error("wrong ratio", logging.ErrAttr(err), logging.StringAttr("ratio", ratio))
			return c.Send("Oops. Wrong ratio specified. Try again.", menu)
		}

		err = c.Send("Generating prompt for neural network based on your request...")
		if err != nil {
			logging.L(ctx).Error("failed to send message", logging.ErrAttr(err))
			return err
		}

		s1 := time.Now()

		openaiPrompt, err := w.openaiSvc.GenerateImagePrompt(ctx, prompt)
		if err != nil {
			logging.L(ctx).Error("failed to generate prompt", logging.ErrAttr(err))
			return c.Send("Oops. There is a problem with your request.", menu)
		}

		s2 := time.Now()

		err = c.Send("Prompt generated in " + s2.Sub(s1).String())
		if err != nil {
			logging.L(ctx).Error("failed to send message", logging.ErrAttr(err))
			return err
		}

		err = c.Send("Generating image...")
		if err != nil {
			logging.L(ctx).Error("failed to send message", logging.ErrAttr(err))
			return err
		}

		res, err := w.replicateSvc.GenerateImage(
			ctx,
			replicate.NewRequest(
				replicate.WithPrompt(openaiPrompt),
				replicate.WithRatio(aspectRatio),
			),
		)
		if err != nil {
			logging.L(ctx).Error("failed to generate image", logging.ErrAttr(err))
			return err
		}

		s3 := time.Now()

		err = c.Send("Image generated in " + s3.Sub(s2).String())
		if err != nil {
			logging.L(ctx).Error("failed to send message", logging.ErrAttr(err))
			return err
		}

		photo := &tb.Photo{
			File: tb.FromURL(res.Output),
		}

		return c.Send(photo)
	default:
		return c.Send("Choose what do you want", menu)
	}
}

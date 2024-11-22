package replicate

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/theartofdevel/logging"
)

var ErrBadRequest = errors.New("bad request")

type Service struct {
	token, url string
}

func NewService(cfg *Config) (*Service, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	url := "https://api.replicate.com/v1/models/black-forest-labs/flux-1.1-pro-ultra/predictions"

	return &Service{token: cfg.Token, url: url}, nil
}

func (s *Service) GenerateImage(ctx context.Context, reqGen *Request) (Response, error) {
	data, err := json.Marshal(reqGen)
	if err != nil {
		return Response{}, errors.Wrap(err, "json.Marshal")
	}

	// Создаем HTTP-запрос
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewBuffer(data))
	if err != nil {
		return Response{}, errors.Wrap(err, "http.NewRequestWithContext")
	}

	// Устанавливаем заголовки
	req.Header.Set("Authorization", "Bearer "+s.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "wait")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Response{}, errors.Wrap(err, "http.DefaultClient.Do")
	}

	defer func() {
		errBodyClose := resp.Body.Close()
		if errBodyClose != nil {
			logging.WithAttrs(
				ctx,
				logging.ErrAttr(errBodyClose),
			).Error("failed to close response body")
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, errors.Wrap(err, "io.ReadAll")
	}

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusCreated {
		logging.WithAttrs(
			ctx,
			logging.StringAttr("status", resp.Status),
			logging.IntAttr("status_code", resp.StatusCode),
			logging.StringAttr("body", string(body)),
		).Error("request to generate image failed")

		return Response{}, ErrBadRequest
	}

	// Парсим ответ
	var res Response
	if err = json.Unmarshal(body, &res); err != nil {
		return Response{}, errors.Wrap(err, "json.Unmarshal")
	}

	return res, nil
}

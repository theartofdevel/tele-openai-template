# tele-openai-template

Шаблон Telegram бота с интеграцией с [OpenAI](https://openai.com/) и [Replicate Flux 1.1 Pro Ultra](https://replicate.com/black-forest-labs/flux-1.1-pro-ultra).

Функции бота:
- Общение с OpenAI
- Сброс контекст и начало нового разговора
- Генерация изображения с промтом на любом языке

![SCR-20241201-mqvj (1) 2](https://github.com/user-attachments/assets/3fc573a4-9241-45cf-a6f4-5bf7f003d111)

# Сборка, конфигурация и запуск

1. Конфигурируем что хотим в `configs/config.yml`.
1. Создаем файл `.env` по образу и подобию в `.env.example`

```bash
$ cp .env.example .env
$ nano .env
$ cp docker-compose.example.yml docker-compose.yml
$ nano docker-compose.yml
```

## Docker 

### Собираем образ
```bash
$ make build-docker
```

### Собираем образ под линукс
```bash
$ make build-docker-linux
```

### Запускаем compose

```bash
$ docker compose up -d
```

## Goлый запуск

```bash
$ make lint
$ make test
$ make run
```

# Contributing

1. Клонируем и проверяем
```bash
$ git clone https://github.com/theartofdevel/tele-openai-template
$ make lint
$ make test
```
2. Создаем ветку
3. Повторяем линтер и тесты
```bash
$ make lint
$ make test
```
4. Создаем Pull Request 

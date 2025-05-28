package main

import (
	"log"

	"github.com/fire9900/golang-forum/internal/app"
	"github.com/fire9900/golang-forum/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %s", err.Error())
	}

	// Создание и запуск приложения
	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("Ошибка создания приложения: %s", err.Error())
	}

	if err := application.Run(); err != nil {
		log.Fatalf("Ошибка запуска приложения: %s", err.Error())
	}
}

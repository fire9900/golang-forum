package main

import (
	"log"

	"forum/internal/config"
	"forum/internal/handler"
	"forum/internal/middleware"
	"forum/internal/repository"
	"forum/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %s", err.Error())
	}

	// Подключение к базе данных
	db, err := repository.NewMySQLDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)

	// Инициализация обработчиков
	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)

	// Инициализация роутера
	router := gin.Default()

	// Группа API
	api := router.Group("/api")
	{
		// Публичные маршруты
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)

		// Защищенные маршруты
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(cfg.JWT.SecretKey))
		{
			// Профиль пользователя
			authorized.GET("/profile", userHandler.GetProfile)
			authorized.PUT("/profile", userHandler.UpdateProfile)
			authorized.PUT("/profile/password", userHandler.UpdatePassword)

			// Посты
			posts := authorized.Group("/posts")
			{
				posts.POST("/", postHandler.Create)
				posts.GET("/", postHandler.GetAll)
				posts.GET("/:id", postHandler.GetByID)
				posts.PUT("/:id", postHandler.Update)
				posts.DELETE("/:id", postHandler.Delete)
			}
		}
	}

	// Запуск сервера
	if err := router.Run(cfg.HTTP.Port); err != nil {
		log.Fatalf("Error starting server: %s", err.Error())
	}
}

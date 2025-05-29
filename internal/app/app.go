package app

import (
	"path/filepath"

	"github.com/fire9900/golang-forum/internal/auth"
	"github.com/fire9900/golang-forum/internal/config"
	"github.com/fire9900/golang-forum/internal/controllers"
	"github.com/fire9900/golang-forum/internal/middleware"
	"github.com/fire9900/golang-forum/internal/repository"
	"github.com/fire9900/golang-forum/internal/service"

	"github.com/gin-gonic/gin"
)

// App представляет собой структуру приложения
type App struct {
	cfg        *config.Config
	router     *gin.Engine
	authClient *auth.GrpcAuthClient
}

// New создает новый экземпляр приложения
func New(cfg *config.Config) (*App, error) {
	authClient, err := auth.NewGrpcAuthClient(cfg.Auth.GrpcAddress)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:        cfg,
		router:     gin.Default(),
		authClient: authClient,
	}, nil
}

// Run запускает приложение
func (a *App) Run() error {
	// Подключение к базе данных
	db, err := repository.NewMySQLDB(a.cfg)
	if err != nil {
		return err
	}

	// Выполнение миграций
	migrationsPath := filepath.Join("migrations")
	if err := repository.MigrateDB(db, migrationsPath); err != nil {
		return err
	}

	// Инициализация репозиториев
	postRepo := repository.NewPostRepository(db)

	// Инициализация сервисов
	postService := service.NewPostService(postRepo)

	// Инициализация обработчиков
	authController := controllers.NewAuthController(a.authClient)
	postController := controllers.NewPostController(postService)

	// Настройка маршрутов
	a.setupRoutes(authController, postController)

	// Запуск сервера
	return a.router.Run(a.cfg.HTTP.Port)
}

// setupRoutes настраивает маршруты приложения
func (a *App) setupRoutes(authController *controllers.AuthController, postController *controllers.PostController) {
	// Группа API
	api := a.router.Group("/api")
	{
		// Маршруты аутентификации
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshTokens)
		}

		// Публичные маршруты для постов
		posts := api.Group("/posts")
		{
			posts.GET("/", postController.GetAll)
			posts.GET("/:id", postController.GetByID)
		}

		// Защищенные маршруты
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(a.authClient))
		{
			// Защищенные маршруты для постов
			authorizedPosts := authorized.Group("/posts")
			{
				authorizedPosts.POST("/", postController.Create)
				authorizedPosts.PUT("/:id", postController.Update)
				authorizedPosts.DELETE("/:id", postController.Delete)
			}
		}
	}
}

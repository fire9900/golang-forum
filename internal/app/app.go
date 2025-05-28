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
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Инициализация сервисов
	userService := service.NewUserService(userRepo)
	postService := service.NewPostService(postRepo)

	// Инициализация обработчиков
	userController := controllers.NewUserController(userService)
	postController := controllers.NewPostController(postService)

	// Настройка маршрутов
	a.setupRoutes(userController, postController)

	// Запуск сервера
	return a.router.Run(a.cfg.HTTP.Port)
}

// setupRoutes настраивает маршруты приложения
func (a *App) setupRoutes(userController *controllers.UserController, postController *controllers.PostController) {
	// Группа API
	api := a.router.Group("/api")
	{
		// Публичные маршруты
		api.POST("/register", userController.Register)
		api.POST("/login", userController.Login)
		api.POST("/refresh", userController.RefreshTokens)

		// Защищенные маршруты
		authorized := api.Group("/")
		authorized.Use(middleware.AuthMiddleware(a.authClient))
		{
			// Профиль пользователя
			authorized.GET("/profile", userController.GetProfile)
			authorized.PUT("/profile", userController.UpdateProfile)
			authorized.PUT("/profile/password", userController.UpdatePassword)

			// Посты
			posts := authorized.Group("/posts")
			{
				posts.POST("/", postController.Create)
				posts.GET("/", postController.GetAll)
				posts.GET("/:id", postController.GetByID)
				posts.PUT("/:id", postController.Update)
				posts.DELETE("/:id", postController.Delete)
			}
		}
	}
}

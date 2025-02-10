package main

import (
	"context"
	"flag"
	"github.com/Dendyator/AntiBF/infrastructure/database"
	"github.com/Dendyator/AntiBF/infrastructure/infRepositories"
	"github.com/Dendyator/AntiBF/internal/delivery/grpc"
	api2 "github.com/Dendyator/AntiBF/internal/delivery/http"
	"github.com/Dendyator/AntiBF/internal/repositories"
	"github.com/Dendyator/AntiBF/pkg/config"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Dendyator/AntiBF/docs"           //nolint
	"github.com/Dendyator/AntiBF/internal/usecase" //nolint
	httpSwagger "github.com/swaggo/http-swagger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

// @title Anti-Bruteforce Service API
// @version 1.0
// @description API для сервиса, защищающего от брутфорс-атак
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	flag.Parse()

	logLevel := "info"
	appLogger := logger.New(logLevel)

	cfg := config.LoadConfig(configFile, appLogger)
	appLogger.Infof("Configuration loaded: %+v", cfg)

	db := database.NewDB(cfg.Database.DSN, cfg.Redis.Address, appLogger)
	defer db.Close()

	api2.InitLogger(appLogger)

	redisRepo := infRepositories.NewRedisRepo(&database.DB{}, appLogger)
	userRepo := repositories.NewUserRepository(&database.DB{}, appLogger)
	rateLimiter := usecase.NewRateLimiter(redisRepo, userRepo, cfg.RateLimiter, appLogger)

	ok := rateLimiter.CheckAuthorization("user1", "password123", "192.168.1.1")
	if ok {
		appLogger.Info("Authorization successful")
	} else {
		appLogger.Warn("Authorization failed")
	}

	appLogger.Debugf("RateLimiter initialized with limits: LoginLimit=%d, PasswordLimit=%d, IPLimit=%d",
		cfg.RateLimiter.LoginLimit, cfg.RateLimiter.PasswordLimit, cfg.RateLimiter.IPLimit)

	go func() {
		grpc.RunGRPCServer(rateLimiter, appLogger)
	}()

	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		api2.HandleAuth(rateLimiter).ServeHTTP(w, r)
	})
	http.HandleFunc("/manage_list", func(w http.ResponseWriter, r *http.Request) {
		api2.HandleManageList(rateLimiter).ServeHTTP(w, r)
	})
	http.HandleFunc("/check_list", func(w http.ResponseWriter, r *http.Request) {
		api2.HandleCheckList(rateLimiter).ServeHTTP(w, r)
	})
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	server := &http.Server{Addr: ":8080"}
	go func() {
		appLogger.Info("Starting HTTP server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatalf("Could not start HTTP server: %v", err)
		}
	}()

	<-stop
	appLogger.Info("Получен сигнал остановки. Завершение работы...")

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		appLogger.Errorf("Ошибка при закрытии сервера: %v", err)
	}
}

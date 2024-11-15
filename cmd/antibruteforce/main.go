package main

import (
	"flag"
	"net/http"

	"github.com/Dendyator/AntiBF/api"             //nolint
	_ "github.com/Dendyator/AntiBF/docs"          //nolint
	"github.com/Dendyator/AntiBF/internal/config" //nolint
	"github.com/Dendyator/AntiBF/internal/core"   //nolint
	"github.com/Dendyator/AntiBF/internal/db"     //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint
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

	db.InitDB(cfg.Database.DSN, appLogger)
	defer db.CloseDB()

	db.InitRedis(cfg.Redis.Address, appLogger)
	defer db.CloseRedis()

	api.InitLogger(appLogger)
	core.InitLogger(appLogger)
	core.InitRateLimiter(cfg.RateLimiter)
	appLogger.Debugf("RateLimiter initialized with limits: LoginLimit=%d, PasswordLimit=%d, IPLimit=%d",
		cfg.RateLimiter.LoginLimit, cfg.RateLimiter.PasswordLimit, cfg.RateLimiter.IPLimit)

	go func() {
		api.RunGRPCServer(appLogger)
	}()

	http.HandleFunc("/auth", api.HandleAuth)
	http.HandleFunc("/manage_list", api.HandleManageList)
	http.HandleFunc("/check_list", api.HandleCheckList)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("swagger-ui"))))

	appLogger.Info("Starting HTTP server on :8080")
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			appLogger.Fatalf("Could not start HTTP server: %v", err)
		}
	}()

	select {}
}

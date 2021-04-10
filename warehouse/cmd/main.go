package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"

	_binDeliveryHTTP "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/bin/delivery/http"
	_commodityDeliveryHTTP "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/commodity/delivery/http"
	_skuDeliveryHTTP "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/sku/delivery/http"
	_warehouseDeliveryHTTP "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/warehouse/delivery/http"

	_binRepository "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/bin/repository"
	_commodityRepository "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/commodity/repository"
	_skuRepository "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/sku/repository"
	_warehouseRepository "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/warehouse/repository"

	_binUsecase "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/bin/usecase"
	_commodityUsecase "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/commodity/usecase"
	_skuUsecase "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/sku/usecase"
	_warehouseUsecase "github.com/alvinradeka/jamblang-hakenton/warehouse/internal/warehouse/usecase"
)

var (
	viperInstance  *viper.Viper
	configData     AppConfig
	logrusInstance *logrus.Logger
	dbInstance     *sqlx.DB
	routerInstance *mux.Router
)

type (
	AppConfig struct {
		Logger LoggerConfig
		HTTP   HTTPConfig
		SQL    SQLConfig
	}

	LoggerConfig struct {
		Level string
	}

	HTTPConfig struct {
		Host string
		Port int64
	}

	SQLConfig struct {
		Host     string
		Port     int64
		Username string
		Password string
		DBName   string
	}
)

func main() {
	// Run Viper (Config Reader)
	// Config should be loaded in ./configs/config.yaml
	viperInstance = viper.New()
	viperInstance.AddConfigPath("./configs")
	viperInstance.SetConfigType("yaml")
	viperInstance.SetConfigName("config")

	if err := viperInstance.ReadInConfig(); err != nil {
		panic(err)
	}
	if err := viperInstance.Unmarshal(&configData); err != nil {
		panic(err)
	}

	// Run Logger
	logrusInstance = logrus.New()
	_logrusLevel, err := logrus.ParseLevel(configData.Logger.Level)
	if err != nil {
		panic(err)
	}
	logrusInstance.SetLevel(_logrusLevel)
	logrusInstance.SetFormatter(&logrus.TextFormatter{})
	logrusInstance.SetOutput(os.Stdout)

	// Run DB Instance
	dbInstance, err = sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", configData.SQL.Username, configData.SQL.Password, configData.SQL.Host, configData.SQL.Port, configData.SQL.DBName))
	if err != nil {
		logrusInstance.Fatalln(err)
	}

	// Build Repositories
	warehouseRepository := _warehouseRepository.NewSQL(logrusInstance, dbInstance)
	skuRepository := _skuRepository.NewSQL(logrusInstance, dbInstance)
	binRepository := _binRepository.NewSQL(logrusInstance, dbInstance)
	commodityRepository := _commodityRepository.NewSQL(logrusInstance, dbInstance)

	// Build Usecases
	warehouseUsecase := _warehouseUsecase.NewUsecase(logrusInstance, warehouseRepository, binRepository)
	skuUsecase := _skuUsecase.NewUsecase(logrusInstance, skuRepository)
	binUsecase := _binUsecase.NewUsecase(logrusInstance, binRepository, warehouseRepository)
	commodityUsecase := _commodityUsecase.NewUsecase(logrusInstance, commodityRepository)

	// Build Deliveries for HTTP
	routerInstance = mux.NewRouter()
	http.Handle("/", buildRouterHandle(logrusInstance, routerInstance))
	_warehouseDeliveryHTTP.NewHTTPDelivery(routerInstance, logrusInstance, warehouseUsecase)
	_skuDeliveryHTTP.NewHTTPDelivery(routerInstance, logrusInstance, skuUsecase)
	_binDeliveryHTTP.NewHTTPDelivery(routerInstance, logrusInstance, binUsecase)
	_commodityDeliveryHTTP.NewHTTPDelivery(routerInstance, logrusInstance, commodityUsecase)

	// Small Health Check
	routerInstance.HandleFunc("/sys/_health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Run HTTP Server
	logrusInstance.Infoln(fmt.Sprintf("HTTP Server Running At %s:%d", configData.HTTP.Host, configData.HTTP.Port))
	logrusInstance.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", configData.HTTP.Host, configData.HTTP.Port), nil))
}

func buildRouterHandle(log *logrus.Logger, h http.Handler) http.Handler {
	// Build Recover Function
	recover := handlers.RecoveryHandler(handlers.RecoveryLogger(log))

	// Wrap Handler with several middleware (Acts like Global Middleware)
	return handlers.LoggingHandler(log.Writer(),
		handlers.ProxyHeaders(
			handlers.CompressHandler(
				recover(h),
			),
		))
}

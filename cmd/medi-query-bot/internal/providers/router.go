package providers

//
//import (
//	"fa-quiz-next-gen/cmd/fa-quiz-next-gen/internal/cfg"
//	"fa-quiz-next-gen/cmd/fa-quiz-next-gen/internal/handlers"
//	"fa-quiz-next-gen/internal/ginLogger"
//	commonMiddleware "fa-quiz-next-gen/internal/ginMiddleware"
//	"fa-quiz-next-gen/internal/ginServer"
//	"fmt"
//	"os"
//
//	"github.com/gin-gonic/gin"
//	newrelic "github.com/newrelic/go-agent"
//	"github.com/rs/zerolog"
//	"github.com/spf13/viper"
//	_swaggerFiles "github.com/swaggo/files"
//	_ginSwagger "github.com/swaggo/gin-swagger"
//)
//
//// setupRouter setup router.
//func setupRouter(hs *handlers.Handlers) ginServer.GinRoutingFn {
//	return func(router *gin.Engine) {
//		cfgNR := newrelic.NewConfig(viper.GetString(cfg.ConfigNewRelicAppName), viper.GetString(cfg.ConfigNewRelicKey))
//		cfgNR.Logger = newrelic.NewDebugLogger(os.Stdout)
//		app, err := newrelic.NewApplication(cfgNR)
//		if nil != err {
//			fmt.Printf("add newrelic error: %s", err.Error())
//		}
//
//		router.Use(
//			commonMiddleware.RequestIDLoggingMiddleware(),
//			ginLogger.MiddlewareGin(AppName, zerolog.InfoLevel),
//			commonMiddleware.Recovery(),
//			nrgin.Middleware(app),
//		)
//
//		baseRoute := router.Group(viper.GetString(cfg.ConfigKeyContextPath))
//		baseRoute.GET("health-check", hs.HealthCheck.HealthCheck())
//		baseRoute.GET("swagger/*any", _ginSwagger.WrapHandler(_swaggerFiles.Handler))
//
//		v1 := baseRoute.Group("/v1")
//
//		medbotURL := v1.Group("/medbot")
//		{
//			medbotURL.GET("query", hs.SearchHandler.GetLinksMedbot())
//		}
//	}
//}

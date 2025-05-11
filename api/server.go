package api

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	db "p2platform/db/sqlc"
	"p2platform/errors"
	"p2platform/token"
	"p2platform/util"
	"strings"

	appErr "p2platform/errors"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
	producer   sarama.SyncProducer
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TelegramBotToken)
	if err != nil {
		return nil, appErr.ErrInternalServer
	}
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(config.KafkaBrokers, ","), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create Kafka producer: %w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		producer:   producer,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("source", validSource)
	}

	server.setupRouter()
	return server, nil
}

var StaticFS embed.FS

func (server *Server) setupRouter() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{viper.GetString("BASE_URL")},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.Use(func(ctx *gin.Context) {
		ctx.Set("BaseUrl", viper.GetString("BASE_URL"))
		ctx.Next()
	})

	router.SetHTMLTemplate(template.Must(template.ParseFS(StaticFS, "static/*.html")))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/", server.renderAuthPage)
	router.GET("/telegram-auth", server.renderIndexPage)
	router.GET("/sell-request/:id", server.renderSellRequestPage)
	router.GET("/create-sell-request", server.renderCreateSellRequestPage)
	router.GET("/create-buy-request", server.renderCreateBuyRequestPage)
	router.GET("/list-sell-requests", server.renderListSellRequestsPage)
	router.GET("/buy-request/:id", server.renderBuyRequestPage)
	router.GET("/list-buy-requests", server.renderListBuyRequestsPage)
	router.GET("/list-spaces", server.renderListSpacesPage)
	router.GET("/create-space", server.renderCreateSpacePage)

	api := router.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.GET("/telegram/init", server.initTelegramAuth)
			authGroup.GET("/telegram/status", server.checkAuthStatus)
			authGroup.POST("/telegram/webhook", server.handleTelegramWebhook)
		}

		authRoutes := api.Group("/").Use(CookieAuthMiddleware(server.tokenMaker))
		authRoutes.GET("/auth/me", server.getCurrentUser)
		authRoutes.GET("/sell-request/:id", server.getSellRequest)
		authRoutes.GET("/sell-requests", server.listSellRequests)
		authRoutes.GET("/buy-request/:id", server.getBuyRequest)
		authRoutes.GET("/buy-requests", server.listBuyRequests)
		authRoutes.POST("/sell-request", server.createSellRequest)
		authRoutes.POST("/buy-request", server.createBuyRequest)
		authRoutes.PATCH("/sell-request/:id", server.updateSellRequest)
		authRoutes.DELETE("/sell-request/:id", server.deleteSellRequest)
		authRoutes.POST("/buy-request/:id/close-confirm/seller", server.closeBuyRequestBySeller)
		authRoutes.POST("/buy-request/:id/close-confirm/buyer", server.closeBuyRequestByBuyer)
		authRoutes.POST("/buy-request/:id/close-confirm", server.CloseBuyRequestSellerBuyer)
		authRoutes.DELETE("/buy-request/:id", server.DeleteBuyRequest)
		authRoutes.GET("/sell-requests/my", server.listMySellRequests)
		authRoutes.GET("/buy-requests/my", server.listMyBuyRequests)
		authRoutes.POST("/space", server.createSpace)
		authRoutes.GET("/spaces", server.listSpaces)
		authRoutes.GET("/spaces/my", server.listMySpaces)
		authRoutes.GET("/space/:id", server.getSpace)
		authRoutes.POST("/space/:id/join", server.joinToSpace)

	}

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func ErrorResponse(err error) gin.H {
	if e, ok := err.(*errors.AppError); ok {
		return gin.H{
			"error": gin.H{
				"code":    e.Code,
				"message": e.Message,
			},
		}
	}
	// fallback
	return gin.H{
		"error": gin.H{
			"code":    "UNKNOWN_ERROR",
			"message": err.Error(),
		},
	}
}

func HandleAppError(ctx *gin.Context, err error) {
	if ae, ok := err.(*errors.AppError); ok {
		ctx.JSON(ae.Status, ErrorResponse(ae))
	} else {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
	}
}

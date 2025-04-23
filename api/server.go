package api

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	db "p2platform/db/sqlc"
	"p2platform/errors"
	"p2platform/util"
	"strings"

	"github.com/IBM/sarama"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	config   util.Config
	store    db.Store
	router   *gin.Engine
	producer sarama.SyncProducer
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(config.KafkaBrokers, ","), saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("cannot create Kafka producer: %w", err)
	}
	server := &Server{
		config:   config,
		store:    store,
		producer: producer,
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
		AllowOrigins:     []string{"http://localhost:8080", "https://facf-37-110-210-189.ngrok-free.app"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	router.SetHTMLTemplate(template.Must(template.ParseFS(StaticFS, "static/*.html")))
	router.SetHTMLTemplate(template.Must(template.ParseFS(StaticFS, "static/*.html")))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	router.GET("/", server.renderIndexPage)
	router.GET("/sell-request/:id", server.renderSellRequestPage)
	router.GET("/create-sell-request", server.renderCreateSellRequestPage)
	router.GET("/list-sell-requests", server.renderListSellRequestsPage)

	api := router.Group("/api/v1")
	{
		api.GET("/sell-request/:id", server.getSellRequest)
		api.GET("/sell-requests", server.listSellRequests)
		api.GET("/buy-request/:id", server.getBuyRequest)
		api.GET("/buy-requests", server.listBuyRequests)
		api.POST("/users/telegram", server.telegramAuth)

		authRoutes := api.Group("/").Use(CookieAuthMiddleware())

		authRoutes.POST("/sell-request", server.createSellRequest)
		authRoutes.POST("/buy-request", server.createBuyRequest)
		authRoutes.PATCH("/sell-request/:id", server.updateSellRequest)
		authRoutes.DELETE("/sell-request/:id", server.deleteSellRequest)
		authRoutes.POST("/buy-request/:id/close-confirm/seller", server.closeBuyRequestBySeller)
		authRoutes.POST("/buy-request/:id/close-confirm/buyer", server.closeBuyRequestByBuyer)
		authRoutes.DELETE("/buy-request/:id", server.DeleteBuyRequest)
		authRoutes.GET("/sell-requests/my", server.listMySellRequests)
		authRoutes.GET("/buy-requests/my", server.listMyBuyRequests)
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

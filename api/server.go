package api

import (
	db "p2platform/db/sqlc"
	"p2platform/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config util.Config
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store, config util.Config) (*Server, error) {
	server := &Server{
		config: config,
		store:  store,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("source", validSource)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.GET("/sell-request/:id", server.getSellRequest)
	router.GET("/sell-requests", server.listSellRequests)
	router.GET("/buy-request/:id", server.getBuyRequest)
	router.GET("/buy-requests", server.listBuyRequests)
	router.POST("/users/telegram", server.telegramAuth)

	authRoutes := router.Group("/").Use(CookieAuthMiddleware())

	authRoutes.POST("/sell-request", server.createSellRequest)
	authRoutes.POST("/buy-request", server.createBuyRequest)
	authRoutes.PATCH("/sell-request/:id", server.updateSellRequest)
	authRoutes.DELETE("/sell-request/:id", server.deleteSellRequest)
	authRoutes.POST("/buy-request/:id/close-confirm/seller", server.closeBuyRequestBySeller)
	authRoutes.POST("/buy-request/:id/close-confirm/buyer", server.closeBuyRequestByBuyer)
	authRoutes.DELETE("/buy-request/:id", server.DeleteBuyRequest)
	authRoutes.GET("/sell-requests/my", server.listMySellRequests)
	authRoutes.GET("/buy-requests/my", server.listMyBuyRequests)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}

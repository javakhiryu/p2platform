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
	router.POST("/sell-request", server.createSellRequest)
	router.GET("/sell-request/:id", server.getSellRequest)
	router.GET("/sell-requests", server.listSellRequest)
	router.PATCH("/sell-request/:id", server.updateSellRequest)
	router.DELETE("/sell-request/:id", server.deleteSellRequest)
	router.POST("/buy-request", server.createBuyRequest)
	router.GET("/buy-request/:id", server.getBuyRequest)
	router.GET("/buy-requests", server.listBuyRequests)
	router.POST("/buy-request/:id/close-confirm/seller", server.closeBuyRequestBySeller)
	router.POST("/buy-request/:id/close-confirm/buyer", server.closeBuyRequestByBuyer)
	//router.DELETE("/buy-request/:id", server.DeleteBuyRequest)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}

}

package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (server *Server) renderIndexPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.html", nil)
}

func (server *Server) renderSellRequestPage(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.HTML(http.StatusOK, "sell_request_page.html", gin.H{
		"id": id,
	})
}

func (server *Server) renderCreateSellRequestPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create_sell_request_page.html", nil)
}

func (server *Server) renderListSellRequestsPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "list_sell_requests_page.html", nil)
}

func (server *Server) renderCreateBuyRequestPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "create_buy_request_page.html", nil)
}

func (server *Server) renderBuyRequestPage(ctx *gin.Context) {
	id := ctx.Param("id")
	ctx.HTML(http.StatusOK, "buy_request_page.html", gin.H{
		"id": id,
	})
}

func (server *Server) renderListBuyRequestsPage(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "list_buy_requests_page.html", nil)
}
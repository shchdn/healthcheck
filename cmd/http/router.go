package main

import (
	"github.com/gin-gonic/gin"
	"template/internal/healthcheck"
)

func initRouter() *gin.Engine {
	router := gin.Default()
	setupRouter(router)
	return router
}

func setupRouter(router *gin.Engine) {
	h := healthcheck.New()
	go h.Start()
	router.GET(
		"/get_info", h.GetInfo,
	)
	router.GET(
		"/get_min", h.GetMin,
	)
	router.GET(
		"/get_max", h.GetMax,
	)
	router.GET(
		"/get_stats", h.GetMax,
	)
	router.GET(
		"/get_request_stats", h.GetRequestStats,
	)
}

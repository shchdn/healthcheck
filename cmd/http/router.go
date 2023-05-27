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
		"/getInfo", h.GetInfo,
	)
	router.GET(
		"/getMin", h.GetMin,
	)
	router.GET(
		"/getMax", h.GetMax,
	)
	router.GET(
		"/getStats", h.GetMax,
	)
	router.GET(
		"/getRequestStats", h.GetRequestStats,
	)
}

package http

import (
	"github.com/gin-gonic/gin"
)

// Router конфигурирует маршруты приложения
type Router struct {
	engine  *gin.Engine
	handler Handler
}

// NewRouter создает новый роутер
func NewRouter(h Handler) *Router {
	return &Router{
		engine:  gin.Default(),
		handler: h,
	}
}

// Setup настраивает все маршруты
func (r *Router) Setup() {

	// API v1
	v1 := r.engine.Group("/api/v1")
	{
		v1.POST("/webhook/bot-update", r.handler.SaveWebhookBotUpdate)
	};
}

// Run запускает HTTP сервер
func (r *Router) Run(port string) error {
	return r.engine.Run(":" + port)
}
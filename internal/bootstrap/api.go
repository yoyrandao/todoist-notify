package bootstrap

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPHost struct {
	server *gin.Engine
}

func NewHTTPHost() *HTTPHost {
	return &HTTPHost{server: gin.Default()}
}

func (s *HTTPHost) WithRouting(routing func(*gin.Engine), configureDefaultHealthEndpoint bool) *HTTPHost {
	if configureDefaultHealthEndpoint {
		s.server.GET("/health", func(ctx *gin.Context) {
			ctx.Status(http.StatusOK)
		})
	}

	routing(s.server)
	return s
}

func (h *HTTPHost) Run(port string) error {
	return h.server.Run(":" + port)
}

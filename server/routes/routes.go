package routes

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/scch94/Grecharge-gateway/server/handler"
	"github.com/scch94/Grecharge-gateway/server/middleware"
	"github.com/scch94/ins_log"
)

func SetUpRouter(ctx context.Context) *gin.Engine {

	//creamos el router de gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	//agregamos middleware de cors
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Permite todos los orígenes (cámbialo si quieres restringirlo)
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	//agregamos middleware gloables
	router.Use(middleware.GlobalMiddleware())

	h := handler.NewHandler()

	//metodos
	router.GET("/", h.Welcome)
	router.POST("/rechargeMobile", h.RechargeAccount)
	router.GET("/searchTransaction", h.SearchTransaction)
	router.NoRoute(notFoundHandler)

	return router
}

// Controlador para manejar rutas no encontradas
func notFoundHandler(c *gin.Context) {

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	ins_log.Errorf(ctx, "Route  not found: url: %v, method: %v", c.Request.RequestURI, c.Request.Method)
	c.JSON(http.StatusNotFound, nil)
}

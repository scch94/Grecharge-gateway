package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/scch94/ins_log"
)

func GlobalMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//generamos utdi y agregamos al contexto
		ctx := c.Request.Context()
		ctx = ins_log.SetPackageNameInContext(ctx, "middleware")
		utfi := ins_log.GenerateUTFI()
		ctx = ins_log.SetUTFIInContext(ctx, utfi)

		//copiamos el contexto de la solicitud
		c.Request = c.Request.WithContext(ctx)

		ins_log.Infof(ctx, "New petition received")
		ins_log.Tracef(ctx, "url: %v, method: %v", c.Request.RequestURI, c.Request.Method)
		startTime := time.Now() //registro de inicio de tiempo del request

		//pasamos la solicitud al siguiente middleware o al controlador final
		c.Next()
		elapsedTime := time.Since(startTime)
		ins_log.Infof(ctx, "Request took %v", elapsedTime)
	}
}

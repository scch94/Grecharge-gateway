package server

import (
	"context"
	"net/http"

	"github.com/scch94/Grecharge-gateway/config"
	"github.com/scch94/Grecharge-gateway/server/routes"
	"github.com/scch94/ins_log"
)

func StartServer(ctx context.Context) error {
	ctx = ins_log.SetPackageNameInContext(ctx, "server")

	ins_log.Infof(ctx, "Starting server on Port: %s", config.Config.ServPort)
	//usamos las rutas
	router := routes.SetUpRouter(ctx)
	serverConfig := &http.Server{
		Addr:    config.Config.ServPort,
		Handler: router,
	}
	err := serverConfig.ListenAndServe()
	if err != nil {
		ins_log.Errorf(ctx, "cant connect to the server: %+v", err)
		return err
	}
	return nil
}

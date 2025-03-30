package config

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/scch94/Gconfiguration"
	"github.com/scch94/ins_log"
)

var Config RechargeGatewayConfig

type RechargeGatewayConfig struct {
	LogLevel             string         `json:"log_level"`
	Log_name             string         `json:"log_name"`
	ServPort             string         `json:"server_port"`
	Client               Client         `json:"client"`
	RealizarVenta        EndpointConfig `json:"realizar_venta"`
	ConsultarTransaccion EndpointConfig `json:"consultar_transaccion"`
	MobileRegex          string         `json:"mobil_regex"`
	Acg                  Acg            `json:"acg"`
	TimeZone             string         `json:"time_zone"`
}
type Client struct {
	MaxIdleConns           int  `json:"maxIdleConns"`
	MaxConnsPerHost        int  `json:"maxConnsPerHost"`
	MaxIdleConnsPerHost    int  `json:"maxIdleConnsPerHost"`
	IdleConnTimeoutSeconds int  `json:"idleConnTimeoutSeconds"`
	DisableCompression     bool `json:"disableCompression"`
	PetitionsTimeOut       int  `json:"petitionsTimeOut"`
	DisableKeepAlives      bool `json:"disableKeepAlives"`
}
type Acg struct {
	IdMayorista int    `json:"id_mayorista"`
	IdProducto  int    `json:"id_producto"`
	Usuario     string `json:"usuario"`
	Clave       string `json:"clave"`
	IdCliente   string `json:"id_cliente"`
	Moneda      string `json:"moneda"`
	Canal       string `json:"canal"`
}

type EndpointConfig struct {
	URL    string `json:"url"`
	Method string `json:"method"`
}

func (r RechargeGatewayConfig) ConfigurationString() string {
	configJSON, err := json.Marshal(r)
	if err != nil {
		return fmt.Sprintf("Error al convertir la configuración a JSON: %v", err)
	}
	return string(configJSON)
}

func Upconfig(ctx context.Context) error {

	//traemos el contexto y le setiamos el contexto actual
	// Agregamos el valor "packageName" al contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "config")

	ins_log.Info(ctx, "starting to get the config struct ")
	err := Gconfiguration.GetConfig(&Config, "../../config", "recharGatewayConfig.json")

	if err != nil {
		ins_log.Fatalf(ctx, "error in Gconfiguration.GetConfig() ", err)
		return err
	}
	return nil
}

func WatchConfig(ctx context.Context) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ins_log.Errorf(ctx, "error creating watcher: %v", err)
		return
	}
	defer watcher.Close()

	// Agrega el directorio de configuración al watcher
	err = watcher.Add("../../config/recharGatewayConfig.json")
	if err != nil {
		ins_log.Errorf(ctx, "error adding config directory to watcher: %v", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			// Verificar si el archivo fue modificado
			if event.Op&fsnotify.Write == fsnotify.Write {
				ins_log.Infof(ctx, "Detected change in configuration file, reloading config")
				err := Upconfig(ctx)
				if err != nil {
					ins_log.Errorf(ctx, "Error reloading config: %v", err)
				} else {
					ins_log.Infof(ctx, "Configuration reloaded successfully")
					ins_log.SetLevel(Config.LogLevel)
					ins_log.SetService("")
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			ins_log.Errorf(ctx, "error watching config directory: %v", err)
		}
	}
}

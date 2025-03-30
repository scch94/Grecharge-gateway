package client

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/scch94/Grecharge-gateway/config"
	"github.com/scch94/Grecharge-gateway/models/request"
	"github.com/scch94/Grecharge-gateway/models/response"
	"github.com/scch94/ins_log"
)

func RealizarVenta(ctx context.Context, rechargeRequest request.RechargeMobile) (response.RechargeMobileResponse, error) {

	//traemos el contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "client")

	//preparamos la peticion
	var realizarVentaResponse response.RealizarVenta2Response
	var response response.RechargeMobileResponse

	ins_log.Infof(ctx, "starting to prepare the call to 'REALIZAR VENTA' whit the number: %v, and amount %v", rechargeRequest.Line, rechargeRequest.Amount)
	req, err := prepareRealizarVentaRequest(ctx, rechargeRequest.Line, rechargeRequest.Amount, rechargeRequest.IdTRN)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to prepareRequest()to 'REALIZAR VENTA'")
		return response, err
	}

	//hacemos el llamado
	realizarVentaResponse, err = callRealizarVenta(ctx, req)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to call 'REALIZAR VENTA'()")
		return response, err
	}

	//llamamos la funcion que nos logarea informacion completa sobre la solicitud
	response = generateResponse(ctx, realizarVentaResponse)

	return response, nil
}

func prepareRealizarVentaRequest(ctx context.Context, number string, amount int, transactionId string) (*http.Request, error) {

	//GENERAMOS EL CUERPO DE LA SOLICITUD
	venta := createRequestBody(number, amount, transactionId)

	ventaBody := &request.VentaBody{
		RealizarVenta: *venta,
	}

	bodyToRealizarVenta, err := request.CreateBodyToVenta(ventaBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(config.Config.RealizarVenta.Method, config.Config.RealizarVenta.URL, strings.NewReader(bodyToRealizarVenta))
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to 'REALIZAR VENTA': %v", err.Error())
		return nil, err
	}
	ins_log.Tracef(ctx, "petition http created")

	req.Header.Set("Content-Type", "text/xml;charset=UTF-8")
	req.Header.Del("soapaction")
	req.Header["SOAPAction"] = []string{`""`}

	//logueamos la url y el body final de la peticion
	ins_log.Infof(ctx, "Final url : %s", config.Config.RealizarVenta.URL)
	ins_log.Infof(ctx, "Final BODY: %s", bodyToRealizarVenta)

	return req, nil

}

func callRealizarVenta(ctx context.Context, req *http.Request) (response.RealizarVenta2Response, error) {

	//creamos la variable donde guardaremos la respuesta
	var realizarVentaResponse response.RealizarVenta2Response

	start := time.Now()

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to 'REALIZAR VENTA': %s", err)
		return realizarVentaResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "Request to 'REALIZAR VENTA' took %v", duration)

	// Confirmamos que la respuesta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "error due to non-200 status code: %v", err)
		return realizarVentaResponse, err
	}

	// Logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return realizarVentaResponse, err
	}

	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "HTTP Status Response: %d", statusCode)
	ins_log.Infof(ctx, "RESPONSE BODY: %s", string(responseBody))

	// Parseamos el resultado con lo que esperamos recibir
	err = xml.Unmarshal(responseBody, &realizarVentaResponse)
	if err != nil {
		ins_log.Errorf(ctx, " Error decoding the response: %s", err)
		return realizarVentaResponse, err
	}
	return realizarVentaResponse, nil
}

func createRequestBody(number string, amount int, transactionId string) *request.RealizarVenta {
	loc, err := time.LoadLocation(config.Config.TimeZone)
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)
	formatted := now.Format("2006-01-02T15:04:05.000-07:00")

	venta := &request.RealizarVenta{
		IDMayorista:                            config.Config.Acg.IdMayorista,
		IDCliente:                              config.Config.Acg.IdCliente,
		Usuario:                                config.Config.Acg.Usuario,
		Clave:                                  config.Config.Acg.Clave,
		FechaCliente:                           formatted,
		IDTransaccionCliente:                   transactionId,
		IDProducto:                             config.Config.Acg.IdProducto,
		Importe:                                amount,
		CaracteristicaTelefono:                 0,
		NumeroTelefono:                         0,
		Certificado:                            "",
		TipoTRX:                                "ON",
		Moneda:                                 config.Config.Acg.Moneda,
		IDClienteExt:                           "",
		IdentifTerminal:                        "",
		SaldoDeSubeAntesDeRealizarUnaVenta:     "",
		TelefonoCompletoOTarjeta:               number,
		ProdIdentif:                            "",
		ProdDatosAdic:                          "",
		ModeloDeTerminal:                       "",
		ImeiDeSIM:                              "",
		ImeiTerminal:                           "",
		SaldoDelCliParaInformarAlOperador:      0,
		IDTipoTerminalParaInformarAlOperador:   0,
		DescTipoTerminalParaInformarAlOperador: "",
		HashDeSeguridad:                        "",
		SubeIDDistribuidorRed:                  "",
		Canal:                                  config.Config.Acg.Canal,
	}

	return venta
}
func generateResponse(ctx context.Context, v response.RealizarVenta2Response) response.RechargeMobileResponse {
	var response response.RechargeMobileResponse

	out := v.Body.RealizarVentaResponse.Out
	state := "-1010" // Default unknown error
	trnId := ""
	balance := "-1.0"

	switch {
	// ✅ Successful transaction
	case out.IDEstadoTransaccion == APROBADA && out.Error.HayError == "false":
		state = "-1000"
		trnId = out.IDTransaccion
		balance = out.SaldoDisponibleDelCliente
		ins_log.Infof(ctx, "Recharge completed successfully. External ID: %v | Final user balance: %v", out.IDTransaccion, out.SaldoDisponibleDelCliente)

	// 🔁 Duplicate transaction
	case out.Error.HayError == "true" && out.Error.CodigoError == TRANSACCION_DUPLICADA:
		state = "-1008"
		ins_log.Errorf(ctx, "Recharge rejected due to duplication. Error message: %v", out.Error.MsgError)

	// ⏳ Recharge in process
	case out.IDEstadoTransaccion == CARGA_EN_PROCESO && out.Error.HayError == "false":
		state = "-1009"
		ins_log.Warnf(ctx, "Recharge is in progress. It should be consulted later.")

	// ⚠️ Invalid line / subscriber does not exist
	case out.IDEstadoTransaccion == LINEA_INVALIDA && out.Error.HayError == "false":
		state = "-1005"
		ins_log.Errorf(ctx, "Subscriber does not exist. Description: %v", out.DescEstadoTransaccion)

	// ❌ Any other backend error
	case out.Error.HayError == "true":
		state = "-1008"
		ins_log.Errorf(ctx, "Recharge failed due to backend error. Code: %v | Message: %v", out.Error.CodigoError, out.Error.MsgError)

	// ❗ Rejected by business rules
	case out.Error.HayError == "false":
		ins_log.Errorf(ctx, "Recharge rejected by business rules. Status: %v | Description: %v", out.IDEstadoTransaccion, out.DescEstadoTransaccion)
	}

	response.ReloadResult.State = state
	response.ReloadResult.TrnId = trnId
	response.ReloadResult.Balance = balance

	return response
}

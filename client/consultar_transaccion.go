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

func ConsultarTransaccion(ctx context.Context, searchTransactionRequest request.SearchTransaction) (response.ConsultarTransaccion2Response, error) {

	//traemos el contexto
	ctx = ins_log.SetPackageNameInContext(ctx, "client")

	//preparamos la peticion
	var consultarTransaccionResponse response.ConsultarTransaccion2Response

	ins_log.Infof(ctx, "starting to prepare the call to 'CONSULTAR TRANSACCION' whit the number: %v, and amount, %v", searchTransactionRequest.Line, searchTransactionRequest.Amount)
	req, err := prepareConsultarTransaccion(ctx, searchTransactionRequest.Line, searchTransactionRequest.Amount, searchTransactionRequest.IdTRN, searchTransactionRequest.IdTRNClient)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to prepareRequest()to 'CONSULTAR TRANSACCION'")
		return consultarTransaccionResponse, err
	}

	//hacemos llamado
	consultarTransaccionResponse, err = callConsultarTransaccion(ctx, req)
	if err != nil {
		ins_log.Errorf(ctx, "error when we try to call 'CONSULTAR TRANSACCION'()")
		return consultarTransaccionResponse, err
	}
	return consultarTransaccionResponse, nil
}

func prepareConsultarTransaccion(ctx context.Context, number string, amount string, transactionId string, transactionIdClient string) (*http.Request, error) {

	//GENERAMOS EL CUERPO DE LA SOLICITUD
	ConsultarTransaccion := createConsultarTransaccionRequestBody(number, amount, transactionId, transactionIdClient)

	ConsultarTransaccionBody := &request.ConsultarTransaccionBody{
		ConsultarTransaccion: *ConsultarTransaccion,
	}

	bodyToRealizarVenta, err := request.CreateBodyToRealizarCuenta(ConsultarTransaccionBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(config.Config.ConsultarTransaccion.Method, config.Config.ConsultarTransaccion.URL, strings.NewReader(bodyToRealizarVenta))
	if err != nil {
		ins_log.Errorf(ctx, "Error creating request to 'CONSULTAR TRANSACCION': %v", err.Error())
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

func callConsultarTransaccion(ctx context.Context, req *http.Request) (response.ConsultarTransaccion2Response, error) {

	//creamos la variable dodne guardaremos la respuesta del consultar transaccion
	var consultarTransaccionResponse response.ConsultarTransaccion2Response

	start := time.Now()

	resp, err := Client.Do(req)
	if err != nil {
		ins_log.Errorf(ctx, "Error when we do the petition to 'CONSULTAR TRANSACCION': %s", err)
		return consultarTransaccionResponse, err
	}
	defer resp.Body.Close()
	duration := time.Since(start)
	ins_log.Infof(ctx, "Request to 'CONSULTAR TRANSACCION' took %v", duration)

	// Confirmamos que la respuesta sea 200
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
		ins_log.Errorf(ctx, "error due to non-200 status code: %v", err)
		return consultarTransaccionResponse, err
	}

	// Logueamos lo que recibimos
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		ins_log.Errorf(ctx, "Error reading response body: %s", err)
		return consultarTransaccionResponse, err
	}
	// Imprimir la respuesta recibida
	statusCode := resp.StatusCode
	ins_log.Infof(ctx, "HTTP Status Response: %d", statusCode)
	ins_log.Infof(ctx, "RESPONSE BODY: %s", string(responseBody))

	// Parseamos el resultado con lo que esperamos recibir
	err = xml.Unmarshal(responseBody, &consultarTransaccionResponse)
	if err != nil {
		ins_log.Errorf(ctx, " Error decoding the response: %s", err)
		return consultarTransaccionResponse, err
	}
	return consultarTransaccionResponse, nil
}
func createConsultarTransaccionRequestBody(number string, amount string, transactionId string, transactionIdClient string) *request.ConsultarTransaccion {
	loc, err := time.LoadLocation(config.Config.TimeZone)
	if err != nil {
		panic(err)
	}

	now := time.Now().In(loc)
	formatted := now.Format("2006-01-02T15:04:05.000-07:00")

	consulta := &request.ConsultarTransaccion{
		IDMayorista:            config.Config.Acg.IdMayorista,
		IDCliente:              config.Config.Acg.IdCliente,
		IDUsuario:              "0",
		Usuario:                config.Config.Acg.Usuario,
		Clave:                  config.Config.Acg.Clave,
		FechaCliente:           formatted,
		IDTransaccion:          transactionId,
		IDTransaccionCliente:   transactionIdClient,
		IDProducto:             config.Config.Acg.IdProducto,
		Importe:                amount,
		CaracteristicaTelefono: 0,
		NumeroTelefono:         number,
		Certificado:            "",
		UltimaTransaccion:      false,
	}

	return consulta
}

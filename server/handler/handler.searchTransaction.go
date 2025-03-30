package handler

import (
	"context"
	"errors"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/scch94/Grecharge-gateway/client"
	"github.com/scch94/Grecharge-gateway/config"
	"github.com/scch94/Grecharge-gateway/models/request"
	"github.com/scch94/ins_log"
)

func (h *Handler) SearchTransaction(c *gin.Context) {
	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	ins_log.Infof(ctx, "startint to SearchTransaction")
	ins_log.Tracef(ctx, "starting to get the params of the petition")

	//creamos la variable de tipo struct para guardar los datos de la peticon
	searchTransactionRequest := request.SearchTransaction{
		Line:        c.Query("line"),
		Amount:      c.Query("amount"),
		IdTRN:       c.Query("transactionId"),
		IdTRNClient: c.Query("transactionIdClient"),
	}

	ins_log.Tracef(ctx, "params in the petition: %s", searchTransactionRequest)

	//validamos la solicitud
	err := checkRequestSearchTransaction(ctx, searchTransactionRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error in the function checkRequestSearchTransaction() err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	response, err := client.ConsultarTransaccion(ctx, searchTransactionRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error triyin to cole ConsultarTransaccion() err:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.XML(http.StatusOK, response)
}

func checkRequestSearchTransaction(ctx context.Context, searchTransactionRequest request.SearchTransaction) error {

	//chequeamos que el numero pase la expresion regular
	regex, err := regexp.Compile(config.Config.MobileRegex)
	if err != nil {
		ins_log.Errorf(ctx, "error to compilate the regex expression function regexp.Compile(): , err: %v", err)
		return errors.New("error to compilate the regex expression function")
	}
	if !regex.MatchString(searchTransactionRequest.Line) {
		ins_log.Errorf(ctx, "Mobile did not match in the regex expression")
		return errors.New("mobile did not match in the regex expression regex.MatchString()")
	}
	//si llegamos aca la expresion regular si valido todo.
	ins_log.Tracef(ctx, "value match with the regex expression")

	//chequeamos que el amount no venga vacio
	if searchTransactionRequest.Amount == "" {
		ins_log.Errorf(ctx, "the Amount is empity")
		return errors.New("the Amount is empity")
	}
	ins_log.Tracef(ctx, "amount check pass")

	//chequeamos el idtranssaccion no venga vacio
	if searchTransactionRequest.IdTRN == "" {
		ins_log.Errorf(ctx, "the transaction id is empity")
		return errors.New("the transaction id is empity")
	}
	ins_log.Tracef(ctx, "transaction id check pass")

	return nil
}

package handler

import (
	"context"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/scch94/Grecharge-gateway/client"
	"github.com/scch94/Grecharge-gateway/config"
	"github.com/scch94/Grecharge-gateway/models/request"
	"github.com/scch94/ins_log"
)

func (h *Handler) RechargeAccount(c *gin.Context) {

	//traemos el contexto y le setiamos el contexto actual
	ctx := c.Request.Context()
	ctx = ins_log.SetPackageNameInContext(ctx, "handler")

	ins_log.Infof(ctx, "startint to recharge account")
	ins_log.Tracef(ctx, "starting to get the body of the petition")

	//leemos el cuerpo de la solicitud
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		ins_log.Errorf(ctx, "error reading the body of the petition err: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	//creamos la variable de tipo struct para guardar los datos de la peticon
	var rechargeRequest request.RechargeMobile

	err = xml.Unmarshal(body, &rechargeRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error reading the body of the petition err: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ins_log.Tracef(ctx, "Body of the petition: %s", body)

	//validamos la solicitud
	err = checkRequestRechargeAccount(ctx, rechargeRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error in the function checkRequestRechargeAccount() err: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	ins_log.Infof(ctx, "the petition pass de validation")

	//hacemos el llamado a la api, la cual devuelve el xml de respuesta
	responseRealizarVenta, err := client.RealizarVenta(ctx, rechargeRequest)
	if err != nil {
		ins_log.Errorf(ctx, "error triyin to cole RealizarVenta() err:%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ins_log.Infof(ctx, "finish Recharge account")

	c.XML(http.StatusOK, responseRealizarVenta)
}

func checkRequestRechargeAccount(ctx context.Context, rechargeRequest request.RechargeMobile) error {

	//chequeamos que el numero pase la expresion regular
	regex, err := regexp.Compile(config.Config.MobileRegex)
	if err != nil {
		ins_log.Errorf(ctx, "error to compilate the regex expression function regexp.Compile(): , err: %v", err)
		return errors.New("error to compilate the regex expression function")
	}
	if !regex.MatchString(rechargeRequest.Line) {
		ins_log.Errorf(ctx, "Mobile did not match in the regex expression")
		return errors.New("mobile did not match in the regex expression regex.MatchString()")
	}
	//si llegamos aca la expresion regular si valido todo.
	ins_log.Tracef(ctx, "mobile: %v match with the regex expression:%v", rechargeRequest.Line, config.Config.MobileRegex)

	//chequeamos que el amount no venga vacio
	if rechargeRequest.Amount == 0 {
		ins_log.Errorf(ctx, "the Amount is empity")
		return errors.New("the Amount is empity")
	}
	ins_log.Tracef(ctx, "amount:%v check pass", rechargeRequest.Amount)

	//chequeamos el idtranssaccion no venga vacio
	if rechargeRequest.IdTRN == "" {
		ins_log.Errorf(ctx, "the transaction id is empity")
		return errors.New("the transaction id is empity")
	}
	ins_log.Tracef(ctx, "transaction id:%v check pass", rechargeRequest.IdTRN)

	return nil
}

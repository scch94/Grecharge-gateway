package response

import (
	"encoding/xml"
)

type RealizarVenta2Response struct {
	XMLName xml.Name          `xml:"Envelope"`
	Body    BodyRealizarVenta `xml:"Body"`
}

type BodyRealizarVenta struct {
	RealizarVentaResponse Nsl `xml:"realizarVenta2Response"`
}

type Nsl struct {
	Out Result `xml:"out"`
}

type Result struct {
	CaracteristicaTelefono    string     `xml:"caracteristicaTelefono"`
	DescEstadoTransaccion     string     `xml:"descEstadoTransaccion"`
	DescripcionProducto       string     `xml:"descripcionProducto"`
	Error                     ErrorBlock `xml:"error"`
	Fecha                     string     `xml:"fecha"`
	IDCliente                 string     `xml:"idCliente"`
	IDEstadoTransaccion       string     `xml:"idEstadoTransaccion"`
	IDProducto                string     `xml:"idProducto"`
	IDTransaccion             string     `xml:"idTransaccion"`
	Importe                   string     `xml:"importe"`
	ImporteDestino            string     `xml:"importeDestino"`
	Moneda                    string     `xml:"moneda"`
	MonedaDestino             string     `xml:"monedaDestino"`
	NumeroTelefono            string     `xml:"numeroTelefono"`
	RazonSocialCliente        string     `xml:"razonSocialCliente"`
	TipoOperacion             string     `xml:"tipoOperacion"`
	Usuario                   string     `xml:"usuario"`
	SaldoDisponibleDelCliente string     `xml:"saldoDisponibleDelCliente"`
}

type ErrorBlock struct {
	CodigoError string `xml:"codigoError"`
	HayError    string `xml:"hayError"`
	MsgError    string `xml:"msgError"`
}

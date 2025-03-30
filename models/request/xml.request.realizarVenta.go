package request

import (
	"encoding/xml"
)

type CallVenta struct {
	XMLName xml.Name  `xml:"soapenv:Envelope"`
	Soapenv string    `xml:"xmlns:soapenv,attr"`
	Q0      string    `xml:"xmlns:q0,attr"`
	Xsd     string    `xml:"xmlns:xsd,attr"`
	Xsi     string    `xml:"xmlns:xsi,attr"`
	Header  string    `xml:"soapenv:Header"`
	Body    VentaBody `xml:"soapenv:Body"`
}

type VentaBody struct {
	RealizarVenta RealizarVenta `xml:"q0:realizarVenta2"`
}

type RealizarVenta struct {
	IDMayorista                            int    `xml:"idMayorista"`
	IDCliente                              string `xml:"idCliente"`
	Usuario                                string `xml:"usuario"`
	Clave                                  string `xml:"clave"`
	FechaCliente                           string `xml:"fechaCliente"`
	IDTransaccionCliente                   string `xml:"idTransaccionCliente"`
	IDProducto                             int    `xml:"idProducto"`
	Importe                                int    `xml:"importe"`
	CaracteristicaTelefono                 int    `xml:"caracteristicaTelefono"`
	NumeroTelefono                         int    `xml:"numeroTelefono"`
	Certificado                            string `xml:"certificado"`
	TipoTRX                                string `xml:"tipoTRX"`
	Moneda                                 string `xml:"moneda"`
	IDClienteExt                           string `xml:"idClienteExt"`
	IdentifTerminal                        string `xml:"identifTerminal"`
	SaldoDeSubeAntesDeRealizarUnaVenta     string `xml:"saldoDeSubeAntesDeRealizarUnaVenta"`
	TelefonoCompletoOTarjeta               string `xml:"telefonoCompletoOTarjeta"`
	ProdIdentif                            string `xml:"prodIdentif"`
	ProdDatosAdic                          string `xml:"prodDatosAdic"`
	ModeloDeTerminal                       string `xml:"modeloDeTerminal"`
	ImeiDeSIM                              string `xml:"imeiDeSIM"`
	ImeiTerminal                           string `xml:"imeiTerminal"`
	SaldoDelCliParaInformarAlOperador      int    `xml:"saldoDelCliParaInformarAlOperador"`
	IDTipoTerminalParaInformarAlOperador   int    `xml:"idTipoTerminalParaInformarAlOperador"`
	DescTipoTerminalParaInformarAlOperador string `xml:"descTipoTerminalParaInformarAlOperador"`
	HashDeSeguridad                        string `xml:"hashDeSeguridad"`
	SubeIDDistribuidorRed                  string `xml:"sube_idDistribuidorRed"`
	Canal                                  string `xml:"canal"`
}

func CreateBodyToVenta(data *VentaBody) (string, error) {
	envelope := CallVenta{
		Soapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		Q0:      "http://service.core.cargavirtual.americacg.com",
		Xsd:     "http://www.w3.org/2001/XMLSchema",
		Xsi:     "http://www.w3.org/2001/XMLSchema-instance",
		Header:  "",
		Body:    *data,
	}

	xmlData, err := xml.MarshalIndent(envelope, "", "\t")
	if err != nil {
		return "", err
	}

	return string(xmlData), nil
}

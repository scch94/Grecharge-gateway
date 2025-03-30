package request

import "encoding/xml"

type CallConsultarTransaccion struct {
	XMLName xml.Name                 `xml:"soapenv:Envelope"`
	Soapenv string                   `xml:"xmlns:soapenv,attr"`
	Q0      string                   `xml:"xmlns:q0,attr"`
	Xsd     string                   `xml:"xmlns:xsd,attr"`
	Xsi     string                   `xml:"xmlns:xsi,attr"`
	Header  string                   `xml:"soapenv:Header"`
	Body    ConsultarTransaccionBody `xml:"soapenv:Body"`
}

type ConsultarTransaccionBody struct {
	ConsultarTransaccion ConsultarTransaccion `xml:"consultarTransaccion2"`
}

type ConsultarTransaccion struct {
	IDMayorista            int    `xml:"idMayorista"`
	IDCliente              string `xml:"idCliente"`
	IDUsuario              string `xml:"idUsuario"`
	Clave                  string `xml:"clave"`
	FechaCliente           string `xml:"fechaCliente"`
	IDTransaccion          string `xml:"idTransaccion"`
	IDTransaccionCliente   string `xml:"idTransaccionCliente"`
	IDProducto             int    `xml:"idProducto"`
	Importe                string `xml:"importe"`
	CaracteristicaTelefono int    `xml:"caracteristicaTelefono"`
	NumeroTelefono         string `xml:"numeroTelefono"`
	Certificado            string `xml:"certificado"`
	Usuario                string `xml:"usuario"`
	UltimaTransaccion      bool   `xml:"ultimaTransaccion"`
}

func CreateBodyToRealizarCuenta(data *ConsultarTransaccionBody) (string, error) {
	envelope := CallConsultarTransaccion{
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

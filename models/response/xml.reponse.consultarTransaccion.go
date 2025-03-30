package response

import "encoding/xml"

type ConsultarTransaccion2Response struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    Body     `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
}

type Body struct {
	Response Nsl `xml:"http://service.core.cargavirtual.americacg.com consultarTransaccion2Response"`
}

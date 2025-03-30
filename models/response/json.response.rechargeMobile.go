package response

import "encoding/xml"

type RechargeMobileResponse struct {
	XMLName      xml.Name             `xml:"http://192.168.3.11/RcrgMyrst/ ReloadResponse"`
	ReloadResult rechargeMobileResult `xml:"ReloadResult"`
}
type rechargeMobileResult struct {
	State   string `xml:"State"`
	TrnId   string `xml:"TrnId"`
	Balance string `xml:"Balance"`
}

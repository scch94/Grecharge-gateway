package request

import "encoding/xml"

type RechargeMobile struct {
	XMLName       xml.Name `xml:"http://192.168.3.11/RcrgMyrst/ Reload"`
	Prefix        string   `xml:"Prefix"`
	Line          string   `xml:"Line"`
	Amount        int      `xml:"Amount"`
	IdTRN         string   `xml:"IdTRN"`
	Product       int      `xml:"Product"`
	IdDistributor int      `xml:"IdDistributor"`
	IdLocation    int      `xml:"IdLocation"`
	Password      string   `xml:"Password"`
	IdPOS         int      `xml:"IdPOS"`
}

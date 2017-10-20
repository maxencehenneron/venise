package structures

import "encoding/xml"

type Ref struct {
	XMLName xml.Name `xml:"nd"`
	Ref     int64    `xml:"ref,attr"`
}

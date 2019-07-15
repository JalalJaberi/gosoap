package gosoap

import (
	"encoding/xml"
	"fmt"
)

// Soap Response
type Response struct {
	Body    []byte
	Header  []byte
	Payload []byte
}

// Unmarshal get the body and unmarshal into the interface
func (r *Response) Unmarshal(v interface{}) error {

	// fmt.Printf(string(r.Body[:]))

	if len(r.Body) == 0 {
		return fmt.Errorf("Body is empty")
	}

	// fmt.Printf("r.Body: %s\r", r.Body)

	var f Fault
	xml.Unmarshal(r.Body, &f)
	if f.Code != "" {
		return fmt.Errorf("[%s]: %s", f.Code, f.Description)
	}

	return xml.Unmarshal(r.Body, v)
}

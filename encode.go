package gosoap

import (
	"encoding/xml"
	"fmt"
	"reflect"
	"strconv"
)

// MarshalXML envelope the body and encode to xml
func (c process) MarshalXML(e *xml.Encoder, _ xml.StartElement) error {
	tokens := &tokenData{}

	//start envelope
	if c.Client.Definitions == nil {
		return fmt.Errorf("definitions is nil")
	}

	tokens.startEnvelope()
	if len(c.Client.HeaderParams) > 0 {
		tokens.startHeader(c.Client.HeaderName, c.Client.Definitions.Types[0].XsdSchema[0].TargetNamespace)
		for k, v := range c.Client.HeaderParams {
			t := xml.StartElement{
				Name: xml.Name{
					Space: "",
					Local: k,
				},
			}

			tokens.data = append(tokens.data, t, xml.CharData(v), xml.EndElement{Name: t.Name})
		}

		tokens.endHeader(c.Client.HeaderName)
	}

	err := tokens.startBody(c.Request.Method, c.Client.Definitions.Types[0].XsdSchema[0].TargetNamespace)
	if err != nil {
		return err
	}

	tokens.recursiveEncode(c.Request.Params)

	//end envelope
	tokens.endBody(c.Request.Method)
	tokens.endEnvelope()

	for _, t := range tokens.data {
		err := e.EncodeToken(t)
		if err != nil {
			return err
		}
	}

	return e.Flush()
}

type tokenData struct {
	data []xml.Token
}

func (tokens *tokenData) recursiveEncode(hm interface{}) {
	v := reflect.ValueOf(hm)

	switch v.Kind() {
	case reflect.Map:
		for _, key := range v.MapKeys() {
			t := xml.StartElement{
				Name: xml.Name{
					Space: "",
					Local: key.String(),
				},
			}

			tokens.data = append(tokens.data, t)
			tokens.recursiveEncode(v.MapIndex(key).Interface())
			tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			inner := v.Index(i).Interface()
			v2 := reflect.ValueOf(inner)
			switch v2.Kind() {
			case reflect.String:
				t := xml.StartElement{
					Name: xml.Name{
						Space: "",
						Local: "string",
					},
				}
				tokens.data = append(tokens.data, t)
				tokens.recursiveEncode(inner)
				tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
			case reflect.Bool:
				t := xml.StartElement{
					Name: xml.Name{
						Space: "",
						Local: "bool",
					},
				}
				tokens.data = append(tokens.data, t)
				tokens.recursiveEncode(inner)
				tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64:
				t := xml.StartElement{
					Name: xml.Name{
						Space: "",
						Local: "int",
					},
				}
				tokens.data = append(tokens.data, t)
				tokens.recursiveEncode(inner)
				tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
			case reflect.Float32,
				reflect.Float64:
				t := xml.StartElement{
					Name: xml.Name{
						Space: "",
						Local: "float",
					},
				}
				tokens.data = append(tokens.data, t)
				tokens.recursiveEncode(inner)
				tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
			//case reflect
			default:
				tokens.recursiveEncode(inner)
			}
		}
	case reflect.Struct:
		tp := reflect.TypeOf(hm)
		t := xml.StartElement{
			Name: xml.Name{
				Space: "",
				Local: reflect.TypeOf(hm).Name(),
			},
		}
		tokens.data = append(tokens.data, t)
		for i := 0; i < v.NumField(); i++ {
			inner := v.Field(i).Interface()
			v2 := reflect.ValueOf(inner)
			t2 := xml.StartElement{
				Name: xml.Name{
					Space: "",
					Local: tp.Field(i).Name,
				},
			}
			tokens.data = append(tokens.data, t2)
			show := true
			tag := reflect.TypeOf(hm).Field(i).Tag.Get("typeAsTag")
			if "false" == tag {
				show = false
			}
			var t3 xml.StartElement
			switch v2.Kind() {
			case reflect.String:
				if show {
					t3 = xml.StartElement{
						Name: xml.Name{
							Space: "",
							Local: "string",
						},
					}
					tokens.data = append(tokens.data, t3)
				}
				tokens.recursiveEncode(inner)
				if show {
					tokens.data = append(tokens.data, xml.EndElement{Name: t3.Name})
				}
			case reflect.Bool:
				if show {
					t3 = xml.StartElement{
						Name: xml.Name{
							Space: "",
							Local: "bool",
						},
					}
					tokens.data = append(tokens.data, t3)
				}
				tokens.recursiveEncode(inner)
				if show {
					tokens.data = append(tokens.data, xml.EndElement{Name: t3.Name})
				}
			case reflect.Int,
				reflect.Int8,
				reflect.Int16,
				reflect.Int32,
				reflect.Int64,
				reflect.Uint,
				reflect.Uint8,
				reflect.Uint16,
				reflect.Uint32,
				reflect.Uint64:
				if show {
					t3 = xml.StartElement{
						Name: xml.Name{
							Space: "",
							Local: "int",
						},
					}
					tokens.data = append(tokens.data, t3)
				}
				tokens.recursiveEncode(inner)
				if show {
					tokens.data = append(tokens.data, xml.EndElement{Name: t3.Name})
				}
			case reflect.Float32,
				reflect.Float64:
				if show {
					t3 = xml.StartElement{
						Name: xml.Name{
							Space: "",
							Local: "float",
						},
					}
					tokens.data = append(tokens.data, t3)
				}
				tokens.recursiveEncode(inner)
				if show {
					tokens.data = append(tokens.data, xml.EndElement{Name: t3.Name})
				}
			//case reflect
			default:
				tokens.recursiveEncode(inner)
			}
			tokens.data = append(tokens.data, xml.EndElement{Name: t2.Name})
		}
		tokens.data = append(tokens.data, xml.EndElement{Name: t.Name})
	case reflect.String:
		content := v.String()
		tokens.data = append(tokens.data, xml.CharData(content))
	case reflect.Bool:
		content := strconv.FormatBool(v.Bool())
		tokens.data = append(tokens.data, xml.CharData(content))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		content := xml.CharData(strconv.FormatInt(v.Int(), 10))
		tokens.data = append(tokens.data, xml.CharData(content))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		content := xml.CharData(strconv.FormatUint(v.Uint(), 10))
		tokens.data = append(tokens.data, xml.CharData(content))
	case reflect.Float32, reflect.Float64:
		content := fmt.Sprintf("%f", v.Float())
		tokens.data = append(tokens.data, xml.CharData(content))
	}
}

func (tokens *tokenData) startEnvelope() {
	e := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Envelope",
		},
		Attr: []xml.Attr{
			{Name: xml.Name{Space: "", Local: "xmlns:xsi"}, Value: "http://www.w3.org/2001/XMLSchema-instance"},
			{Name: xml.Name{Space: "", Local: "xmlns:xsd"}, Value: "http://www.w3.org/2001/XMLSchema"},
			{Name: xml.Name{Space: "", Local: "xmlns:soap"}, Value: "http://schemas.xmlsoap.org/soap/envelope/"},
		},
	}

	tokens.data = append(tokens.data, e)
}

func (tokens *tokenData) endEnvelope() {
	e := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Envelope",
		},
	}

	tokens.data = append(tokens.data, e)
}

func (tokens *tokenData) startHeader(m, n string) {
	h := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Header",
		},
	}

	if m == "" || n == "" {
		tokens.data = append(tokens.data, h)
		return
	}

	r := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
		Attr: []xml.Attr{
			{Name: xml.Name{Space: "", Local: "xmlns"}, Value: n},
		},
	}

	tokens.data = append(tokens.data, h, r)

	return
}

func (tokens *tokenData) endHeader(m string) {
	h := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Header",
		},
	}

	if m == "" {
		tokens.data = append(tokens.data, h)
		return
	}

	r := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
	}

	tokens.data = append(tokens.data, r, h)
}

// startToken initiate body of the envelope
func (tokens *tokenData) startBody(m, n string) error {
	b := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Body",
		},
	}

	if m == "" || n == "" {
		return fmt.Errorf("method or namespace is empty")
	}

	r := xml.StartElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
		Attr: []xml.Attr{
			{Name: xml.Name{Space: "", Local: "xmlns"}, Value: n},
		},
	}

	tokens.data = append(tokens.data, b, r)

	return nil
}

// endToken close body of the envelope
func (tokens *tokenData) endBody(m string) {
	b := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: "soap:Body",
		},
	}

	r := xml.EndElement{
		Name: xml.Name{
			Space: "",
			Local: m,
		},
	}

	tokens.data = append(tokens.data, r, b)
}

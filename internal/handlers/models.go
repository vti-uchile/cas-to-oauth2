package handlers

import "encoding/xml"

type CommonResponse struct {
	XMLName xml.Name `xml:"cas:serviceResponse"`
	XMLNS   string   `xml:"xmlns:cas,attr"`
}

type CommonFailure struct {
	Code        string `xml:"code,attr"`
	Description string `xml:",chardata"`
}

type CASResponse struct {
	CommonResponse
	Success *AuthenticationSuccess `xml:"cas:authenticationSuccess,omitempty"`
	Failure *AuthenticationFailure `xml:"cas:authenticationFailure,omitempty"`
}

type ProxyResponse struct {
	CommonResponse
	Success *ProxySuccess `xml:"cas:proxySuccess,omitempty"`
	Failure *ProxyFailure `xml:"cas:proxyFailure,omitempty"`
}

type ProxyValidateResponse struct {
	CommonResponse
	Success *ProxyValidateSuccess `xml:"cas:authenticationSuccess,omitempty"`
	Failure *ProxyValidateFailure `xml:"cas:authenticationFailure,omitempty"`
}

type AuthenticationSuccess struct {
	User string `xml:"cas:user"`
}

type ProxySuccess struct {
	ProxyTicket string `xml:"cas:proxyTicket"`
}

type ProxyValidateSuccess struct {
	User                string   `xml:"cas:user"`
	ProxyGrantingTicket string   `xml:"cas:proxyGrantingTicket"`
	Proxies             []string `xml:"cas:proxies>proxy"`
}

type AuthenticationFailure CommonFailure
type ProxyFailure CommonFailure
type ProxyValidateFailure CommonFailure

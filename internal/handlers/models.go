package handlers

import (
	"encoding/xml"
	"time"
)

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

type SAMLRequest struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Header  Header   `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`
	Body    Body     `xml:"Body"`
}

type Header struct{}

type Body struct {
	Request Request `xml:"Request"`
}

type Request struct {
	MajorVersion      int       `xml:"MajorVersion,attr"`
	MinorVersion      int       `xml:"MinorVersion,attr"`
	RequestID         string    `xml:"RequestID,attr"`
	IssueInstant      time.Time `xml:"IssueInstant,attr"`
	AssertionArtifact string    `xml:"AssertionArtifact"`
}

type AuthenticationFailure CommonFailure
type ProxyFailure CommonFailure
type ProxyValidateFailure CommonFailure

type ResponseEnvelope struct {
	XMLName xml.Name     `xml:"SOAP-ENV:Envelope"`
	XMLNS   string       `xml:"xmlns:SOAP-ENV,attr"`
	Body    ResponseBody `xml:"SOAP-ENV:Body"`
}

type ResponseBody struct {
	Response SAMLResponse
}

type SAMLResponse struct {
	XMLName      xml.Name  `xml:"saml1p:Response"`
	XMLNS        string    `xml:"xmlns:saml1p,attr"`
	ResponseID   string    `xml:"ResponseID,attr"`
	Recipient    string    `xml:"Recipient,attr"`
	MajorVersion int       `xml:"MajorVersion,attr"`
	MinorVersion int       `xml:"MinorVersion,attr"`
	IssueInstant time.Time `xml:"IssueInstant,attr"`
	Status       Status    `xml:"saml1p:Status"`
	Assertion    Assertion `xml:"saml1:Assertion"`
}

type Status struct {
	StatusCode StatusCode `xml:"saml1p:StatusCode"`
}

type StatusCode struct {
	Value string `xml:"Value,attr"`
}

type Assertion struct {
	AssertionID             string                  `xml:"AssertionID,attr"`
	XMLNS                   string                  `xml:"xmlns:saml1,attr"`
	Issuer                  string                  `xml:"Issuer,attr"`
	IssueInstant            time.Time               `xml:"IssueInstant,attr"`
	MajorVersion            int                     `xml:"MajorVersion,attr"`
	MinorVersion            int                     `xml:"MinorVersion,attr"`
	Conditions              Conditions              `xml:"saml1:Conditions,omitempty"`
	AuthenticationStatement AuthenticationStatement `xml:"saml1:AuthenticationStatement,omitempty"`
}

type Conditions struct {
	NotBefore                    time.Time                    `xml:"NotBefore,attr,omitempty"`
	NotOnOrAfter                 time.Time                    `xml:"NotOnOrAfter,attr,omitempty"`
	AudienceRestrictionCondition AudienceRestrictionCondition `xml:"saml1:AudienceRestrictionCondition,omitempty"`
}

type AudienceRestrictionCondition struct {
	Audience string `xml:"saml1:Audience,omitempty"`
}

type AuthenticationStatement struct {
	AuthenticationMethod  string    `xml:"AuthenticationMethod,attr"`
	AuthenticationInstant time.Time `xml:"AuthenticationInstant,attr"`
	Subject               Subject   `xml:"saml1:Subject"`
}

type Subject struct {
	NameIdentifier      string              `xml:"saml1:NameIdentifier"`
	SubjectConfirmation SubjectConfirmation `xml:"saml1:SubjectConfirmation"`
}

type SubjectConfirmation struct {
	ConfirmationMethod string `xml:"saml1:ConfirmationMethod"`
}

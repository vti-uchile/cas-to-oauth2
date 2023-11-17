package handlers

import (
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyValidateHandler(c *gin.Context) {
	ticket := c.DefaultQuery("ticket", "")
	service := c.DefaultQuery("service", "")

	var response ProxyValidateResponse
	response.XMLNS = constants.XML_CAS_NAMESPACE

	if ticket == "" || service == "" {
		response.Failure = &ProxyValidateFailure{
			Code:        "INVALID_REQUEST",
			Description: "'ticket' and 'service' parameters are both required",
		}
		c.XML(http.StatusBadRequest, response)
		return
	}

	// Validate the ticket
	user, pgt, proxies, err := utils.ValidateProxyTicket(ticket, service)
	if err != nil {
		response.Failure = &ProxyValidateFailure{
			Code:        "INTERNAL_ERROR",
			Description: "An internal error occurred during ticket validation",
		}
		c.XML(http.StatusInternalServerError, response)
		return
	}

	if user == "" {
		response.Failure = &ProxyValidateFailure{
			Code:        "INVALID_TICKET",
			Description: "Ticket not recognized",
		}
		c.XML(http.StatusUnauthorized, response)
		return
	}

	response.Success = &ProxyValidateSuccess{
		User:                user,
		ProxyGrantingTicket: pgt,
		Proxies:             proxies,
	}
	c.XML(http.StatusOK, response)
}

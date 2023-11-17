package handlers

import (
	"cas-to-oauth2/constants"
	"cas-to-oauth2/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Proxy(c *gin.Context) {
	pgt := c.DefaultQuery("pgt", "")
	targetService := c.DefaultQuery("targetService", "")
	var response ProxyResponse
	response.XMLNS = constants.XML_CAS_NAMESPACE

	if pgt == "" || targetService == "" {
		c.XML(http.StatusBadRequest, ProxyResponse{
			Failure: &ProxyFailure{
				Code:        "INVALID_REQUEST",
				Description: "'pgt' and 'targetService' parameters are both required",
			},
		})
		return
	}

	isValidPGT, err := utils.ValidatePGT(pgt)
	if err != nil {
		c.XML(http.StatusInternalServerError, ProxyResponse{
			Failure: &ProxyFailure{
				Code:        "INTERNAL_ERROR",
				Description: "An internal error occurred during ticket validation",
			},
		})
		return
	}

	if !isValidPGT {
		c.XML(http.StatusUnauthorized, ProxyResponse{
			Failure: &ProxyFailure{
				Code:        "BAD_PGT",
				Description: "The pgt provided was invalid",
			},
		})
		return
	}

	proxyTicket := utils.GenerateProxyTicket(pgt, targetService)
	if proxyTicket == "" {
		c.XML(http.StatusInternalServerError, ProxyResponse{
			Failure: &ProxyFailure{
				Code:        "INTERNAL_ERROR",
				Description: "An internal error occurred during proxy ticket generation",
			},
		})
		return
	}

	c.XML(http.StatusOK, ProxyResponse{
		Success: &ProxySuccess{
			ProxyTicket: proxyTicket,
		},
	})
}

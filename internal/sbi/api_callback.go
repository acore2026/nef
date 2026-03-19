package sbi

import (
	"net/http"

	"github.com/free5gc/nef/internal/logger"
	"github.com/acore2026/openapi"
	"github.com/acore2026/openapi/models"
	"github.com/acore2026/util/metrics/sbi"
	"github.com/gin-gonic/gin"
)

func (s *Server) getCallbackRoutes() []Route {
	return []Route{
		{
			Method:  http.MethodPost,
			Pattern: "/notification/smf",
			APIFunc: s.apiPostSmfNotification,
		},
	}
}

func (s *Server) apiPostSmfNotification(gc *gin.Context) {
	var eeNotif models.NsmfEventExposureNotification
	reqBody, err := gc.GetRawData()
	if err != nil {
		logger.SBILog.Errorf("Get Request Body error: %+v", err)
		pd := openapi.ProblemDetailsSystemFailure(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusInternalServerError, pd)
		return
	}

	err = openapi.Deserialize(&eeNotif, reqBody, "application/json")
	if err != nil {
		logger.SBILog.Errorf("Deserialize Request Body error: %+v", err)
		pd := openapi.ProblemDetailsMalformedReqSyntax(err.Error())
		gc.Set(sbi.IN_PB_DETAILS_CTX_STR, pd.Cause)
		gc.JSON(http.StatusBadRequest, pd)
		return
	}

	s.Processor().SmfNotification(gc, &eeNotif)
}

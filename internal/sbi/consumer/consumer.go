package consumer

import (
	"github.com/free5gc/nef/internal/logger"
	"github.com/free5gc/nef/pkg/app"
	"github.com/acore2026/openapi"
	"github.com/acore2026/openapi/nrf/NFDiscovery"
	"github.com/acore2026/openapi/nrf/NFManagement"
	"github.com/acore2026/openapi/pcf/PolicyAuthorization"
	"github.com/acore2026/openapi/udr/DataRepository"
)

type nef interface {
	app.App
}

type Consumer struct {
	nef

	// consumer services
	*nnrfService
	*npcfService
	*nudrService
}

func NewConsumer(nef nef) (*Consumer, error) {
	c := &Consumer{
		nef: nef,
	}

	c.nnrfService = &nnrfService{
		consumer:        c,
		nfDiscClients:   make(map[string]*NFDiscovery.APIClient),
		nfMngmntClients: make(map[string]*NFManagement.APIClient),
	}

	c.npcfService = &npcfService{
		consumer: c,
		clients:  make(map[string]*PolicyAuthorization.APIClient),
	}

	c.nudrService = &nudrService{
		consumer: c,
		clients:  make(map[string]*DataRepository.APIClient),
	}
	return c, nil
}

func handleAPIServiceNoResponse(err error) (int, interface{}) {
	detail := "server no response"
	if err != nil {
		detail = err.Error()
	}
	logger.ConsumerLog.Errorf("APIService error: %s", detail)
	pd := openapi.ProblemDetailsSystemFailure(detail)
	return int(pd.Status), pd
}

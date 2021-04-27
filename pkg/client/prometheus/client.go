package prometheus

import (
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

type Client struct {
	v1.API
}

func New(address string) Client {
	promClient, err := api.NewClient(api.Config{Address: address})
	if err != nil {
		panic(err.Error())
	}
	prometheusAPI := v1.NewAPI(promClient)
	return Client{prometheusAPI}
}

package prometheus

import (
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
)

type Options struct {
	Address string
}

type Client struct {
	v1.API
}

func New(options Options) Client {
	promClient, err := api.NewClient(api.Config{Address: options.Address})
	if err != nil {
		panic(err.Error())
	}
	prometheusAPI := v1.NewAPI(promClient)
	return Client{prometheusAPI}
}

package api

import (
	"github.com/MQasimSarfraz/metrics-aggregator/internal/aggregator"
	"github.com/flachnetz/startup/lib/httputil"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type API struct {
	aggregator *aggregator.Aggregator
}

func NewAPI(aggregator *aggregator.Aggregator) API {
	return API{aggregator}
}

type AggregatorOpts struct {
	NodeName string `path:"nodeName"`
}

type AggregatorView struct {
	Timeslice  float32 `json:"timeslice" validate:"required,gte=0"`
	UsedCPU    float32 `json:"cpu" validate:"required,gte=0"`
	UsedMemory float32 `json:"mem" validate:"required,gte=0"`
}

func routing(api API) http.Handler {
	router := httprouter.New()
	register(router, api)
	return router
}

func register(r *httprouter.Router, api API) {
	r.GET("/ping", api.ping)
	r.POST("/v1/metrics/node/:nodeName", api.createNodeMetric)
}

func (api API) ping(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	httputil.WriteResponseValue(w, "pong", nil)
}

func (api API) createNodeMetric(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var opts AggregatorOpts
	var body AggregatorView

	httputil.ExtractAndCallWithBody(&opts, &body, w, req, params, func() (interface{}, error) {

		nodeMetrics := aggregator.NodeMetric{
			UsedCPU:   body.UsedCPU,
			UsedMem:   body.UsedMemory,
			Timeslice: body.Timeslice,
		}
		err := api.aggregator.StoreNodeMetric(opts.NodeName, nodeMetrics)

		if err != nil {
			return nil, err
		}
		return nil, nil
	})

}

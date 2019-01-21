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
	NodeName    string `path:"nodeName"`
	ProcessName string `path:"processName"`
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
	r.POST("/v1/metrics/node/:nodeName/process/:processName", api.createProcessMetric)
}

// ping is used as a health service
func (api API) ping(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	httputil.WriteResponseValue(w, "pong", nil)
}

// createNodeMetric handler is used for storing node metrics
func (api API) createNodeMetric(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var opts AggregatorOpts
	var body AggregatorView

	httputil.ExtractAndCallWithBody(&opts, &body, w, req, params, func() (interface{}, error) {

		nodeMetrics := aggregator.NodeMetric{
			Name:      opts.NodeName,
			UsedCPU:   body.UsedCPU,
			UsedMem:   body.UsedMemory,
			Timeslice: body.Timeslice,
		}
		err := api.aggregator.StoreNodeMetric(nodeMetrics)

		if err != nil {
			return nil, err
		}
		return httputil.WriteJSON(w, http.StatusCreated, body), nil
	})

}

// createProcessMetric handler is used to store process metric
func (api API) createProcessMetric(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
	var opts AggregatorOpts
	var body AggregatorView

	httputil.ExtractAndCallWithBody(&opts, &body, w, req, params, func() (interface{}, error) {

		processMetrics := aggregator.ProcessMetric{
			Name:      opts.ProcessName,
			UsedCPU:   body.UsedCPU,
			UsedMem:   body.UsedMemory,
			Timeslice: body.Timeslice,
			NodeName:  opts.NodeName,
		}
		err := api.aggregator.StoreProcessMetric(processMetrics)

		if err != nil {
			return nil, err
		}

		return httputil.WriteJSON(w, http.StatusCreated, body), nil
	})

}

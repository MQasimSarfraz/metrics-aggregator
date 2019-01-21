// Package aggregate provides primitives for storing/retrieving metrics from database.
package aggregator

import (
	"encoding/json"
	"fmt"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"time"
)

const (
	DbName      = "aggregator"
	NanoSeconds = 1000000000
)

// Aggregator is the main domain type.
// It takes InfluxDB client as a dependency.
type Aggregator struct {
	dbClient client.Client
}

// NewAggregator is used to perform metrics operations.
// It returns a aggregator.Aggregator pointer.
func NewAggregator(dbClient client.Client) *Aggregator {
	return &Aggregator{dbClient}
}

// NodeMetric is a type for metrics collected for a node.
type NodeMetric struct {
	Name      string  `json:"name,omitempty"`
	UsedCPU   float64 `json:"cpu_used"`
	UsedMem   float64 `json:"mem_used"`
	Timeslice float64 `json:"timeslice"`
}

// ProcessMetric is a type for metrics collected for a process.
type ProcessMetric struct {
	Name      string
	UsedCPU   float64
	UsedMem   float64
	Timeslice float64
	NodeName  string
}

// StoreNodeMetric stores the NodeMetric in the influxDB.
// It returns error in case there was an issue writing to database.
func (a *Aggregator) StoreNodeMetric(metric NodeMetric) error {

	batchPoints, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: DbName})
	dbAggregator := map[string]interface{}{
		"cpu_used":  metric.UsedCPU,
		"mem_used":  metric.UsedMem,
		"timeslice": metric.Timeslice,
	}
	point, _ := client.NewPoint("nodes", map[string]string{"name": metric.Name}, dbAggregator, time.Now())
	batchPoints.AddPoint(point)

	err := a.dbClient.Write(batchPoints)
	if err != nil {
		errors.WithMessage(err, "error storing node metric")
	}

	return nil
}

// StoreProcessMetric stores the Process in the influxDB.
// It returns error in case there was an issue writing to database.
func (a *Aggregator) StoreProcessMetric(metric ProcessMetric) error {

	batchPoints, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: DbName})
	dbAggregator := map[string]interface{}{
		"cpu_used":  metric.UsedCPU,
		"mem_used":  metric.UsedMem,
		"timeslice": metric.Timeslice,
	}
	point, _ := client.NewPoint("processes", map[string]string{"name": metric.Name, "nodes": metric.NodeName}, dbAggregator, time.Now())
	batchPoints.AddPoint(point)

	err := a.dbClient.Write(batchPoints)
	if err != nil {
		errors.WithMessage(err, "error storing process metric")
	}

	return nil
}

// GetNodesAverage returns average for cup/mem for all the nodes
// It returns average from current time to the timeslice provided
func (a *Aggregator) GetNodesAverage(timeslice float64) (NodeMetric, error) {

	var metric NodeMetric
	command := fmt.Sprintf("SELECT cpu_used,mem_used FROM nodes WHERE time > Now() - %dns", int(timeslice*NanoSeconds))
	query := client.NewQuery(command, DbName, "ns")
	resp, err := a.dbClient.Query(query)
	if err != nil || resp.Results[0].Series == nil {
		return metric, errors.WithMessage(err, "error fetching node metrics average")
	}

	metrics := resp.Results[0].Series[0].Values
	var totalCPU float64
	var totalMem float64
	for i := 0; i < len(metrics); i++ {
		metricCPU, _ := metrics[i][1].(json.Number).Float64()
		totalCPU = totalCPU + metricCPU

		metricMem, _ := metrics[i][2].(json.Number).Float64()
		totalMem = totalMem + metricMem
	}

	metric.UsedCPU = totalCPU / float64(int(len(metrics)))
	metric.UsedMem = totalMem / float64(int(len(metrics)))
	metric.Timeslice = timeslice

	return metric, nil
}

// Package aggregate provides primitives for storing/retrieving metrics from database.
package aggregator

import (
	"github.com/influxdata/influxdb/client/v2"
	"github.com/pkg/errors"
	"time"
)

// DbName is the database name used in influxDB.
const DbName = "aggregator"

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
	UsedCPU   float32
	UsedMem   float32
	Timeslice float32
}

// StoreNodeMetric stores the NodeMetric in the influxDB.
// It returns error in case there was an issue writing to database.
func (a *Aggregator) StoreNodeMetric(nodeName string, metric NodeMetric) error {

	batchPoints, _ := client.NewBatchPoints(client.BatchPointsConfig{Database: DbName})
	dbAggregator := map[string]interface{}{
		"cpu_used":  metric.UsedCPU,
		"mem_used":  metric.UsedMem,
		"timeslice": metric.Timeslice,
	}
	point, _ := client.NewPoint("node", map[string]string{"name": nodeName}, dbAggregator, time.Now())
	batchPoints.AddPoint(point)

	err := a.dbClient.Write(batchPoints)
	if err != nil {
		errors.WithMessage(err, "error storing node metric")
	}

	return nil
}

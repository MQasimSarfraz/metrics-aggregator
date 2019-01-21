package main

import (
	"fmt"
	"github.com/MQasimSarfraz/metrics-aggregator/internal/aggregator"
	"github.com/MQasimSarfraz/metrics-aggregator/internal/api"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/jessevdk/go-flags"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func main() {

	var opts struct {
		Address         string        `long:"http-address" default:":3080" description:"Address to listen on."`
		InfluxDBAddress string        `long:"influxdb-address" default:"http://influxdb.monitoring.svc.cluster.local:8086" description:"Address used to connect to db."`
		InfluxDBTimeout time.Duration `long:"influxdb-timeout" default:"3" description:"Duration to connect to db."`
	}

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	dbClient, _ := client.NewHTTPClient(client.HTTPConfig{Addr: opts.InfluxDBAddress})
	defer dbClient.Close()

	setupDb(dbClient, opts.InfluxDBTimeout)

	api.Serve(api.NewAPI(aggregator.NewAggregator(dbClient)), opts.Address)
}

func setupDb(c client.Client, duration time.Duration) {

	_, _, err := c.Ping(duration)
	if err != nil {
		logrus.Fatal("Error connecting to database")
		os.Exit(1)
	}

	// TODO(@mqasimsarfraz) clean this may be?
	// ignore error
	query := client.NewQuery(fmt.Sprintf("CREATE DATABASE %s", aggregator.DbName), aggregator.DbName, "ns")
	c.Query(query)

}

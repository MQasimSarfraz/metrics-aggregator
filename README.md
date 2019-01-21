[![Build Status](https://travis-ci.com/MQasimSarfraz/metrics-aggregator.svg?branch=master)](https://travis-ci.com/MQasimSarfraz/metrics-aggregator)
# Metrics Aggregator 
A simple metric aggregator service
## Installation:
The service can be installed in a Kubernetes environment using:
```
kubectl -f deploy/influxdb/
kubectl -f deploy/
```
## Requirements:

- InfluxDB

## API(s):

POST `/v1/metrics/node/{nodename}/` with payload e.g:
```
{  
   "timeslice":"10",
   "cpu":"0.5",
   "mem":"0.3"
}
```
POST `/v1/metrics/node/{nodename}/process/{processname}` with payload e.g:
```
{  
   "timeslice":"10",
   "cpu":"0.5",
   "mem":"0.3"
}
```
GET `/v1/analytics/nodes/average?timeslice=120` with response e.g:
```
{  
   "timeslice":"120",
   "cpu":"0.5",
   "mem":"0.3"
}
```
this returns the average resources usage in last 2 minutes. 
package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus"
	"flag"
)

var containerName = flag.String("docker container name", "functions", "cn")
var dockerApiVersion = flag.String("docker API version", "1.37", "cn")

func main() {

	//Create a new instance of the foocollector and
	//register it with the prometheus client.
	foo := newIronCollector(*dockerApiVersion, *containerName)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	registry := prometheus.NewRegistry()
	registry.Register(foo)
	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	http.Handle("/metrics", handler)
	log.Info("Beginning to serve on port :8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

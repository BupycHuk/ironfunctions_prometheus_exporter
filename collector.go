package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"github.com/docker/docker/client"
	"context"
	"github.com/docker/docker/api/types"
	"fmt"
	"strings"
	"regexp"
	"github.com/Sirupsen/logrus"
)

var (
	succeededRegexp = regexp.MustCompile(`name\[0m=run\..*\.succeeded`)
	requestsRegexp  = regexp.MustCompile(`name\[0m=run\..*\.requests`)
)

//Define a struct for you collector that contains pointers
//to prometheus descriptors for each metric you wish to expose.
//Note you can also include fields of other types if they provide utility
//but we just won't be exposing them as metrics.
type ironCollector struct {
	requestsMetric   *prometheus.Desc
	failedMetric     *prometheus.Desc
	successMetric    *prometheus.Desc
	containerName    string
	dockerApiVersion string
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newIronCollector(dockerApiVersion, containerName string) *ironCollector {
	return &ironCollector{
		successMetric: prometheus.NewDesc("success_runs_count",
			"Shows success runs count",
			nil, nil,
		),
		requestsMetric: prometheus.NewDesc("requests_count",
			"Shows requests count",
			nil, nil,
		),
		failedMetric: prometheus.NewDesc("failed_task_count",
			"Shows failed tasks count",
			nil, nil,
		),
		dockerApiVersion: dockerApiVersion,
		containerName:    containerName,
	}
}

//Each and every collector must implement the Describe function.
//It essentially writes all descriptors to the prometheus desc channel.
func (collector *ironCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.requestsMetric
	ch <- collector.successMetric
}

//Collect implements required collect function for all promehteus collectors
func (collector *ironCollector) Collect(ch chan<- prometheus.Metric) {

	lines, err := collector.getLogs()
	if err != nil {
		logrus.Errorln("Couldn't get logs from iron functions")
		return
	}

	successRunCount := 0
	requestsCount := 0
	failedCount := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "[34mINFO") {
			if succeededRegexp.MatchString(line) {
				successRunCount++
			} else if requestsRegexp.MatchString(line) {
				requestsCount++
			}
		} else if strings.HasPrefix(line, "[31mERRO") {
			if strings.Contains(line, "Failed to run task") {
				failedCount++
			}
		} else {
			continue
		}
	}

	ch <- prometheus.MustNewConstMetric(collector.successMetric, prometheus.CounterValue, float64(successRunCount))
	ch <- prometheus.MustNewConstMetric(collector.requestsMetric, prometheus.CounterValue, float64(requestsCount))
	ch <- prometheus.MustNewConstMetric(collector.failedMetric, prometheus.CounterValue, float64(failedCount))
}
func (collector *ironCollector) getLogs() ([]string, error) {
	cli, err := client.NewClientWithOpts(client.WithVersion(collector.dockerApiVersion))
	if err != nil {
		panic(err)
	}

	reader, err := cli.ContainerLogs(context.Background(), collector.containerName, types.ContainerLogsOptions{
		ShowStdout: true,
		Follow:     false,
	})
	if err != nil {
		fmt.Printf("error: %v", err)
		return nil, err
	}
	defer reader.Close()
	content, _ := ioutil.ReadAll(reader)
	return strings.Split(string(content), "\n"), nil
}

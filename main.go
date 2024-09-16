//go:build alert

package main

/*
Currently only have implemented half of expr
If have time, using alert manager's design or just use alert manager instead
rules:
  - alert: HighRequestLatency
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
    severity: page
    annotations:
    summary: High request latency
*/
import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

var (
	PROMETHEUS_URL = "http://192.168.68.53:9090"
	SLEEP_INTERVAL = 10 * time.Second
	rules          = []string{"(avg by(instance) (rate(cpu_usage[5m])) * 100)"}
)

func queryPrometheus(rule string) model.Value {
	client, err := api.NewClient(api.Config{
		Address: PROMETHEUS_URL,
	})
	if err != nil {
		slog.Error("Error creating client: %v\n", err)
		return nil
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// query by rules
	r := v1.Range{
		Start: time.Now().Add(-time.Hour),
		End:   time.Now(),
		Step:  time.Minute,
	}
	queryResult, warnings, err := v1api.QueryRange(ctx, rule, r)
	if err != nil {
		slog.Error("Error querying Prometheus: %v\n", err)
		return nil
	}
	if len(warnings) > 0 {
		slog.Warn("%v\n", warnings)
	}
	return queryResult
}

// all the implements after query Prometheus are BS*#@,
// sorry don't have enough time to finish the lexical & syntactic analysis in 1 day
func queryAlertRules(jobChan chan) {
	rule := <- jobChan
	queryResult := queryPrometheus(rule)
	metrics := parseMetrics(queryResult.String(), ".*@")

	var flag = false
	for metric := range metrics {
		if metric > 80 {
			flag = true
			break
		}
	}
	if flag {
		slog.Error("CPU usage over 80% \n")
	}

}

func parseMetrics(s string, reg string) []int {
	var result []int
	var valid = regexp.MustCompile(reg)

	metrics := valid.FindAllStringSubmatch(s, -1)
	for _, m := range metrics {
		str := m[0]

		length := len(str)
		if length < 2 {
			continue
		}
		fmt.Println(str[:length-2])
		if r, err := strconv.ParseInt(str[:length-2], 10, 32); err == nil {
			fmt.Println(r)
			result = append(result, int(r))
		}
	}
	return result
}

func main() {
	slog.Info("alert service started")
	for {

		jobChan := make(chan string, 100)

		for i := 0; i < 8; i++ {
			go queryAlertRules(jobChan)
		}

		for _, rule := range rules {
			jobChan <- rule
		}
		close(jobChan)
		time.Sleep(SLEEP_INTERVAL)
	}
}

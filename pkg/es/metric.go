package es

import "github.com/zeromicro/go-zero/core/metric"

const (
	esNamespace = "es_client"
	esSubsystem = "requests"
)

var (
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: esNamespace,
		Subsystem: esSubsystem,
		Name:      "duration_ms",
		Help:      "es client requests duration(ms).",
		Labels:    []string{"index"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricClientReqErrTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: esNamespace,
		Subsystem: esSubsystem,
		Name:      "error_total",
		Help:      "es client requests error count.",
		Labels:    []string{"index", "is_error"},
	})
)

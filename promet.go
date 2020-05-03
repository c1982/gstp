package main

import (
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	countsubject *prometheus.CounterVec
)

func init() {
	NewPrometheusMetrics()
}

func NewPrometheusMetrics() {

	countsubject = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gstp_counter",
			Help: "subject count of filtered email message",
		}, []string{"subject"})

	prometheus.MustRegister(countsubject)
}

func run(svc *GmailService, cfg *GstpConfig) {

	watchFilters(svc, cfg)

	http.Handle(cfg.WebPath, promhttp.Handler())
	log.Println("starting server at " + cfg.Port)
	log.Fatal(http.ListenAndServe(cfg.Port, nil))
}

func watchFilters(srv *GmailService, config *GstpConfig) {
	go func(svc *GmailService, cfg *GstpConfig) {
		for range time.Tick(cfg.CheckInterval) {
			subjects, err := svc.ReadSubjects(cfg.UserID)
			if err != nil {
				log.Fatal(err)
			}
			executeFilters(cfg.Filters, subjects)
		}
	}(srv, config)
}

func executeFilters(filters []GstpFilter, subjects []string) {
	for i := 0; i < len(filters); i++ {
		f := filters[i]

		for j := 0; j < len(subjects); j++ {

			s := subjects[j]

			ok, err := regexp.MatchString(f.SubjectRegex, s)
			if err != nil {
				log.Printf("match error: %s\n", err)
				continue
			}

			if ok {
				countsubject.WithLabelValues(f.Label).Inc()
				log.Printf("increments the %s label by 1. subject: %s\n", f.Label, s)
			}
		}
	}
}

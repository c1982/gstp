package main

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func Test_Filters(t *testing.T) {

	testfilters := []GstpFilter{
		{SubjectRegex: `^\[SOCRadar Incident\].+`, Label: "socradar_incidend"},
		{SubjectRegex: `mysql backup`, Label: "mysql_backup"},
	}

	testsubjects := []string{"[SOCRadar Incident] Last action",
		"mysql backup",
		"spam spam"}

	countsubject.Reset()

	executeFilters(testfilters, testsubjects)

	for _, filter := range testfilters {
		c := testutil.ToFloat64(countsubject.WithLabelValues(filter.Label))
		if c != 1 {
			t.Errorf("want: 1, got: %f", c)
		}
	}
}

// Copyright 2023 RetailNext, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package puppetreport

import "github.com/prometheus/client_golang/prometheus"

var (
	catalogVersionDesc = prometheus.NewDesc(
		"puppet_last_catalog_version",
		"The version of the last attempted Puppet catalog.",
		[]string{"version"},
		nil,
	)

	runAtDesc = prometheus.NewDesc(
		"puppet_last_run_at_seconds",
		"Time of the last Puppet run.",
		nil,
		nil,
	)

	runDurationDesc = prometheus.NewDesc(
		"puppet_last_run_duration_seconds",
		"Duration of the last Puppet run.",
		nil,
		nil,
	)

	runSuccessDesc = prometheus.NewDesc(
		"puppet_last_run_success",
		"1 if the last Puppet run was successful.",
		nil,
		nil,
	)

	resourcesTotalDesc = prometheus.NewDesc(
		"puppet_last_catalog_resources_total",
		"Resources managed during the last Puppet run",
		nil,
		nil,
	)

	resourcesStateDesc = prometheus.NewDesc(
		"puppet_last_catalog_resources",
		"Resource states encountered during the last Puppet run",
		[]string{"state"},
		nil,
	)

	changesTotalDesc = prometheus.NewDesc(
		"puppet_last_catalog_changes_total",
		"Applied node changes during the last Puppet run",
		nil,
		nil,
	)

	eventsTotalDesc = prometheus.NewDesc(
		"puppet_last_catalog_events_total",
		"Events fired during the last Puppet run",
		nil,
		nil,
	)

	eventsStateDesc = prometheus.NewDesc(
		"puppet_last_catalog_events",
		"Events states encountered during the last Puppet run",
		[]string{"state"},
		nil,
	)
)

type Collector struct {
	Logger     Logger
	ReportPath string
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- catalogVersionDesc
	ch <- runAtDesc
	ch <- runDurationDesc
	ch <- runSuccessDesc
	ch <- resourcesTotalDesc
	ch <- resourcesStateDesc
	ch <- changesTotalDesc
	ch <- eventsTotalDesc
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	var result interpretedReport
	if report, err := load(c.reportPath()); err != nil {
		c.Logger.Errorw("puppet_read_run_report_failed", "err", err)
	} else {
		result = report.interpret()
	}
	result.collect(ch)
}

func (c Collector) reportPath() string {
	if c.ReportPath != "" {
		return c.ReportPath
	}
	return "/opt/puppetlabs/puppet/cache/state/last_run_report.yaml"
}

type Logger interface {
	Errorw(msg string, keysAndValues ...interface{})
}

type interpretedReport struct {
	RunAt          float64
	RunDuration    float64
	CatalogVersion string
	RunSuccess     float64
	ResourceCount  float64
	ChangeCount    float64
	EventCount     float64
	ResourceStates map[string]float64
	EventStates    map[string]float64
}

func (r interpretedReport) collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(catalogVersionDesc, prometheus.GaugeValue, 1, r.CatalogVersion)
	ch <- prometheus.MustNewConstMetric(runAtDesc, prometheus.GaugeValue, r.RunAt)
	ch <- prometheus.MustNewConstMetric(runDurationDesc, prometheus.GaugeValue, r.RunDuration)
	ch <- prometheus.MustNewConstMetric(runSuccessDesc, prometheus.GaugeValue, r.RunSuccess)
	ch <- prometheus.MustNewConstMetric(resourcesTotalDesc, prometheus.GaugeValue, r.ResourceCount)
	ch <- prometheus.MustNewConstMetric(changesTotalDesc, prometheus.GaugeValue, r.ChangeCount)
	ch <- prometheus.MustNewConstMetric(eventsTotalDesc, prometheus.GaugeValue, r.EventCount)

	for state, count := range r.ResourceStates {
		ch <- prometheus.MustNewConstMetric(resourcesStateDesc, prometheus.GaugeValue, count, state)
	}

	for state, count := range r.EventStates {
		ch <- prometheus.MustNewConstMetric(eventsStateDesc, prometheus.GaugeValue, count, state)
	}
}

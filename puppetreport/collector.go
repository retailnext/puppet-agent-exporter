// Copyright 2021 RetailNext, Inc.
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
}

func (r interpretedReport) collect(ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(catalogVersionDesc, prometheus.UntypedValue, 0, r.CatalogVersion)
	ch <- prometheus.MustNewConstMetric(runAtDesc, prometheus.GaugeValue, r.RunAt)
	ch <- prometheus.MustNewConstMetric(runDurationDesc, prometheus.GaugeValue, r.RunDuration)
	ch <- prometheus.MustNewConstMetric(runSuccessDesc, prometheus.GaugeValue, r.RunSuccess)
}

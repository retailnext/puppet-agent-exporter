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

import (
	"os"
	"strconv"
	"time"

	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
)

type runReport struct {
	ConfigurationVersion string                      `yaml:"configuration_version"`
	Time                 time.Time                   `yaml:"time"`
	TransactionCompleted bool                        `yaml:"transaction_completed"`
	ReportFormat         int                         `yaml:"report_format"`
	ResourceStatuses     map[string]resourceStatus   `yaml:"resource_statuses"`
	Metrics              map[string]puppetUtilMetric `yaml:"metrics"`
	Logs                 []puppetUtilLog             `yaml:"logs"`
}

func (r runReport) interpret() interpretedReport {
	result := interpretedReport{
		RunAt:                   asUnixSeconds(r.Time),
		RunDuration:             -1,
		CatalogVersion:          r.ConfigurationVersion,
		ConfigRetrievalDuration: -1,
	}
	if r.success() {
		result.RunSuccess = 1
	}

	resourceMetrics, ok := r.Metrics["resources"]
	if ok {
		interpretResourceMetrics(resourceMetrics.Values(), &result)
	}

	timeMetrics, ok := r.Metrics["time"]
	if ok {
		interpretTimeMetrics(timeMetrics.Values(), &result)
	}

	changeMetrics, ok := r.Metrics["changes"]
	if ok {
		interpretChangeMetrics(changeMetrics.Values(), &result)
	}

	eventMetrics, ok := r.Metrics["events"]
	if ok {
		interpretEventMetrics(eventMetrics.Values(), &result)
	}

	return result
}

func interpretResourceMetrics(m map[string]float64, r *interpretedReport) {
	r.ResourceStates = make(map[string]float64, len(m))

	for l, v := range m {
		if l == "total" {
			r.ResourceCount = v
		} else {
			r.ResourceStates[l] = v
		}
	}
}

func interpretTimeMetrics(m map[string]float64, r *interpretedReport) {
	total, ok := m["total"]
	if ok {
		r.RunDuration = total
	}
	config_retrieval, ok := m["config_retrieval"]
	if ok {
		r.ConfigRetrievalDuration = config_retrieval
	}
}

func interpretChangeMetrics(m map[string]float64, r *interpretedReport) {
	total, ok := m["total"]
	if ok {
		r.ChangeCount = total
	}
}

func interpretEventMetrics(m map[string]float64, r *interpretedReport) {
	r.EventStates = make(map[string]float64, len(m))

	for l, v := range m {
		if l == "total" {
			r.EventCount = v
		} else {
			r.EventStates[l] = v
		}
	}
}

func asUnixSeconds(t time.Time) float64 {
	return float64(t.Unix()) + (float64(t.Nanosecond()) / 1e+9)
}

func (r runReport) success() bool {
	if !r.TransactionCompleted {
		return false
	}
	var failed int
	for _, item := range r.ResourceStatuses {
		if item.Failed {
			failed++
		}
	}
	return failed == 0
}

type resourceStatus struct {
	Failed         bool    `yaml:"failed"`
	EvaluationTime float64 `yaml:"evaluation_time"`
}

type puppetUtilMetric struct {
	Name      string     `yaml:"name"`
	Label     string     `yaml:"label"`
	RawValues [][]string `yaml:"values"`
}

func (s puppetUtilMetric) Values() map[string]float64 {
	result := make(map[string]float64, len(s.RawValues))
	for _, item := range s.RawValues {
		if len(item) == 3 {
			value, err := strconv.ParseFloat(item[2], 64)
			if err == nil {
				result[item[0]] = value
			}
		}
	}
	return result
}

type puppetUtilLog struct {
	Time time.Time `yaml:"time"`
}

func load(path string) (runReport, error) {
	file, err := os.Open(path)
	if err != nil {
		return runReport{}, err
	}

	decoder := yaml.NewDecoder(file)
	var report runReport
	err = decoder.Decode(&report)
	return report, multierr.Append(err, file.Close())
}

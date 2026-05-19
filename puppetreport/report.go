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

	"github.com/retailnext/puppet-agent-exporter/internal/utils"
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
		RunAt:          utils.UnixSeconds(r.Time),
		RunDuration:    r.totalDuration(),
		CatalogVersion: r.ConfigurationVersion,
	}
	if r.success() {
		result.RunSuccess = 1
	}
	return result
}

func (r runReport) totalDuration() float64 {
	timeMetrics, ok := r.Metrics["time"]
	if !ok {
		return -1
	}
	values := timeMetrics.Values()
	total, ok := values["total"]
	if !ok {
		return -1
	}
	return total
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

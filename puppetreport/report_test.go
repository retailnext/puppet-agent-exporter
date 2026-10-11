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
	"fmt"
	"testing"
)

func TestLoadReport(t *testing.T) {
	testCases := map[string]interpretedReport{
		"last_run_report": interpretedReport{
			RunAt:          1618957125.5901103,
			RunDuration:    17.199882286,
			CatalogVersion: "1618957129",
			RunSuccess:     1,
		},
		"last_run_report-5.4.0": interpretedReport{
			RunAt:          1725776230.602652,
			RunDuration:    0.03196727,
			CatalogVersion: "1725776230",
			RunSuccess:     1,
			ResourceCount:  8,
			ResourceStates: map[string]float64{
				"skipped":           0,
				"failed":            0,
				"failed_to_restart": 0,
				"restarted":         0,
				"changed":           0,
				"out_of_sync":       0,
				"scheduled":         0,
				"corrective_change": 0,
			},
			EventStates: map[string]float64{
				"failure": 0,
				"success": 0,
			},
		},
		"last_run_report-6.28.0": interpretedReport{
			RunAt:          1725776354.854867,
			RunDuration:    0.004820335,
			CatalogVersion: "1725776354",
			RunSuccess:     1,
			ResourceCount:  8,
			ResourceStates: map[string]float64{
				"skipped":           0,
				"failed":            0,
				"failed_to_restart": 0,
				"restarted":         0,
				"changed":           0,
				"out_of_sync":       0,
				"scheduled":         0,
				"corrective_change": 0,
			},
			EventStates: map[string]float64{
				"failure": 0,
				"success": 0,
			},
		},
		"last_run_report-7.32.1": interpretedReport{
			RunAt:          1725776438.356112,
			RunDuration:    0.006013873,
			CatalogVersion: "1725776438",
			RunSuccess:     1,
			ResourceCount:  8,
			ResourceStates: map[string]float64{
				"skipped":           0,
				"failed":            0,
				"failed_to_restart": 0,
				"restarted":         0,
				"changed":           0,
				"out_of_sync":       0,
				"scheduled":         0,
				"corrective_change": 0,
			},
			EventStates: map[string]float64{
				"failure": 0,
				"success": 0,
			},
		},
		"last_run_report-8.8.1": interpretedReport{
			RunAt:          1725776515.039312,
			RunDuration:    0.005837204,
			CatalogVersion: "1725776515",
			RunSuccess:     1,
			ResourceCount:  8,
			ResourceStates: map[string]float64{
				"skipped":           0,
				"failed":            0,
				"failed_to_restart": 0,
				"restarted":         0,
				"changed":           0,
				"out_of_sync":       0,
				"scheduled":         0,
				"corrective_change": 0,
			},
			EventStates: map[string]float64{
				"failure": 0,
				"success": 0,
			},
		},
	}

	for name, tc := range testCases {
		want := tc
		t.Run(name, func(t *testing.T) {
			report, err := load("testdata/" + name + ".yaml")
			if err != nil {
				t.Fatal(err)
			}

			got := report.interpret()

			if want.RunAt != got.RunAt {
				t.Fatalf("RunAt: want %f; got %f", want.RunAt, got.RunAt)
			}

			if want.RunDuration != got.RunDuration {
				t.Fatalf("RunDuration: want %f; got %f", want.RunDuration, got.RunDuration)
			}

			if want.CatalogVersion != got.CatalogVersion {
				t.Fatalf("CatalogVersion: want %q; got %q", want.CatalogVersion, got.CatalogVersion)
			}

			if want.RunSuccess != got.RunSuccess {
				t.Fatalf("RunSuccess: want %f; got %f", want.RunSuccess, got.RunSuccess)
			}

			if want.ResourceCount != got.ResourceCount {
				t.Fatalf("ResourceCount: want %f; got %f", want.ResourceCount, got.ResourceCount)
			}

			if want.ChangeCount != got.ChangeCount {
				t.Fatalf("ChangeCount: want %f; got %f", want.ChangeCount, got.ChangeCount)
			}

			if want.EventCount != got.EventCount {
				t.Fatalf("EventCount: want %f; got %f", want.EventCount, got.EventCount)
			}

			if err := mapCompare(want.ResourceStates, got.ResourceStates); err != "" {
				t.Fatalf("ResourceStates: %s", err)
			}

			if err := mapCompare(want.EventStates, got.EventStates); err != "" {
				t.Fatalf("EventStates: %s", err)
			}
		})
	}
}

func mapCompare(l, r map[string]float64) string {
	if len(l) != len(r) {
		return fmt.Sprintf("length: want %d, got %d", len(l), len(r))
	}

	for k, want := range l {
		got, ok := r[k]
		if !ok {
			return fmt.Sprintf("key %q is missing", k)
		}

		if want != got {
			return fmt.Sprintf("key %q: want %f, got %f", k, want, got)
		}
	}

	return ""
}

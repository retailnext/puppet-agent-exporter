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

import "testing"

func TestLoadReport(t *testing.T) {
	report, err := load("last_run_report.yaml")
	if err != nil {
		t.Fatal(err)
	}

	ir := report.interpret()
	expected := interpretedReport{
		RunAt:          1618957125.5901103,
		RunDuration:    17.199882286,
		CatalogVersion: 1618957129,
		RunSuccess:     1,
	}
	if ir != expected {
		t.Fatalf("%+v != %+v", ir, expected)
	}
}

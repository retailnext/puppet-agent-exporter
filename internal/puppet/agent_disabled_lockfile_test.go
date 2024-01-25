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

package puppet

import (
	"testing"
)

func TestParseAgentDisabledLockfile(t *testing.T) {
	testCases := map[string]struct {
		wantDisabledMessage string
		wantError           bool
	}{
		"404":              {wantError: true},
		"empty":            {wantError: true},
		"malformed":        {wantError: true},
		"empty_hash":       {},
		"disabled_message": {wantDisabledMessage: "testing unmarshalling"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			got, err := ParseAgentDisabledLockfile("testdata/agent_disabled_lockfile/" + name + ".json")

			if tc.wantError {
				if err == nil {
					t.Errorf("expected fixture %q to produce error but got none", name)
				}
			} else {
				if err != nil {
					t.Errorf("expected fixture %q to parse properly but got error\n%s", name, err)
				}
				if got == nil {
					t.Errorf("expected fixture %q to parse properly but got nil as result", name)
				}
			}
		})
	}
}

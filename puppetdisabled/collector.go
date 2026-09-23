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

package puppetdisabled

import (
	"errors"
	"os"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/retailnext/puppet-agent-exporter/internal/logging"
	"github.com/retailnext/puppet-agent-exporter/internal/puppet"
	"github.com/retailnext/puppet-agent-exporter/internal/utils"
)

var disabledSinceDesc = prometheus.NewDesc(
	"puppet_disabled_since_seconds",
	"Time since when puppet has been disabled.",
	nil,
	nil,
)

type Collector struct {
	Logger       logging.Logger
	LockfilePath string
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- disabledSinceDesc
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	dto, err := puppet.ParseAgentDisabledLockfile(c.lockfilePath())
	if errors.Is(err, os.ErrNotExist) {
		return // nothing to report
	} else if err != nil {
		c.Logger.Errorw("puppet_open_lockfile_failed", "err", err)
		return
	}

	disabledSince := utils.UnixSeconds(dto.DisabledSince)
	ch <- prometheus.MustNewConstMetric(disabledSinceDesc, prometheus.GaugeValue, disabledSince)
}

func (c Collector) lockfilePath() string {
	if c.LockfilePath != "" {
		return c.LockfilePath
	}
	return puppet.DefaultAgentDisabledLockfile
}

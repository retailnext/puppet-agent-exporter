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

package main

import (
	"context"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/term"

	"github.com/retailnext/puppet-agent-exporter/puppetconfig"
	"github.com/retailnext/puppet-agent-exporter/puppetdisabled"
	"github.com/retailnext/puppet-agent-exporter/puppetreport"
)

func setupLogger() func() {
	var logger *zap.Logger
	var err error
	if term.IsTerminal(int(os.Stdin.Fd())) {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.RedirectStdLog(logger)

	return func() {
		_ = logger.Sync()
	}
}

func setupInterruptContext() (context.Context, func()) {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		select {
		case sig := <-c:
			zap.S().Infow("shutting_down", "signal", sig)
			cancel()
		case <-ctx.Done():
		}
	}()
	onExit := func() {
		signal.Stop(c)
		cancel()
	}
	return ctx, onExit
}

var rootTemplate = template.Must(template.New("/").Parse(`<html>
<head><title>puppet-agent-exporter</title></head>
<body>
<h1>puppet-agent-exporter</h1>
<p><a href="{{.}}">Metrics</a></p>
</body>
</html>
`))

func run(ctx context.Context, listenAddress, telemetryPath string) (ok bool) {
	lgr := zap.S()

	prometheus.DefaultRegisterer.MustRegister(puppetconfig.Collector{
		Logger: lgr,
	})
	prometheus.DefaultRegisterer.MustRegister(puppetreport.Collector{
		Logger: lgr,
	})
	prometheus.DefaultRegisterer.MustRegister(puppetdisabled.Collector{
		Logger: lgr,
	})

	mux := http.NewServeMux()
	mux.Handle(telemetryPath, promhttp.Handler())
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_ = rootTemplate.Execute(writer, telemetryPath)
	})
	server := &http.Server{Addr: listenAddress, Handler: mux}

	resultCh := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err == http.ErrServerClosed {
			resultCh <- true
			return
		}
		lgr.Errorw("listen_and_serve_failed", "err", err)
		resultCh <- false
	}()

	stopCh := ctx.Done()
	for {
		select {
		case <-stopCh:
			shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 5*time.Second)
			if shutdownErr := server.Shutdown(shutdownCtx); shutdownErr != nil {
				lgr.Warnw("server_shutdown_failed", "err", shutdownErr)
			}
			cancelShutdownCtx()
			ok = <-resultCh
			return
		case ok = <-resultCh:
			return
		}
	}
}

func main() {
	// Flag definitions copied from github.com/prometheus/node_exporter
	var (
		listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9819").String()
		telemetryPath = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics.").Default("/metrics").String()
	)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	syncLogger := setupLogger()
	defer syncLogger()

	ctx, onExit := setupInterruptContext()
	defer onExit()

	if ok := run(ctx, *listenAddress, *telemetryPath); !ok {
		os.Exit(1)
	}
}

# puppet-agent-exporter

## Puppet Agent Prometheus Exporter

This [Prometheus](https://prometheus.io/)
[exporter](https://prometheus.io/docs/instrumenting/exporters/)
exposes the status of the Puppet Agent on the host it is running on.

Unlike [puppet-prometheus_reporter](https://github.com/voxpupuli/puppet-prometheus_reporter)
and other solutions that rely on the Puppet Agent successfully a report to the
Puppet Server, this allows you to actively monitor the status of every node in
your environment that is discoverable by Prometheus.

### Metrics Exposed

```
# HELP puppet_config Puppet configuration.
# TYPE puppet_config gauge
puppet_config{environment="",server="puppet.redacted"} 1
# HELP puppet_last_catalog_version The version of the last attempted Puppet catalog.
# TYPE puppet_last_catalog_version gauge
puppet_last_catalog_version{version="1680640107"} 1
# HELP puppet_last_run_at_seconds Time of the last Puppet run.
# TYPE puppet_last_run_at_seconds gauge
puppet_last_run_at_seconds 1.6806401024160552e+09
# HELP puppet_last_run_duration_seconds Duration of the last Puppet run.
# TYPE puppet_last_run_duration_seconds gauge
puppet_last_run_duration_seconds 28.023470087
# HELP puppet_last_run_success 1 if the last Puppet run was successful.
# TYPE puppet_last_run_success gauge
puppet_last_run_success 1
```

### Example Alert Rules

```yaml
groups:
  - name: Puppet
    rules:
      - alert: MultiplePuppetFailing
        expr: count(puppet_last_run_success) - count(puppet_last_run_success == 1) > 4
      - alert: PuppetEnvironmentSet
        expr: puppet_config{environment!=""}
        for: 4h
      - alert: PuppetFailing
        expr: puppet_last_run_success == 0
        for: 40m
      - alert: LastPuppetTooLongAgo
        expr: time() - puppet_last_run_at_seconds > 3*60*60
        for: 40m
```

We've found it worthwhile to avoid alerting on conditions affecting individual
nodes until that condition has persisted long enough to affect more than one
run. (For example, a transient hiccup in an APT proxy, etc...)

Alerting on the catalog being too old helps catch situations like a node being
set to an environment that no longer exists, or if it's having TLS or network
issues contacting the Puppet Server.

Alerting on a non-default environment being set helps catch operator error,
for example when a node is used to test changes from a branch environment
but forgotten about after that branch is merged.

## Project Status: **Works For Us**

This is an open-source release of something we've used for quite a while
internally.

*   The package APIs functions are likely to be refactored drastically,
    potentially without warning in release notes.

    (This isn't really meant to be a library.)

*   We may add a systemd unit to the packaging at some point.

*   We probably won't make breaking changes to the arguments without warning.

### Areas needing improvement

*   Test coverage
*   General code hygiene and refactoring
*   Better packaging (supporting RPM distros, including a systemd unit, etc...)

## Contributing

Contributions considered, but be aware that this is mostly just something we
needed. It's public because there's no reason anyone else should have to waste
an afternoon (or more) building something similar, and we think the approach
is good enough that others would benefit from adopting.

This project is licensed under the [Apache License, Version 2.0](LICENSE).

Please include a `Signed-off-by` in all commits, per
[Developer Certificate of Origin version 1.1](DCO).

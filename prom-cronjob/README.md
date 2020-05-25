# prom-cronjob

A shell script (and some rules) for cronjob monitoring with Prometheus.  Successful completion of the job updates a file read by the node_exporter textfile collector.

It intentionally redirects the output to `/tmp/${job}.stdout` and
`/tmp/${job}.stderr` so no email is sent by cron.

## Usage

```crontab
*/5 * * * * cronwrap foobar someprogram --foo=bar
```

`foobar` will be the value of the cronjob label on the metric.

## Setup

Enable the node_exporter textfile collector:

```shell
# mkdir /var/spool/nodexporter
# chmod a+rwxt /var/spool/nodeexporter
# node_exporter --collector.textfile.directory=/var/spool/nodeexporter
```

## Rules

```yaml
groups:
- name: cron
  rules:
    # instead of making an entirely new rule for every cronjob, you can
    # just setup a simple threshold value.
    - record: cron_threshold
      labels:
        cronjob: update-duck
      expr: 7200
    - alert: CronJobStale
      expr: time() - last_success > on(cronjob) cron_threshold
      for: 10m
      labels:
        severity: page
        job: cronjob
      annotations:
        summary: >
          no successful {{ $labels.cronjob }}
          in {{ humanizeDuration $value }}!
    # If you want, you can set a default value for jobs without a threshold.
    - alert: CronJobStale
      expr: (time() - last_success > 3600) unless on(cronjob) cron_threshold
      for: 10m
      labels:
        severity: page
        job: cronjob
      annotations:
        summary: no successful {{ $labels.cronjob }} in {{ humanizeDuration $value }}!
```

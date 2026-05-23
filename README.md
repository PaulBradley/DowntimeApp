<img src="./docs/images/icons/go.svg?raw=true"  width="64" height="64" />&nbsp;
<img src="./docs/images/icons/aws.svg?raw=true" width="64" height="64" />

# DowntimeApp

This repository contains the code required to run the SaaS [Downtime Solution](https://downtimeapp.cloud).

A cloud-native downtime solution providing a critical backup solution for when key hospital systems become unavailable due to either planned/un-planned downtime.

<img src="./docs/images/downtime-solution.jpg?raw=true" width="512" height="314" />


## Progress

- [2026-05-23](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-23&until=2026-05-23): Write DSQL endpoints out to .mk file for the schema creation tool
- [2026-05-22](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-22&until=2026-05-22): Refactored S3 & DSQL clients; updated makefile to inject Git commit hash into binary
- [2026-05-21](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-21&until=2026-05-21): Added code to provision/teardown S3 buckets; added flag to just report status of AWS
- [2026-05-19](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-21&until=2026-05-21): Added code to provision & teardown AWS DSQL databases

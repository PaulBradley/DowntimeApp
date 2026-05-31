<img src="./docs/images/icons/go.svg?raw=true"  width="64" height="64" />&nbsp;
<img src="./docs/images/icons/aws.svg?raw=true" width="64" height="64" />

# DowntimeApp

This repository contains the code required to run the SaaS [Downtime Solution](https://downtimeapp.cloud).

A cloud-native downtime solution providing a critical backup solution for when key hospital systems become unavailable due to either planned/un-planned downtime.

<img src="./docs/images/downtime-solution.jpg?raw=true" width="512" height="314" />


---

## Architecture Diagram

Our [cloud infrastructure provisioning tool](https://github.com/PaulBradley/DowntimeApp/tree/main/src/aws-infrastructure) builds the AWS architecture detailed below. This allow health organisations to provision the infrastructure within their own AWS accounts if required.

<img src="./docs/images/architecture-diagram.png?raw=true" width="892" height="388" />


## Progress

- [2026-05-30](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-31&until=2026-05-31): Add `environments` option; extended database schema
- [2026-05-27](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-27&until=2026-05-27): Add `list-tables` option; support multi DDL statements; extended migrations
- [2026-05-26](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-26&until=2026-05-26): Add initial SQL schema migration utility
- [2026-05-23](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-23&until=2026-05-23): Write DSQL endpoints out to .mk file for the schema creation tool
- [2026-05-22](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-22&until=2026-05-22): Refactored S3 & DSQL clients; updated makefile to inject Git commit hash into binary
- [2026-05-21](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-21&until=2026-05-21): Added code to provision/teardown S3 buckets; added flag to just report status of AWS
- [2026-05-19](https://github.com/PaulBradley/DowntimeApp/commits/main?since=2026-05-21&until=2026-05-21): Added code to provision & teardown AWS DSQL databases

# Advanced URL shortener service

## Key features
- For sure, url shortening;
- Metrics - clicks, user data etc;
  - Metrics are stored in Clickhouse and delivered there via event-sourcing by Kafka.

## How to run locally

Clone a project
```shell
    git clone <project repo url>
```
In a root folder start a local docker env

```shell
    docker compose -f ./deployment/local/docker-compose.yml up -d
```
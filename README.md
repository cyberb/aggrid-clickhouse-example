# aggrid-clickhouse-example
AG Grid with ClickHouse Example

### run clickhouse db

```
docker run -d -p 8123:8123 -p9000:9000 --name clickhouse --ulimit nofile=262144:262144 clickhouse/clickhouse-server
```
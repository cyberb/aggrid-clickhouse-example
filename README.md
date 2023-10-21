# aggrid-clickhouse-example
AG Grid with ClickHouse Example

### run clickhouse server

```
docker run -d -p 8123:8123 -p9000:9000 --name clickhouse --ulimit nofile=262144:262144 clickhouse/clickhouse-server
```

### run dev web
```
npm i
npm run dev
```

### run clickhouse client

```
sudo docker run -it --rm --link clickhouse:clickhouse-server --entrypoint clickhouse-client clickhouse/clickhouse-server --host clickhouse-server
```
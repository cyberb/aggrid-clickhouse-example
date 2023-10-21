# AG Grid with ClickHouse Example

![](https://github.com/cyberb/aggrid-clickhouse-example/releases/download/demo/demo.gif)

### run clickhouse server

```
docker run -d -p 8123:8123 -p9000:9000 --name clickhouse --ulimit nofile=262144:262144 clickhouse/clickhouse-server
```

### run api

```
cd api
go build -o api ./cmd
./api
```

### run dev web
```
cd web
npm i
npm run dev
```

### run clickhouse client (optional)

```
sudo docker run -it --rm --link clickhouse:clickhouse-server --entrypoint clickhouse-client clickhouse/clickhouse-server --host clickhouse-server
```
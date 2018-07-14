# Very simple prometheus exporter for iron functions

Run iron functions using:
```bash
docker run --rm -it --name functions -v ${PWD}/data:/app/data -v /var/run/docker.sock:/var/run/docker.sock -p 8010:8080 iron/functions
```

Run our exporter:
```bash
go run main.go collector.go
```

Run prometheus using:
```bash
docker run -p 9090:9090 --net=host -v $PWD/prometheus-data:/prometheus-data prom/prometheus --config.file=/prometheus-data/prometheus.yml
```

Grafana dashboard json located in `./dashboard`

to run a few tasks in iron functions use commands below:

```bash
curl http://localhost:8010/r/myapp/hello #runs failed task
```

```bash
curl http://localhost:8010/r/hello2/hello2 #runs succeeded task
```
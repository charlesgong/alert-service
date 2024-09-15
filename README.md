# Assignment description
collector includes 
4. alert


# Alert Metrics
Query from

## Dependencies check, please refer https://github.com/charlesgong/metrics-service 
Promethues

Metric storage

### Grafana

Visualization, please import sample-grafana.json after startup

## Getting Start
In our Go project, we just need to import the dependent library.

```go mod download```

Attach suitable metrics to relevant functions to evaluate the performance of the scope. Check example code [here](./main.go).

## Step One

Start the server to collect the metrics.

```go run main.go```

and it will query Prometheus metrics itself

## Step Two

Server metrics need to be collected and visualized to analyze the data. Prometheus and Grafana need to be started and configured along with Docker.

### Docker Setup

Install Docker with Docker Compose. Check out [here](https://docs.docker.com/engine/install/)

```bash 
docker build -t [image_name] 

```

Then initiate the Docker Compose by using the following command,









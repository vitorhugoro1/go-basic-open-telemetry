# Go Basic OpenTelemetry

This project is a simple project made just to show a implementation of OpenTelemetry on Golang.

## How to run

To run this project you'll need:

- Docker
- Go 1.23

With this tools installed you'll need to run:

```bash
docker compose up -d

OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4317 go run main.go

# Request API to endpoint on API
curl --location 'http://localhost:8080/up'
```

## Credits

- [Vitor Merencio](https://github.com/vitorhugoro1)

## License

The GNU GENERAL PUBLIC LICENSE (GNU). Please see [License File](LICENSE.md) for more information.

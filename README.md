# Zerolog Loki hook

Example for a hook that forwards log events from zerolog to you trusty local Loki instance.

For usage checkout the main.go file. To test the setup run 
```sh
docker-compose up -d
```

When running the expample directly the loki instance is expected to be reachable at `localhost:3100`. You can overwrite this using the `LOKI_ENDPOINT` environment variable.

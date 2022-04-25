# nuntius

`nuntius` is the feeder accumulating data from several sources and sending it to the [data server](https://github.com/skybi/pluteo)

## Running in production (Docker)

To run `nuntius` in production, we strongly recommend simply using our stable Docker images released to GHCR.

To make the setup more straightforward, an [example `docker-compose.yml`](docker-compose.example.yml) is provided.

## Running in development (local)

To run a local development copy, simply clone the commit/branch you want to run and create a `.env` file in the same
directory you run the `go run` or binary command.

See the [example `.env`](.env.example) for a quick overview.

## Configuration variables

| Environment variable | Type            | Default                 | Description                                                                    |
|----------------------|-----------------|-------------------------|--------------------------------------------------------------------------------|
| `SBF_ENVIRONMENT`    | `prod` or `dev` | `prod`                  | Whether the worker starts in development or production mode                    |
| `SBF_API_ADDRESS`    | `URL`           | `http://localhost:8082` | The URL of the data API to feed the data into                                  |
| `SBF_API_KEY`        | `string`        | `<none>`                | The API key to use for the data API (unlimited quota & rate limit is required) |
| `SBF_FEED_METARS`    | `bool`          | `false`                 | Whether or not to feed METARs                                                  |

# Country Info Service

A REST web service that provides country information, exchange rates for neighbouring countries, and diagnostics for
dependent services.

## Prerequisites

- Go 1.26.0 (or Docker)

## Configuration

The service is configured via environment variables:

| Variable             | Required | Default | Description                                                                 |
|----------------------|----------|---------|-----------------------------------------------------------------------------|
| `PORT`               | No       | `8080`  | Port the HTTP server listens on                                             |
| `COUNTRIES_ENDPOINT`  | Yes      | -       | Base URL for the REST Countries API (e.g. `http://129.241.150.113:8080/v3.1`) |
| `CURRENCY_ENDPOINT`   | Yes      | -       | Base URL for the Currency Exchange API (e.g. `http://129.241.150.113:9090/currency`) |

## Running

### Locally

```sh
export COUNTRIES_ENDPOINT="http://129.241.150.113:8080/v3.1"
export CURRENCY_ENDPOINT="http://129.241.150.113:9090/currency"
export PORT="8080"

go run ./cmd/server
```

### Docker

```sh
docker build -t countryinfo .

docker run -p 8080:8080 \
    -e COUNTRIES_ENDPOINT="http://129.241.150.113:8080/v3.1" \
    -e CURRENCY_ENDPOINT="http://129.241.150.113:9090/currency" \
    -e PORT="8080" \
    countryinfo

```

### Running Tests

```sh
go test ./...
```

## API Endpoints

All endpoints are prefixed with `/countryinfo/v1`.

```
http://localhost:8080/countryinfo/v1/status/
http://localhost:8080/countryinfo/v1/info/{two_letter_country_code}
http://localhost:8080/countryinfo/v1/exchange/{two_letter_country_code}
```

---

### Diagnostics

Reports the health of upstream dependencies and service uptime.

**Request**

```
Method: GET
Path:   /countryinfo/v1/status/
```

**Response**

- Content-Type: `application/json`
- Status: `200` if all upstream services are reachable, `503` if one or more are unavailable.

```json
{
  "restcountriesapi": 200,
  "currenciesapi": 200,
  "version": "v1",
  "uptime": 3600
}
```

| Field              | Type    | Description                                            |
|--------------------|---------|--------------------------------------------------------|
| `restcountriesapi` | integer | HTTP status code returned by the REST Countries API    |
| `currenciesapi`    | integer | HTTP status code returned by the Currency Exchange API |
| `version`          | string  | API version                                            |
| `uptime`           | integer | Seconds since the service was last started             |

**Example**

```sh
curl http://localhost:8080/countryinfo/v1/status/
```

---

### Country Info

Returns general information about a country identified by its two-letter country
code ([ISO 3166-2](https://en.wikipedia.org/wiki/ISO_3166-2)).

**Request**

```
Method: GET
Path:   /countryinfo/v1/info/{two_letter_country_code}
```

| Parameter                 | Description                                     |
|---------------------------|-------------------------------------------------|
| `two_letter_country_code` | ISO 3166-2 country code (e.g. `no`, `se`, `us`) |

**Response**

- Content-Type: `application/json`
- Status: `200` on success, `400` for invalid country code, `502` if the upstream API is unreachable.

```json
{
  "name": "Norway",
  "continents": [
    "Europe"
  ],
  "population": 5379475,
  "area": 323802,
  "languages": {
    "nno": "Norwegian Nynorsk",
    "nob": "Norwegian Bokmal",
    "smi": "Sami"
  },
  "borders": [
    "FIN",
    "SWE",
    "RUS"
  ],
  "flag": "https://flagcdn.com/w320/no.png",
  "capital": "Oslo"
}
```

| Field        | Type             | Description                           |
|--------------|------------------|---------------------------------------|
| `name`       | string           | Common name of the country            |
| `continents` | array of strings | Continents the country belongs to     |
| `population` | integer          | Population count                      |
| `area`       | integer          | Area in km²                           |
| `languages`  | object           | Map of language code to language name |
| `borders`    | array of strings | ISO 3166 codes of bordering countries |
| `flag`       | string           | URL to the country flag image (PNG)   |
| `capital`    | string           | Capital city                          |

**Example**

```sh
curl http://localhost:8080/countryinfo/v1/info/no
```

---

### Exchange Rates

Returns currency exchange rates between the input country and its neighbouring countries.

The service looks up the country by its two-letter code, determines its base currency and bordering countries, then
returns the exchange rates from the base currency to each neighbour's currency.

**Request**

```
Method: GET
Path:   /countryinfo/v1/exchange/{two_letter_country_code}
```

| Parameter                 | Description                                     |
|---------------------------|-------------------------------------------------|
| `two_letter_country_code` | ISO 3166-2 country code (e.g. `no`, `se`, `us`) |

**Response**

- Content-Type: `application/json`
- Status: `200` on success, `400` for invalid country code, `502` if an upstream API is unreachable.

```json
{
  "country": "Norway",
  "base-currency": "NOK",
  "exchange-rates": [
    {
      "EUR": 0.086536
    },
    {
      "SEK": 0.914075
    },
    {
      "RUB": 7.711945
    }
  ]
}
```

| Field            | Type             | Description                                                                                                |
|------------------|------------------|------------------------------------------------------------------------------------------------------------|
| `country`        | string           | Common name of the country                                                                                 |
| `base-currency`  | string           | ISO 4217 currency code of the input country                                                                |
| `exchange-rates` | array of objects | Each object maps a neighbour's currency code (ISO 4217) to its exchange rate relative to the base currency |

If a country has no land borders (e.g. Iceland), `exchange-rates` will be an empty array.

**Example**

```sh
curl http://localhost:8080/countryinfo/v1/exchange/no
```

## Project Structure

```
cmd/server/          Application entrypoint
internal/
  config/            Environment-based configuration
  handler/
    info/            Country info endpoint
    exchange/        Exchange rates endpoint
    status/          Diagnostics endpoint
  middleware/        HTTP middleware (logging, request ID)
  restclient/        HTTP clients for upstream APIs
  router/            Route registration
  server/            HTTP server lifecycle
  fp/                Generic functional programming utilities
  util/              Input validation and URL helpers
```

## Dependencies

This service depends on two external APIs:

- **REST Countries API**
  - Endpoint: http://129.241.150.113:8080/v3.1/
  - Documentation: http://129.241.150.113:8080/
- **Currency Exchange API**
  - Endpoint: http://129.241.150.113:9090/currency/
  - Documentation: http://129.241.150.113:9090/

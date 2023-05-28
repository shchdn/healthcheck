# Healthcheck service

Interview task

## Run server

```bash
make http
```

## Get min

```bash
curl "http://localhost:8080/get_min"
```

## Get max

```bash
curl "http://localhost:8080/get_max"
```

## Get info

```bash
curl "http://localhost:8080/get_info?url=https://instagram.com"
```

## Get visited stats

```bash
`curl "http://localhost:8080/get_request_stats"`
```

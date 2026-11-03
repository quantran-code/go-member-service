# go-member-service

Owns library members: who they are and how to reach them.

```sh
PORT=8081 go run .
```

| Route | Behaviour |
|---|---|
| `GET /health` | `{"status":"ok"}` when the process is up |
| `GET /members/{id}` | The member as `{"id":"1","name":"...","email":"..."}`, or `404` with `{"error":"member not found"}` |

The public path for a member is `GET /api/members/{id}`: the gateway proxies `/api/members/{id}`
to this service as `/members/{id}`. Members are seeded in memory — nothing is persisted.

## Request correlation

The gateway mints `X-Request-ID` and forwards it unchanged. This service does not mint an ID of
its own: it reads the incoming header, echoes it back on the response when it was present, and
writes one log line per request carrying it:

```
request_id=4f2b1c9e8a7d6f5e4c3b2a1908f7e6d5 method=GET path=/health status=200
```

A request that arrives without the header is still served normally; its log line simply carries
an empty `request_id=`.

Run the tests with `go test ./...`.

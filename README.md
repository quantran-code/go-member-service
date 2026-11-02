# go-member-service

Owns library members. SQLite with GORM; the schema is created at start-up.

```sh
DB_PATH=members.db PORT=8081 go run .
```

| Route | Behaviour |
|---|---|
| `GET /health` | liveness |
| `POST /members` | create; `409` when the email is taken |
| `GET /members` | all members, oldest first |
| `GET /members/{id}` | `404` when absent, `400` for a non-numeric id |
| `PUT /members/{id}` | replace name and email |
| `DELETE /members/{id}` | `204`, `404` when already gone |

Run the tests with `go test ./...`.

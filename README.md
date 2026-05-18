
## Backend setup

Go backend is located in the `backend` directory.

### Environment variables

Create local `.env` file from example:

```bash
cd backend
cp .env.example .env
```

### Start PostgreSQL

```bash
cd backend
docker compose up -d
```

### Run migrations

```bash
migrate -path migrations/place -database "postgres://ryden_user:ryden_password@127.0.0.1:5551/ryden_places?sslmode=disable" up
```

Rollback migrations:

```bash
migrate -path migrations/place -database "postgres://ryden_user:ryden_password@127.0.0.1:5551/ryden_places?sslmode=disable" down
```

### Run place-service

```bash
cd backend
go run ./cmd/place-service
```


### Regenerate protobuf files:

```bash
protoc --proto_path=proto --go_out=gen/go --go_opt=paths=source_relative --go-grpc_out=gen/go --go-grpc_opt=paths=source_relative proto/ryden/places/v1/places.proto
```


FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/todo-api .

FROM alpine:3.20
RUN apk add --no-cache ca-certificates && \
    addgroup -S app && adduser -S app -G app && \
    mkdir -p /data && chown -R app:app /data
WORKDIR /app
COPY --from=build /out/todo-api /app/todo-api
USER app
EXPOSE 9090
ENV TODO_DATA_FILE=/data/todos.json
ENTRYPOINT ["/app/todo-api"]

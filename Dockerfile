FROM golang:1.22 AS builder

WORKDIR /app

COPY . .

# Downloading requirement 
RUN go mod download

# creat binary file
RUN CGO_ENABLED=0 GOOS=linux go build -o /api ./cmd/main.go

FROM alpine:latest AS final

# makefile 
RUN apk add --no-cache make git go

WORKDIR /app

# copy binary file
COPY --from=builder /api /api

# copy migrations for make cmd
COPY --from=builder /app/cmd /app/cmd
COPY --from=builder /app/internal/ /app/internal/
COPY --from=builder /app/Makefile /app/Makefile
COPY --from=builder /app/go.mod /app/go.mod
COPY --from=builder /app/go.sum /app/go.sum

EXPOSE 8080

# run app
CMD ["/api"]

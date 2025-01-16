# Start from the official Go image
FROM golang:1.23-alpine

# Install dependencies
RUN apk add --no-cache git bash postgresql-client
RUN go env -w GOFLAGS=-buildvcs=false

WORKDIR /app


COPY . .
#RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
RUN go install github.com/air-verse/air@latest

# Set the Current Working Directory inside the container


# RUN adduser -D -g '' golang
# RUN mkdir -p /app && chown -R golang:golang /app /go/pkg/mod
# WORKDIR /app
# USER golang
# COPY --chown=golang:golang . .

# Run the executable
CMD ["air", "-d", "-c", ".air.toml"]
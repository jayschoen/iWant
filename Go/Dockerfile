# First build: Installs dependencies
FROM golang:alpine as build
RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/server

# Second build: Copy installed dependencies and build final image
FROM scratch
COPY --from=build /app/server /app/

EXPOSE 8000
CMD ["/app/server"]

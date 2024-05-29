FROM golang:1.22.3-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

EXPOSE 5000

RUN CGO_ENABLED=0 GOOS=linux go build -o /golearn-app ./cmd/http/.

CMD ["/golearn-app"]

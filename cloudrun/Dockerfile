FROM golang:1.22-alpine3.19

WORKDIR /app

COPY ./ ./

RUN go build -o main ./svc/cmd/dev/main.go

CMD ["./main"]

ENV HOST=0.0.0.0

ENV PORT=8080

EXPOSE 8080
FROM golang:1.22-alpine3.19

WORKDIR /app

COPY ./main ./main
COPY ./.env ./.env
COPY ./setting-release.yaml ./setting-release.yaml

EXPOSE 8080
CMD ["./main"]

ENV HOST=0.0.0.0

ENV PORT=8080


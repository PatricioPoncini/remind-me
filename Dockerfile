FROM golang:1.22.3 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o remind_me_app .

FROM alpine:latest

RUN apk add --no-cache libc6-compat

COPY --from=build /app/remind_me_app /app/remind_me_app

ENTRYPOINT ["/app/remind_me_app"]

EXPOSE 8080

CMD ["./main"]
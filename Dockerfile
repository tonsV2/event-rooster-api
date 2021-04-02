FROM golang:1.16-alpine AS build
RUN apk add gcc musl-dev
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/event-rooster-api .

RUN apk --no-cache -U upgrade
FROM alpine:3.9
WORKDIR /app
COPY --from=build /app/event-rooster-api .
COPY ./mail/templates ./mail/templates
USER guest
CMD ["/app/event-rooster-api"]

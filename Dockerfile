FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get ./...
RUN go build -ldflags "-X main.version=production`date -u +.%Y%m%d.%H%M%S`" -o shopping-list-bot 
RUN go build -o migrator cmd/db/main.go 

FROM alpine
COPY --from=builder /build/migrator /app/
COPY --from=builder /build/shopping-list-bot /app/
COPY --from=builder /build/sql/ /app/sql/
WORKDIR /app

EXPOSE 443

ENTRYPOINT ["./migrator -action migration && ./shopping-list-bot"]

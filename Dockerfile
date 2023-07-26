FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' shopping-list-bot .

FROM alpine
COPY --from=builder /build/shopping-list-bot /app/
COPY --from=builder /build/sql/ /app/sql/
WORKDIR /app
CMD ["./shopping-list-bot"]

FROM golang:alpine as builder
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -o ./shopping-list-bot 

FROM alpine
COPY --from=builder /build/shopping-list-bot /app/
COPY --from=builder /build/sql/ /app/sql/
WORKDIR /app

# EXPOSE 8080

USER nonroot:nonroot
ENTRYPOINT ["./shopping-list-bot"]

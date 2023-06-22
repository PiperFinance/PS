ARG GO_VERSION=1.20

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
ENV PORT=8000
COPY  ./src/go.mod .
COPY ./src/go.sum .
RUN go mod download

COPY ./src . 
RUN go build -o ./app ./main.go

FROM alpine:latest

RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

RUN mkdir -p /api
WORKDIR /api
COPY --from=builder /api/app .
ADD https://raw.githubusercontent.com/PiperFinance/CD/main/chains/mainnetV2.json /data/mainnets:q
.json  
RUN rm -rf /var/bs/log/ | true \ 
    && mkdir -p /var/bs/log/ \ 
    && touch /var/bs/log/err.log \ 
    && touch /var/bs/log/debug.log 

EXPOSE 8080

ENTRYPOINT ["./app"]

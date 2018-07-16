
# 1) build environment
FROM golang:1.10.3-alpine3.8 as build-env
# Install SSL ca certificates
RUN apk update && apk add git && apk add -U --no-cache ca-certificates && apk add curl
RUN curl -Ss -X POST https://login.microsoftonline.com/
RUN curl -Ss -X GET https://sbrskeyvault.vault.azure.net/
# copy src files
ADD . /src
# get dependencies
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN go get "github.com/spf13/viper"
# go to folder and build golang app as static binary
RUN cd /src/cmd/keyvault && CGO_ENABLED=0 GOOS=linux go build -a -tags netgo -ldflags '-w' -o keyvaultapp
# use port 8080 for the service
EXPOSE 8080

# 2) ssl certificates environment
FROM alpine:latest as alpine
RUN apk --no-cache add ca-certificates

# 3) deployment environment
FROM scratch 
# change folder reference
WORKDIR /app
# copy static executable
COPY --from=build-env /src/cmd/keyvault/keyvaultapp .
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/app/keyvaultapp"]
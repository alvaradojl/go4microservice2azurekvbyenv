FROM golang:1.9.2 as build-env

ADD . /src

COPY . .
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN go get "github.com/spf13/viper"

RUN cd /src/cmd/keyvault && CGO_ENABLED=0 GOOS=linux go build -o keyvaultapp

EXPOSE 8080

FROM scratch
WORKDIR /app
COPY --from=build-env /src/cmd/keyvault/keyvaultapp .
ENTRYPOINT ["/app/keyvaultapp"]
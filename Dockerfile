FROM golang:1.9.2
WORKDIR $GOPATH/src/github.com/alvaradojl/go4microservice2azurekvbyenv/cmd/keyvault
COPY . .
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN go get "github.com/spf13/viper"
RUN CGO_ENABLED=0 GOOS=linux go build .

EXPOSE 8080

FROM scratch
COPY --from=0 . .
ENTRYPOINT ["/keyvault"]

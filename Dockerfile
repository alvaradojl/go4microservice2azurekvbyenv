FROM golang:1.9.2
WORKDIR $HOME/gopath/src/github.com/alvaradojl/go4microservice2azurekvbyenv
ADD . $GOPATH/src/github.com/alvaradojl/go4microservice2azurekvbyenv
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN go get "github.com/spf13/viper"
RUN CGO_ENABLED=0 GOOS=linux go build $GOPATH/src/github.com/alvaradojl/go4microservice2azurekvbyenv/cmd/keyvault

EXPOSE 8080

FROM scratch
COPY --from=0 $GOPATH/src/github.com/alvaradojl/go4microservice2azurekvbyenv/cmd/keyvault .
ENTRYPOINT ["/keyvault"]

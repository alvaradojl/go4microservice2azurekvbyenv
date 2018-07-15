FROM golang:1.9.2
WORKDIR /go/src/github.com/alvaradojl/go4microservice2azurekvbyenv/cmd/keyvault
COPY . .
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN CGO_ENABLED=0 GOOS=linux go build .

EXPOSE 8080

FROM scratch
COPY --from=0 /go/src/github.com/alvaradojl/go4microservice2azurekvbyenv/cmd/keyvault .
ENTRYPOINT ["/keyvault"]

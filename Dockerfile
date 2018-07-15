FROM golang:1.9.2
WORKDIR ./cmd/keyvault
COPY . .
RUN go get "github.com/opentracing-contrib/go-stdLib/nethttp"
RUN go get "github.com/opentracing/opentracing-go"   
RUN go get "github.com/opentracing/opentracing-go/log"
RUN go get "github.com/uber/jaeger-client-go/config"                  
RUN go get "github.com/uber/jaeger-lib/metrics"
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/keyvault

EXPOSE 8080

FROM scratch
COPY --from=0 ./cmd/keyvault .
ENTRYPOINT ["/keyvault"]

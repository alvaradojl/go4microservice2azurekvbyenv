package main

import (
	"context"
	"log"

	"github.com/opentracing/opentracing-go"
	otlog "github.com/opentracing/opentracing-go/log"
)

func infoCtx(ctx context.Context, msg string) {

	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.LogFields(otlog.String("info", msg))
	}

	log.Println(msg)
}

func errorCtx(ctx context.Context, err error) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.LogFields(otlog.Error(err))
	}

	log.Println(err)
}

func tag(ctx context.Context, key string, value string) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		span.SetTag(key, value)
	}
}

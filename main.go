package main

import (
	"fmt"
	"go-pizza-api/internal/deals"
	"log"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

func main() {
	cfg := &config.Configuration{
		ServiceName: "go-pizza",
		// not sure what this is.
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		// log the spans to output.
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: "localhost:6831",
		},
	}

	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		log.Fatalf("err: cannot init jaeger: %v\n", err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	deals := deals.GetDeals(tracer, "CT27WW")
	fmt.Print(deals)
}

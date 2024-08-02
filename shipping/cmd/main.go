package main

import (
	"github.com/phthaocse/microservices-in-go/shipping/config"
	"github.com/phthaocse/microservices-in-go/shipping/internal/adapters/grpc"
	"github.com/phthaocse/microservices-in-go/shipping/internal/application/core/api"
)

const (
	service     = ""
	environment = "dev"
	id          = 2
)

//	func init() {
//		log.SetFormatter(customLogger{
//			formatter: log.JSONFormatter{FieldMap: log.FieldMap{
//				"msg": "message",
//			}},
//		})
//		log.SetOutput(os.Stdout)
//		log.SetLevel(log.InfoLevel)
//	}
//
//	type customLogger struct {
//		formatter log.JSONFormatter
//	}
//
//	func (l customLogger) Format(entry *log.Entry) ([]byte, error) {
//		span := trace.SpanFromContext(entry.Context)
//		entry.Data["trace_id"] = span.SpanContext().TraceID().String()
//		entry.Data["span_id"] = span.SpanContext().SpanID().String()
//		//Below injection is Just to understand what Context has
//		entry.Data["Context"] = span.SpanContext()
//		return l.formatter.Format(entry)
//	}
func main() {
	//tp, err := tracerProvider("http://jaeger-otel.jaeger.svc.cluster.local:14278/api/traces")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//otel.SetTracerProvider(tp)
	//otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))

	//dbAdapter, err := db.NewAdapter(config.GetDataSourceURL())
	//if err != nil {
	//	log.Fatalf("Failed to connect to database. Error: %v", err)
	//}

	application := api.NewApplication()
	grpcAdapter := grpc.NewAdapter(application, config.GetApplicationPort())
	grpcAdapter.Run()
}

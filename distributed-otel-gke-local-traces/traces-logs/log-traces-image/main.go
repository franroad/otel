//this image is instrumented to correlate logs and traces by injecting the trace id into the log

package main

import (
    "context"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/gorilla/mux"
    "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
    exporter, err := otlptracehttp.New(context.Background())
    if err != nil {
        log.Fatalf("resource.New: %v", err)
    }
    tp := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(resource.Default()),
    )
    otel.SetTracerProvider(tp)
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
    tracer := otel.Tracer("example-server")
    ctx, span := tracer.Start(r.Context(), "mainHandler")
    defer span.End()

    traceID := span.SpanContext().TraceID().String()
    spanID := span.SpanContext().SpanID().String()
    log.Printf("Handling request. TraceID: %s, SpanID: %s\n", traceID, spanID)

    destination := os.Getenv("DESTINATION_URL")
    client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
    req, _ := http.NewRequestWithContext(ctx, "GET", destination, nil)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("could not fetch remote endpoint: %v", err)
    }
    defer resp.Body.Close()

    _, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("could not read response from %v", destination)
    }
    fmt.Fprint(w, strconv.Itoa(resp.StatusCode))
}

func main() {
    initTracer()
    otel.SetTextMapPropagator(propagation.TraceContext{})

    r := mux.NewRouter()
    r.HandleFunc("/", mainHandler)

    var handler http.Handler = r
    handler = otelhttp.NewHandler(handler, os.Getenv("HANDLER_NAME"))

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), handler))
}

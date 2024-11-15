/*
Copyright 2019 Google LLC
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
https://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp" // Updated exporter
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    "go.opentelemetry.io/otel/sdk/trace"
)

func initTracer() {
    // Set up an OTLP exporter (compatible with Beyla)
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

    destination := os.Getenv("DESTINATION_URL")
    client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
    req, _ := http.NewRequestWithContext(ctx, "GET", destination, nil)
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal("could not fetch remote endpoint:", err)
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

    // Set the global propagators
    otel.SetTextMapPropagator(propagation.TraceContext{})

    r := mux.NewRouter()
    r.HandleFunc("/", mainHandler)
    var handler http.Handler = r

    handler = otelhttp.NewHandler(handler, "server")
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", os.Getenv("PORT")), handler))
}

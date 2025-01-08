# OTEL OPERATOR

**Purpose:**
- The intention of this repo is to configure a daemon set using the otel [collector operator](https://github.com/open-telemetry/opentelemetry-operator) 

**Baseline:**
- An instrumented app in go  will be used as example.
- The app and the operator infrastructure will be deployed in a k8's cluster in _GKE_
- The app will be accesed via an external lb-service, the request will generate the traces.

**Notes**

- The Operator will create a *ClusterIp* in the namespace where otel object is created, usually named *otel-collector* we need to set an env variable in our app fowarding the traces to the given service like this:

 ``` yaml
 - name: OTEL_EXPORTER_OTLP_TRACES_ENDPOINT
          value: http://otel-collector.monitoring.svc.cluster.local:4318/v1/traces
 ```
**Insigths:**

- The otel-confgi will be configured to forward the traces to a different project. This approach is useful for centralizing the traces in one place.

## DIAGRAM


## Procedure:

[Follow guide](https://github.com/GoogleCloudPlatform/opentelemetry-operator-sample/tree/main)

[cloud trace example](https://github.com/GoogleCloudPlatform/opentelemetry-operator-sample/tree/main/recipes/cloud-trace)


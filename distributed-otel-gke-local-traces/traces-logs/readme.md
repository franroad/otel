# Centralized

This configuration is intended for **centralizing** the Traces in **Cloud Trace** additionally, the app is instrumented to add the trace ID in the logs thus enabling traceability between the traces and the logs.

1.- The setup consists in two clusters in GCP, the frontend cluster calls the backend cluster, hosted in different projecst. This logic is done at the app level.

2.- Then, thanks to the instrumentation the traces are fowarded to the otel collector through a cluster ip service and then the **collector** forward the infor to the specified gcp projects in its's configuration using the SA specified (creating a secret).

3. - In this example only *one* service account should be used, should be created in the project that is centralizing the traces. 

## We achieve the following:

* Correlate logs with traces
* As the traces has it's span id the full request can be seen in Cloud Trace
* Otel collector configuration for forwarding traces to  **the same** GCP so teh full trace can be visualized at destination.
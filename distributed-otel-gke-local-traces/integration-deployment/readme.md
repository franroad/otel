# Centralized

This configuration is intended for each cluster located in different projects to forward the Traces to their own **Cloud Trace** but using the OTEL Collector

1.- The setup consists in two clusters in GCP, the frontend cluster calls the backend cluster, hosted in different projecst. This logic is done at the app level.

2.- Then, thanks to the instrumentation the traces are fowarded to the otel collector throgught a cluster ip service and then the **collector** forward the infor to the specified gcp projects in its's configuration (otel-config-logs.yaml) using the SA specified (creating a secret).

3. - In this scenario *two* services accoounts should be created one fo each project that is recieving the traces.

## We achieve the following:

* Correlate logs with traces
* This setup should be replicated in both gke  to observe the traces and OTEL collector behavior and capabilities.
* Otel collector configuration for forwarding traces to **different** GCP projects.
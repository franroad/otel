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

# TO-DO: 
- Describe the procedure and add the wkld identity part instead of the sa CLEAN THE PROCEDURE REFERENCE

- ADD AND TEST THE AUTO-INSTUMENTATION WITH PYTHON.

## DIAGRAMs


## Procedure:

[Follow guide](https://github.com/GoogleCloudPlatform/opentelemetry-operator-sample/tree/main)

[cloud trace example](https://github.com/GoogleCloudPlatform/opentelemetry-operator-sample/tree/main/recipes/cloud-trace)



# PROCEDURE REFERENCE:

-INSTALANDO EN UN STANDARD CLUSTER EL OPERATOR:

1 - install certs manager

2 - install the operator : 
```sh
kubectl apply -f https://github.com/open-telemetry/opentelemetry-operator/releases/download/v0.112.0/opentelemetry-operator.yaml 
```
3 - Deploy the app 
4 - create the SA’s the GCP IAM and the  k8’s one the k’8s sa must be created in the ns that will host the collector  deployment
5- Enable workload identity 
	HABILITAR  EL WORKLOAD IDENTIY A NIVEL DE CLUSTER Y A NIVEL DE NODOS 
    WARNING:: Modifying the node pool immediately enables Workload Identity Federation for GKE for any workloads running in the node pool. This prevents the workloads from using the service account that your nodes use and might result in disruptions. 

6-  Create the binding between the 2 SA 

```sh
gcloud iam service-accounts add-iam-policy-binding \
  tracesbackend-pre@afb-backend-pre.iam.gserviceaccount.com \ # <gcp-service-account> 
  --role=roles/iam.workloadIdentityUser \
  --member="serviceAccount:afb-backend-pre.svc.id.goog[monitoring/otel-collector]" # "serviceAccount:<project>.svc.id.goog[<namespace>/<k8s-service-account>]"
  ```


7-  Create the deployment  in the monitoring ns, using the operator
```yaml
apiVersion: opentelemetry.io/v1alpha1
kind: OpenTelemetryCollector
metadata:
  name: otel
  namespace: monitoring 
spec:
  image: otel/opentelemetry-collector-contrib:0.106.0
  config: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
          http:
            endpoint: 0.0.0.0:4318
    processors:

    exporters:
      googlecloud:
        project: afb-backend-pre
      debug:
        verbosity: detailed

    service:
      pipelines:
        traces:
          receivers: [otlp]
          processors: []
          exporters: [debug, googlecloud]
  ```

8- add the annotation to the deployment (Otel object) for allowing the  Collector's ServiceAccount to use Workload Identity

It works because the operator looks for a sa called Otel-collector
```sh
kubectl annotate opentelemetrycollector otel \
    --namespace $COLLECTOR_NAMESPACE \
    iam.gke.io/gcp-service-account=<gcp-sa>
```
Additionally, in case you have deployments that currently uses sa, you can  update the  sa  to use an sa and allow them to use the gcp sa’s:

gcloud iam service-accounts add-iam-policy-binding \
  PROJECT_NUMBER-compute@developer.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:<project>.svc.id.goog[<namespace>/<k8s-service-account>]"

gcloud iam service-accounts add-iam-policy-binding \
  788959310267-compute@developer.gserviceaccount.com \
  --role roles/iam.workloadIdentityUser \
  --member "serviceAccount:afb-backend-pre.svc.id.goog[monitoring/monitoring-sa]”

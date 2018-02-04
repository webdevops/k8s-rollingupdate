# Kubernetes rollingupdate trigger for Deployments, Daemonsets and Statefulsets

This tool triggers a rolling update for Deployments, Daemonsets and Statefulsets with possiblity to filter by annotation.

## Usage

```
Usage:
  k8s-rollingupdate [OPTIONS]

Application Options:
      --kubeconfig=             Path to .kube/config [$KUBECONFIG]
      --kubecontext=            Context of .kube/config [$KUBECONTEXT]
  -n, --namespace=              Namespace to process [$K8S_ROLLINGUPDATE_NAMESPACE]
      --annotation=             Filter Kubernetes object by annotation
                                [$K8S_ROLLINGUPDATE_ANNOTATION]
      --annotation-value=       Filter Kubernetes object by annotation value (needs --annotation)
                                [$K8S_ROLLINGUPDATE_ANNOTATION_VALUE]
      --annotation-autorollout= Annotation which will be added to trigger rolling update (default:
                                rolllingupdate.webdevops.io/trigger)
                                [$K8S_ROLLINGUPDATE_ANNOTATION_TRIGGER]

Help Options:
  -h, --help                    Show this help message
```

```bash
# Trigger rolling update for everything in namespace
k8s-rollingupdate -n your-namespace

# Trigger rolling update for everything multiple namespaces
k8s-rollingupdate -n your-namespace -n other-namespace

# Trigger rolling update for everything with annotation foobar/barfoo
k8s-rollingupdate -n your-namespace --annotation foobar/barfoo

# Trigger rolling update for everything with annotation foobar/barfoo=value
k8s-rollingupdate -n your-namespace --annotation foobar/barfoo --annotation-value value

```

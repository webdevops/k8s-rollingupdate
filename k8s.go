package main

import (
	"time"
	"fmt"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v13 "k8s.io/api/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type Kubernetes struct {
	clientset *kubernetes.Clientset

	Logger *DaemonLogger

	KubeContext string
	KubeConfig string
	AnnotationTrigger string
	AnnotationSelector string
	AnnotationSelectorValue string
}

// Create cached kubernetes client
func (k *Kubernetes) Client() (clientset *kubernetes.Clientset) {
	var err error
	var config *rest.Config

	if k.clientset == nil {
		if k.KubeConfig != "" {
			// KUBECONFIG
			config, err = buildConfigFromFlags(k.KubeContext, k.KubeConfig)
			if err != nil {
				panic(err.Error())
			}
		} else {
			// K8S in cluster
			config, err = rest.InClusterConfig()
			if err != nil {
				panic(err.Error())
			}
		}

		k.clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			panic(err.Error())
		}
	}

	return k.clientset
}

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

func (k *Kubernetes) TriggerRollout(namespace string) (error error) {
	Logger.Println(fmt.Sprintf("Processing namespace %v", namespace))

	if err := k.TriggerRolloutDeployment(namespace); err != nil {
		return err
	}

	if err := k.TriggerRolloutDaemonset(namespace); err != nil {
		return err
	}

	if err := k.TriggerRolloutStatefulset(namespace); err != nil {
		return err
	}

	return
}

func (k *Kubernetes) TriggerRolloutDeployment(namespace string) (error error) {
	list, err := k.ListDeployments(namespace)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		Logger.Println(fmt.Sprintf(" - triggering rolling update for %v/%v", "deployment", item.GetName()))
		item.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		item.Spec.Template.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		if _, err := k.Client().AppsV1().Deployments(namespace).Update(&item); err != nil {
			return err
		}
	}

	return
}

func (k *Kubernetes) TriggerRolloutDaemonset(namespace string) (error error) {
	list, err := k.ListDaemonsets(namespace)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		Logger.Println(fmt.Sprintf(" - triggering rolling update for %v/%v", "daemonset", item.GetName()))
		item.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		item.Spec.Template.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		if _, err := k.Client().AppsV1().DaemonSets(namespace).Update(&item); err != nil {
			return err
		}
	}

	return
}

func (k *Kubernetes) TriggerRolloutStatefulset(namespace string) (error error) {
	list, err := k.ListStatefulsets(namespace)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		Logger.Println(fmt.Sprintf(" - triggering rolling update for %v/%v", "statefulset", item.GetName()))
		item.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		item.Spec.Template.Annotations[k.AnnotationTrigger] =  time.Now().Format(time.RFC3339)
		if _, err := k.Client().AppsV1().StatefulSets(namespace).Update(&item); err != nil {
			return err
		}
	}

	return
}

func (k *Kubernetes) ListDaemonsets(namespace string) (list v13.DaemonSetList, error error) {
	option := v12.ListOptions{}

	if valList, err := k.Client().AppsV1().DaemonSets(namespace).List(option); err == nil {
		// return all if no selector available
		if k.AnnotationSelector == "" {
			list = *valList
			return
		}

		for _, item := range valList.Items {
			if item.Annotations == nil {
				continue
			}

			if val, ok := item.Annotations[k.AnnotationSelector]; ok {
				if k.AnnotationSelectorValue == "" || val == k.AnnotationSelectorValue {
					list.Items = append(list.Items, item)
				}
			}
		}
	} else {
		error = err
	}

	return
}


func (k *Kubernetes) ListDeployments(namespace string) (list v13.DeploymentList, error error) {
	option := v12.ListOptions{}

	if valList, err := k.Client().AppsV1().Deployments(namespace).List(option); err == nil {
		// return all if no selector available
		if k.AnnotationSelector == "" {
			list = *valList
			return
		}

		for _, item := range valList.Items {
			if item.Annotations == nil {
				continue
			}

			if val, ok := item.Annotations[k.AnnotationSelector]; ok {
				if k.AnnotationSelectorValue == "" || val == k.AnnotationSelectorValue {
					list.Items = append(list.Items, item)
				}
			}
		}
	} else {
		error = err
	}

	return
}

func (k *Kubernetes) ListStatefulsets(namespace string) (list v13.StatefulSetList, error error) {
	option := v12.ListOptions{}

	if valList, err := k.Client().AppsV1().StatefulSets(namespace).List(option); err == nil {
		// return all if no selector available
		if k.AnnotationSelector == "" {
			list = *valList
			return
		}

		for _, item := range valList.Items {
			if item.Annotations == nil {
				continue
			}

			if val, ok := item.Annotations[k.AnnotationSelector]; ok {
				if k.AnnotationSelectorValue == "" || val == k.AnnotationSelectorValue {
					list.Items = append(list.Items, item)
				}
			}
		}
	} else {
		error = err
	}

	return
}

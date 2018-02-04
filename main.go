package main

import (
	"os"
	"fmt"
	flags "github.com/jessevdk/go-flags"
)

var opts struct {
	KubeConfig               string   `           long:"kubeconfig"              env:"KUBECONFIG"                                    description:"Path to .kube/config"`
	KubeContext              string   `           long:"kubecontext"             env:"KUBECONTEXT"                                   description:"Context of .kube/config"`
	Namespace                []string `short:"n"  long:"namespace"               env:"K8S_ROLLINGUPDATE_NAMESPACE" env-delim:" "     description:"Namespace to process"              required:"true"`
	AnnotationSelector       string   `           long:"annotation"              env:"K8S_ROLLINGUPDATE_ANNOTATION"                  description:"Filter Kubernetes object by annotation"`
	AnnotationSelectorValue  string   `           long:"annotation-value"        env:"K8S_ROLLINGUPDATE_ANNOTATION_VALUE"            description:"Filter Kubernetes object by annotation value (needs --annotation)"`
	AnnotationTrigger        string   `           long:"annotation-autorollout"  env:"K8S_ROLLINGUPDATE_ANNOTATION_TRIGGER"          description:"Annotation which will be added to trigger rolling update"    default:"rolllingupdate.webdevops.io/trigger"`
}

var (
	argparser *flags.Parser
	args []string
	k8sService = Kubernetes{}
	Logger *DaemonLogger
	ErrorLogger *DaemonLogger
)

func main() {
	var err error
	argparser = flags.NewParser(&opts, flags.Default)
	args, err = argparser.Parse()

	initOpts()

	// Init logger
	Logger = CreateDaemonLogger(0)
	ErrorLogger = CreateDaemonErrorLogger(0)

	// check if there is an parse error
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			fmt.Println()
			argparser.WriteHelp(os.Stdout)
			os.Exit(1)
		}
	}

	k8sService.KubeConfig = opts.KubeConfig
	k8sService.KubeContext = opts.KubeContext
	k8sService.Logger = Logger
	k8sService.AnnotationTrigger = opts.AnnotationTrigger
	k8sService.AnnotationSelector = opts.AnnotationSelector
	k8sService.AnnotationSelectorValue = opts.AnnotationSelectorValue

	for _, namespace := range opts.Namespace {
		if err := k8sService.TriggerRollout(namespace); err != nil {
			ErrorLogger.Error("failed rollout", err)
		}
	}
}

func initOpts() {
	if opts.KubeConfig == "" {
		kubeconfigPath := fmt.Sprintf("%s/.kube/config", UserHomeDir())
		if _, err := os.Stat(kubeconfigPath); err == nil {
			opts.KubeConfig = kubeconfigPath
		}
	}
}

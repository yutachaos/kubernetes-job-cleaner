package main

import (
	"context"
	"flag"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/klog"
	"path/filepath"
	"time"
)

func main() {
	var kubeconfig *string

	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	klog.Info("Start kubernetes job cleaner")

	for {
		jobs, err := clientset.BatchV1().Jobs("").List(context.TODO(), v1.ListOptions{
			FieldSelector: "status.successful=1",
		})

		if err != nil {
			panic(err.Error())
		}

		klog.Infof("There are %d succeeded jobs in the cluster", len(jobs.Items))

		deletePolicy := v1.DeletePropagationForeground

		for _, job := range jobs.Items {
			if err := clientset.BatchV1().Jobs(job.Namespace).Delete(context.TODO(), job.Name, v1.DeleteOptions{
				PropagationPolicy: &deletePolicy,
			}); err != nil {
				panic(err)
			}
			klog.Infof("Delete job Namespame:%s JobName: %s", job.Namespace, job.Name)
		}

		time.Sleep(10 * time.Second)
	}
}

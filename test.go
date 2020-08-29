package main

import (
	"bytes"
	"github.com/parvez0/disk-watcher/watchers"
	"github.com/sirupsen/logrus"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func CreateClient() (*kubernetes.Clientset, *rest.Config) {
	var kubeconfig string
	var config *rest.Config
	var err error
	if os.Getenv("IN_CLUSTER") != "true" {
		if home := homedir.HomeDir(); home != "" {
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err)
		}
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset, config
}

func GetPvcObject(namespace *string) v1.PersistentVolumeClaimInterface {
	client, _ := CreateClient()
	return client.CoreV1().PersistentVolumeClaims(*namespace)
}

func IncreaseWhatsappDiskSize(namespace string) int64 {
	pvcClient := GetPvcObject(&namespace)
	pvc, err := pvcClient.Get("whatsapp-disk", metav1.GetOptions{})
	if err != nil {
		logrus.Panic("Failed to get pvc whatsapp-disk - ", namespace, " ", err)
	}
	quantity := pvc.Spec.Resources.Requests["storage"]
	value, _ := quantity.AsInt64()
	newValue := value + (100 * 1024 * 1024 * 1024)
	quantity.Set(newValue)
	logrus.Info("updating pvc increasing the size to - ", newValue/(1024*1024*1024), "GB")
	pvc.Spec.Resources.Requests["storage"] = quantity
	updatedPvc, err := pvcClient.Update(pvc)
	if err != nil {
		logrus.Panic("failed to update the pvc size - ", err)
	}
	logrus.Printf("pvc updated successfully - %+v", updatedPvc.Spec)
	return newValue / (1024 * 1024 * 1024)
}

func main() {
	logger := watchers.NewLogger()
	client, config := CreateClient()
	execClient := client.CoreV1().RESTClient().Post().Resource("pods").Name("master-0").Namespace("wa-udaan").SubResource("exec")
	exOpts := &v12.PodExecOptions{
		Command: []string{"df", "-h"},
		Stdout: true,
	}
	execClient.VersionedParams(exOpts, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(config, "POST", execClient.URL())
	if err != nil{
		logger.Panic("failed to create an executor - ", err)
	}
	var stdout bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
	})
	if err != nil{
		logger.Panic("failed to exec into the container - ", err)
	}
	logger.Info(stdout.String())
}

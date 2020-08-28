package watchers

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"os"
	"path/filepath"
)

func CreateClient() *kubernetes.Clientset {
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
	return clientset
}

func GetPvcObject(namespace *string) v1.PersistentVolumeClaimInterface {
	client := CreateClient()
	return client.CoreV1().PersistentVolumeClaims(*namespace)
}

func IncreaseWhatsappDiskSize(namespace string) int64 {
	pvcClient := GetPvcObject(&namespace)
	pvc, err := pvcClient.Get("whatsapp-disk", metav1.GetOptions{})
	if err != nil {
		logger.Panic("Failed to get pvc whatsapp-disk - ", namespace, " ", err)
	}
	quantity := pvc.Spec.Resources.Requests["storage"]
	value, _ := quantity.AsInt64()
	newValue := value + (100 * 1024 * 1024 * 1024)
	quantity.Set(newValue)
	logger.Info("updating pvc increasing the size to - ", newValue/(1024*1024*1024), "GB")
	pvc.Spec.Resources.Requests["storage"] = quantity
	updatedPvc, err := pvcClient.Update(pvc)
	if err != nil {
		logger.Panic("failed to update the pvc size - ", err)
	}
	logger.Printf("pvc updated successfully - %+v", updatedPvc.Spec)
	return newValue / (1024 * 1024 * 1024)
}

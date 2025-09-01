package cmd

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig string
var defaultKubeconfigPath string

func getKubeClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func definedDefualtKubeconfigPath() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	viper.BindEnv("default_kubeconfig")
	defaultKubeconfigPath = viper.GetString("default_kubeconfig")
}

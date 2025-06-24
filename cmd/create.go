package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
)

var manifestPath string

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create Kubernetes deployments in the default namespace",
	Run: func(cmd *cobra.Command, args []string) {

		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		deploymentManifest, err := loadDeploymentFromYAML(manifestPath)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to read YAML manifest in %s", manifestPath)
			os.Exit(1)
		}
		deployment, err := clientset.AppsV1().Deployments("default").Create(context.Background(), deploymentManifest, metav1.CreateOptions{})
		if err != nil {
			log.Error().Err(err).Msg("Failed to create deployments")
			os.Exit(1)
		}
		fmt.Printf("Created deployment %q.\n", deployment.GetObjectMeta().GetName())
	},
}

func loadDeploymentFromYAML(path string) (*appsv1.Deployment, error) {
	var deployment appsv1.Deployment

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	if err := yaml.NewYAMLOrJSONDecoder(file, 1024).Decode(&deployment); err != nil {
		return nil, err
	}
	return &deployment, nil
}

func init() {
	rootCmd.AddCommand(createCmd)

	definedDefualtKubeconfigPath()
	createCmd.Flags().StringVar(&manifestPath, "manifest-path", "./deployment.yaml", "Path to deployment manifest")
	createCmd.Flags().StringVar(&kubeconfig, "kubeconfig", defaultKubeconfigPath, "Path to the kubeconfig file")
}

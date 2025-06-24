package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var deploymentName string

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Kubernetes deployments in the default namespace",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		if clientset.AppsV1().Deployments("default").Delete(context.Background(), deploymentName, metav1.DeleteOptions{}) != nil {
			log.Error().Err(err).Msg("Failed to delete deployment")
			os.Exit(1)
		}
		fmt.Printf("%s deployment deleted.\n", deploymentName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	definedDefualtKubeconfigPath()
	deleteCmd.Flags().StringVar(&deploymentName, "deployment", "", "Deployment name")
	deleteCmd.Flags().StringVar(&kubeconfig, "kubeconfig", defaultKubeconfigPath, "Path to the kubeconfig file")
}

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/dimapakki95/k8s-controllers/pkg/informer"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var serverPort int
var serverKubeconfig string
var serverInCluster bool
var timeInterval int
var namespace string
var resourceTypes string

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: fmt.Sprintf("Start a FastHTTP server and informer for the resources: ", resourceTypes),
	Run: func(cmd *cobra.Command, args []string) {
		level := parseLogLevel(logLevel)
		configureLogger(level)
		clientset, err := getServerKubeClient(serverKubeconfig, serverInCluster)
		if err != nil {
			log.Error().Err(err).Msg("Failed to create Kubernetes client")
			os.Exit(1)
		}
		ctx := context.Background()
		informer.StartInformers(ctx, clientset, timeInterval, namespace, resourceTypes)

		handler := func(ctx *fasthttp.RequestCtx) {
			fmt.Fprintf(ctx, "Hello from FastHTTP!")
		}
		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting FastHTTP server on %s (version: %s)", addr, appVersion)
		if err := fasthttp.ListenAndServe(addr, handler); err != nil {
			log.Error().Err(err).Msg("Error starting FastHTTP server")
			os.Exit(1)
		}
	},
}

func getServerKubeClient(kubeconfigPath string, inCluster bool) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	viper.BindEnv("port")
	viper.BindEnv("namespace")
	viper.BindEnv("kubeconfig_path")
	viper.BindEnv("time_interval")
	viper.BindEnv("resource_types")

	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&resourceTypes, "resource-types", viper.GetString("resource_types"), "Kubernetes resource types for the informer")
	serverCmd.Flags().StringVar(&namespace, "namespace", viper.GetString("namespace"), "Kubernetes namespace")
	serverCmd.Flags().IntVar(&timeInterval, "time-interval", viper.GetInt("time_interval"), "Time in seconds the state will be synced with the cache")
	serverCmd.Flags().IntVar(&serverPort, "port", viper.GetInt("port"), "Port to run the server on")
	serverCmd.Flags().StringVar(&serverKubeconfig, "kubeconfig", viper.GetString("kubeconfig_path"), "Path to the kubeconfig file")
	serverCmd.Flags().BoolVar(&serverInCluster, "in-cluster", false, "Use in-cluster Kubernetes config")
}

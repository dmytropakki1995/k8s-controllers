package informer

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

func StartInformers(ctx context.Context, clientset *kubernetes.Clientset, timeInterval int, namespace string, resourceTypes string) {
	resources := strings.Split(resourceTypes, ",")
	for i := range len(resources) {
		go StartInformer(ctx, clientset, timeInterval, namespace, resources[i])
	}
}

// StartInformer starts a shared informer for k8s rsource in the specified namespace.
func StartInformer(ctx context.Context, clientset *kubernetes.Clientset, timeInterval int, namespace string, resourceType string) {
	factory := informers.NewSharedInformerFactoryWithOptions(
		clientset,
		time.Duration(timeInterval)*time.Second,
		informers.WithNamespace(namespace),
		informers.WithTweakListOptions(func(options *metav1.ListOptions) {
			options.FieldSelector = fields.Everything().String()
		}),
	)

	var informer cache.SharedIndexInformer

	switch resourceType {
	case "deployment":
		informer = factory.Apps().V1().Deployments().Informer()
	case "replicaset":
		informer = factory.Apps().V1().ReplicaSets().Informer()
	case "pod":
		informer = factory.Core().V1().Pods().Informer()
	default:
		fmt.Println("Unknown k8s resource type")
		os.Exit(1)
	}

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Info().Msgf("%s added: %s", resourceType, getObjectName(obj))
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Info().Msgf("%s updated: %s", resourceType, getObjectName(newObj))
		},
		DeleteFunc: func(obj interface{}) {
			log.Info().Msgf("%s deleted: %s", resourceType, getObjectName(obj))
		},
	})
	log.Info().Msgf("Starting %s informer...", resourceType)
	factory.Start(ctx.Done())
	for t, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			log.Error().Msgf("Failed to sync informer for %v", t)
			os.Exit(1)
		}
	}
	log.Info().Msgf("%s informer cache synced. Watching for events for every %d seconds...", resourceType, timeInterval)
	<-ctx.Done() // Block until context is cancelled
}

func getObjectName(obj any) string {
	if object, ok := obj.(metav1.Object); ok {
		return object.GetName()
	}
	return "unknown"
}

package controller

import (
	"context"
	"fmt"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/k8s"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/logging"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"time"
)

const (
	// how often we resync informers resources (besides receiving updates)
	resourcesSyncInterval = time.Hour * 1
)

// Informer defines the required SharedIndexInformers that interact with the API server.
type Informer struct {
	Pod cache.SharedIndexInformer
}

// InformerFactory contains shared informer factory
type InformerFactory struct {
	Pod informers.SharedInformerFactory
}

type Controller struct {
	kubeClient *kubernetes.Clientset

	// informer factories
	factories *InformerFactory

	// informer contains the cache Informers
	informers *Informer

	stopChan chan struct{}
}

func NewController(
	kubeClient *kubernetes.Clientset,
	stopChan chan struct{},
) *Controller {
	controller := &Controller{
		kubeClient: kubeClient,
		stopChan:   stopChan,
		factories:  &InformerFactory{},
		informers:  &Informer{},
	}

	// Create informer factories
	controller.factories.Pod = informers.NewSharedInformerFactoryWithOptions(
		kubeClient,
		resourcesSyncInterval,
	)

	podParams := k8s.PodParams{
		InformerFactory: controller.factories.Pod,
		KubeClient:      kubeClient,
		Ctx:             context.Background(),
	}
	controller.informers.Pod = k8s.WatchPod(podParams)

	return controller
}

func (c *Controller) Run() {
	defer runtime.HandleCrash()

	// start informers where we listen to new resources
	go c.informers.Pod.Run(c.stopChan)

	// wait for all involved caches to be synced before processing items
	// from the queue
	if !cache.WaitForCacheSync(c.stopChan, c.informers.Pod.HasSynced) {
		runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
	}

	// wait for SIGTERM
	<-c.stopChan
	logging.GetLogger().Info().Msgf("stopping ingress controller")
}

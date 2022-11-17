package controller

import (
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/client"
	"testing"
)

func TestController_Run(t *testing.T) {
	client.K8sInit()
	stopCh := make(chan struct{}, 1)
	defer close(stopCh)

	kubeClient := client.K8sClient()
	c := NewController(kubeClient, stopCh)
	c.Run()
}

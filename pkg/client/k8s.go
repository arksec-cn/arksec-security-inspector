package client

import (
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

var k8sClient *kubernetes.Clientset

func K8sInit() {
	if k8sClient != nil {
		logging.GetLogger().Debug().Msgf("k8s client already init")
		return
	}
	var err error
	var config *rest.Config

	config, err = rest.InClusterConfig()
	if err != nil {
		home := os.Getenv("HOME")
		configFile := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", configFile)
		if err != nil {
			logging.GetLogger().Error().Msgf("get k8s config err: %v", err)
			return
		}
	}

	k8sClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		logging.GetLogger().Error().Msgf("new k8s client err: %v", err)
		return
	}
}

func K8sClient() *kubernetes.Clientset {
	if k8sClient == nil {
		panic("k8sClient 模块没有初始化")
	}

	return k8sClient
}

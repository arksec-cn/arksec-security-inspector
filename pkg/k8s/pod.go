package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/client"
	"gitlab.arksec.cn/vegeta/security-inspector/pkg/logging"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

type PodParams struct {
	InformerFactory informers.SharedInformerFactory
	KubeClient      *kubernetes.Clientset
	Ctx             context.Context
}

type PodInfo struct {
	Name           string
	Namespace      string
	ControllerName string
	ControllerType string
}

func WatchPod(options PodParams) cache.SharedIndexInformer {
	informer := options.InformerFactory.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod, ok := obj.(*corev1.Pod)
			if !ok {
				return
			}
			logging.GetLogger().Info().Msgf("[add] pod: %s, namespace: %s", pod.Name, pod.Namespace)
			var (
				controllerName string
				controllerType string
				err            error
			)
			if len(pod.OwnerReferences) > 0 {
				controllerType = pod.OwnerReferences[0].Kind
				switch pod.OwnerReferences[0].Kind {
				case "ReplicaSet":
					get, err := options.KubeClient.AppsV1().ReplicaSets(pod.GetNamespace()).Get(options.Ctx, pod.OwnerReferences[0].Name, metav1.GetOptions{})
					if err == nil {
						if len(get.OwnerReferences) > 0 {
							if get.OwnerReferences[0].Kind == "Deployment" {
								controllerName = get.OwnerReferences[0].Name
								controllerType = "Deployment"
							}
						}
					}
				case "Deployment":
					controllerName = pod.OwnerReferences[0].Name
				case "StatefulSet":
					controllerName = pod.OwnerReferences[0].Name
				case "DaemonSet":
					controllerName = pod.OwnerReferences[0].Name
				case "Job":
					jobName := pod.OwnerReferences[0].Name
					get, err := options.KubeClient.BatchV1beta1().CronJobs(pod.GetNamespace()).Get(options.Ctx, jobName, metav1.GetOptions{})
					if err == nil {
						if len(get.OwnerReferences) > 0 {
							if get.OwnerReferences[0].Kind == "CronJob" {
								controllerName = get.OwnerReferences[0].Name
								controllerType = "CronJob"
							}
						}
					}
				}
				if controllerName != "" && controllerType != "" {
					logging.GetLogger().Info().Msgf("controller type: %s, name: %s", controllerType, controllerName)
				}
			}

			// TODO 调用函数
			info := PodInfo{
				Name:           pod.GetName(),
				Namespace:      pod.GetNamespace(),
				ControllerName: controllerName,
				ControllerType: controllerType,
			}
			data, _ := json.Marshal(info)
			id := fmt.Sprintf("%s-%s", pod.Namespace, pod.Name)
			err = client.OsClient.Add(id, data)
			if err != nil {
				logging.GetLogger().Error().Msgf("send %s open search err: %v", id, err)
			}

			return
		},
		UpdateFunc: func(oldObj, newObj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {

		},
	})

	return informer
}

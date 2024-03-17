package scheduler

import (
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/cache"
	"k8s.io/client-go/kubernetes"
)

func NewGPUsharePrioritize(clientset *kubernetes.Clientset, c *cache.SchedulerCache) *Prioritize {
	return &Prioritize{Name: "gpusharingprioritize", cache: c}
}

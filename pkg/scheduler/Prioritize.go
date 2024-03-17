package scheduler

import (
	"fmt"

	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/cache"
	"github.com/AliyunContainerService/gpushare-scheduler-extender/pkg/log"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
)

type Prioritize struct {
	Name  string
	cache *cache.SchedulerCache
}

func (p Prioritize) Handler(args *schedulerapi.ExtenderArgs) (*schedulerapi.HostPriorityList, error) {
	var priorityList schedulerapi.HostPriorityList
	var nodeNames []string
	if args.NodeNames != nil {
		nodeNames = *args.NodeNames
		log.V(3).Info("extender args NodeNames is not nil, result %+v", nodeNames)
	} else if args.Nodes != nil {
		for _, n := range args.Nodes.Items {
			nodeNames = append(nodeNames, n.Name)
		}
		log.V(3).Info("extender args Nodes is not nil, names is %+v", nodeNames)
	} else {
		return &priorityList, fmt.Errorf("cannot get node names")
	}
	for i, nodename := range nodeNames {
		priorityList[i] = schedulerapi.HostPriority{
			Host:  nodename,
			Score: 0,
		}
	}
	// 计算每一个Node上剩余的显存大小，作为分数返回
	for _, nodeName := range nodeNames {
		nodeInfo, err := p.cache.GetNodeInfo(nodeName)
		if err != nil {
			log.V(3).Info("error: failed to get node info for %s, error: %v", nodeName, err)
			continue
		}
		node := nodeInfo.GetNode()
		if node == nil {
			log.V(3).Info("error: failed to get node for %s", nodeName)
			continue
		}
		// 计算每一个Node上剩余的显存大小，作为分数返回
		// 1. 获取Node上的所有GPU设备
		devices := nodeInfo.GetDevs()
		// 2. 计算每一个GPU设备上的剩余显存大小
		for _, dev := range devices {
			// 3. 计算每一个GPU设备上的剩余显存大小
			// 4. 将每一个GPU设备上的剩余显存大小作为分数返回
			totalMemory := int64(dev.GetTotalGPUMemory())
			usedMemory := int64(dev.GetUsedGPUMemory())
			score := totalMemory - usedMemory
			priorityList = append(priorityList, schedulerapi.HostPriority{
				Host:  nodeName,
				Score: score,
			})
		}
		log.V(3).Info("info: The node %s has been added to the priority list", nodeName)
		log.V(100).Info("predicate result for %s, is %+v", nodeNames, priorityList)
	}
	return &priorityList, nil
}

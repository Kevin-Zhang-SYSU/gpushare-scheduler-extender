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
	var nodeNames []string
	if args.NodeNames != nil {
		nodeNames = *args.NodeNames
		log.V(3).Info("info: extender args NodeNames is not nil, result %+v", nodeNames)
	} else if args.Nodes != nil {
		for _, n := range args.Nodes.Items {
			nodeNames = append(nodeNames, n.Name)
		}
		log.V(3).Info("info: extender args Nodes is not nil, names is %+v", nodeNames)
	} else {
		log.V(3).Info("error: cannot get node names")
		return nil, fmt.Errorf("cannot get node names")
	}
	if len(nodeNames) == 0 {
		log.V(3).Info("error: node names slice is empty")
		return nil, fmt.Errorf("node names slice is empty")
	}
	priorityList := make(schedulerapi.HostPriorityList, 0, len(nodeNames))
	log.V(3).Info("info: The node names are %v", nodeNames)
	for _, nodename := range nodeNames {
		hostpriority := schedulerapi.HostPriority{
			Host:  nodename,
			Score: 0,
		}
		priorityList = append(priorityList, hostpriority)
		log.V(3).Info("info: The node %s has been added to the priority list", nodename)
	}
	if len(nodeNames) == 1 {
		log.V(3).Info("info: only one node, return directly")
		return &priorityList, nil
	}
	log.V(3).Info("info: calculate the priority for each node")
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
		score := int64(0)
		// 2. 计算每一个GPU设备上的剩余显存大小
		for _, dev := range devices {
			// 3. 计算每一个GPU设备上的剩余显存大小
			// 4. 将每一个GPU设备上的剩余显存大小作为分数返回
			totalMemory := int64(dev.GetTotalGPUMemory())
			usedMemory := int64(dev.GetUsedGPUMemory())
			score += (totalMemory - usedMemory)
		}
		// 找到priorityList中对应的node，更新score
		for i, hostPriority := range priorityList {
			if hostPriority.Host == nodeName {
				priorityList[i].Score = 100 - score
				break
			}
		}
		log.V(3).Info("info: The node %s has score %d", nodeName, score)
	}
	log.V(100).Info("info: prioritize result for %s, is %+v", nodeNames, priorityList)
	return &priorityList, nil
}

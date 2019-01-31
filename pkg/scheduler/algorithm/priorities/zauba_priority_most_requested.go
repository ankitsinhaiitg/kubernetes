package priorities

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// Labels zaubaprioritytype:highpriority
//		  zaubaprioritytype:lowpriority

const
(
	resourcePercentangeForLowPriority  int64 = 30
	resourcePercentangeForHighPriority int64 = 70
	zaubaprioritytype string = "zaubaprioritytype"
	highpriority string = "highpriority"
	lowpriority string = "lowpriority"
)

func NewZaubaPriorityMostRequest() *ZaubaPriorityMostRequest {
	return &ZaubaPriorityMostRequest{}
}


type ZaubaPriorityMostRequest struct {
}

func (zpmr ZaubaPriorityMostRequest) ZaubaPriorityMostRequestPriorityMap(
	pod *v1.Pod,
	meta interface{},
	nodeInfo *schedulernodeinfo.NodeInfo) (schedulerapi.HostPriority, error) {
	node := nodeInfo.Node()
	if node == nil {
		return schedulerapi.HostPriority{}, fmt.Errorf("node not found")
	}
	score := 0
	highPrioritySelector := labels.SelectorFromSet(map[string]string{zaubaprioritytype:highpriority})
	lowPrioritySelector := labels.SelectorFromSet(map[string]string{zaubaprioritytype:lowpriority})

	allocatable := nodeInfo.AllocatableResource()
	allocatableHighPriority := calculatePercentageAllocatable(resourcePercentangeForHighPriority, allocatable)
	allocatableLowPriority := calculatePercentageAllocatable(resourcePercentangeForLowPriority, allocatable)


	if nodeInfo.Pods() != nil  &&  len(nodeInfo.Pods()) != 0 {
		for _,pod :=  range nodeInfo.Pods() {
			if highPrioritySelector.Matches(labels.Set(pod.Labels)) {
				score = 10;
				break
			} else if lowPrioritySelector.Matches(labels.Set(pod.Labels)) {
				score = 10;
				break
			}
			fmt.Println(allocatableLowPriority,allocatableHighPriority)
		}
	}

	return schedulerapi.HostPriority{
		Host:  node.Name,
		Score: int(score),
	}, nil
}

func calculatePercentageAllocatable(n int64, allocatable schedulernodeinfo.Resource) schedulernodeinfo.Resource {
	return schedulernodeinfo.Resource{
		Memory:           (allocatable.Memory * n) / 100,
		MilliCPU:         (allocatable.MilliCPU * n) / 100,
		EphemeralStorage: allocatable.EphemeralStorage,
		AllowedPodNumber: allocatable.AllowedPodNumber,
		ScalarResources:  allocatable.ScalarResources,
	}
}

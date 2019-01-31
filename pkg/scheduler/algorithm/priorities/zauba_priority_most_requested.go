package priorities

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/klog"
	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
	schedulernodeinfo "k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// Labels zaubaprioritytype:highpriority
//		  zaubaprioritytype:lowpriority

const
(
	resourcePercentangeForLowPriority  int64  = 30
	resourcePercentangeForHighPriority int64  = 70
	zaubaprioritytype                  string = "zaubaprioritytype"
	highpriority                       string = "highpriority"
	lowpriority                        string = "lowpriority"
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
	} else {
		klog.Infof("Zauba: node:%s \n", node.Name)
	}
	score := 5
	highPrioritySelector := labels.SelectorFromSet(map[string]string{zaubaprioritytype: highpriority})
	lowPrioritySelector := labels.SelectorFromSet(map[string]string{zaubaprioritytype: lowpriority})

	allocatable := nodeInfo.AllocatableResource()
	allocatableHighPriority := calculatePercentageAllocatable(resourcePercentangeForHighPriority, allocatable)
	allocatableLowPriority := calculatePercentageAllocatable(resourcePercentangeForLowPriority, allocatable)

	if nodeInfo.Pods() != nil && len(nodeInfo.Pods()) != 0 {
		klog.Infof("Zauba: Scheduling pod name: %s", pod.Name)
		klog.Infof("Zauba: Number of pods,%v", len(nodeInfo.Pods()))
		for _, pod2 := range nodeInfo.Pods() {
			if highPrioritySelector.Matches(labels.Set(pod2.Labels)) {
				if highPrioritySelector.Matches(labels.Set(pod.Labels)) {
					score = 10;
					klog.Infof("Zauba: High Priority match with selector.")
					break
				} else if lowPrioritySelector.Matches(labels.Set(pod.Labels)) {
					score = 0;
					klog.Infof("Zauba: Low priority pod in high priority node.")
					break
				}
			} else if lowPrioritySelector.Matches(labels.Set(pod2.Labels)) {
				if lowPrioritySelector.Matches(labels.Set(pod.Labels)) {
					score = 10;
					klog.Infof("Zauba: Low priority match with selector.")
					break
				} else if highPrioritySelector.Matches(labels.Set(pod.Labels)) {
					score = 0;
					klog.Infof("Zauba: high priority pod in low priority node.")
					break
				}
			}
			fmt.Println(allocatableLowPriority, allocatableHighPriority)
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

// Copyright © 2016 Eder Ávila Prado <eder.prado@luizalabs.com>
//

package cmd

import (
	"fmt"
	"os"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/resource"
	"k8s.io/kubernetes/pkg/fields"

	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("nodes called")
	},
}

// https://github.com/kubernetes/kubernetes/blob/124fb610dcbd445fa710da67508ac6d5b822f61d/pkg/kubectl/describe.go

// nodesCmd represents the nodes command
var topNodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "Allocated resources over all nodes",
	Run: func(cmd *cobra.Command, args []string) {
		nodes, err := K8sClient.Nodes().List(api.ListOptions{})
		if err != nil {
			fmt.Printf("=> %+v\n", err.Error())
			os.Exit(1)
		}

		var allocatableCPU, allocatableMem float64
		numOfNodes := len(nodes.Items)
		for _, n := range nodes.Items {
			allocatableCPU += float64(n.Status.Capacity.Cpu().MilliValue())
			allocatableMem += float64(n.Status.Capacity.Memory().Value())
		}

		fieldSelector, err := fields.ParseSelector("status.phase!=" + string(api.PodSucceeded) + ",status.phase!=" + string(api.PodFailed))
		if err != nil {
			fmt.Printf("=> %+v\n", err.Error())
			os.Exit(1)
		}
		nodeNonTerminatedPodsList, err := K8sClient.Pods(api.NamespaceAll).List(api.ListOptions{FieldSelector: fieldSelector})
		if err != nil {
			fmt.Printf("=> %+v\n", err.Error())
			os.Exit(1)
		}

		reqs, limits, err := getPodsTotalRequestsAndLimits(nodeNonTerminatedPodsList)
		if err != nil {
			fmt.Printf("=> %+v\n", err)
			os.Exit(1)
		}

		cpuReqs, cpuLimits, memoryReqs, memoryLimits := reqs[api.ResourceCPU], limits[api.ResourceCPU], reqs[api.ResourceMemory], limits[api.ResourceMemory]
		fractionCPUReqs := float64(cpuReqs.MilliValue()) / allocatableCPU * 100
		fractionCPULimits := float64(cpuLimits.MilliValue()) / allocatableCPU * 100
		fractionMemoryReqs := float64(memoryReqs.Value()) / allocatableMem * 100
		fractionMemoryLimits := float64(memoryLimits.Value()) / allocatableMem * 100

		fmt.Printf("Allocated resources around all %d minions:\n", numOfNodes)
		fmt.Printf("  CPU Requests: %s (%d%%)\n", cpuReqs.String(), int64(fractionCPUReqs))
		fmt.Printf("  CPU Limits: %s (%d%%)\n", cpuLimits.String(), int64(fractionCPULimits))
		fmt.Printf("  Memory Requests: %s (%d%%)\n", memoryReqs.String(), int64(fractionMemoryReqs))
		fmt.Printf("  Memory Limits: %s (%d%%)\n", memoryLimits.String(), int64(fractionMemoryLimits))
	},
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// NVIDIA GPU, in devices. Alpha, might change: although fractional and allowing values >1, only one whole device per node is assigned.
	ResourceNvidiaGPU ResourceName = "alpha.kubernetes.io/nvidia-gpu"
	// Number of Pods that may be running on this Node: see ResourcePods
)

func getPodsTotalRequestsAndLimits(podList *api.PodList) (reqs map[api.ResourceName]resource.Quantity, limits map[api.ResourceName]resource.Quantity, err error) {
	reqs, limits = map[api.ResourceName]resource.Quantity{}, map[api.ResourceName]resource.Quantity{}
	for _, pod := range podList.Items {
		podReqs, podLimits, err := api.PodRequestsAndLimits(&pod)
		if err != nil {
			return nil, nil, err
		}
		for podReqName, podReqValue := range podReqs {
			if value, ok := reqs[podReqName]; !ok {
				reqs[podReqName] = *podReqValue.Copy()
			} else {
				value.Add(podReqValue)
				reqs[podReqName] = value
			}
		}
		for podLimitName, podLimitValue := range podLimits {
			if value, ok := limits[podLimitName]; !ok {
				limits[podLimitName] = *podLimitValue.Copy()
			} else {
				value.Add(podLimitValue)
				limits[podLimitName] = value
			}
		}
	}
	return
}

func init() {
	getCmd.AddCommand(nodesCmd)

	topCmd.AddCommand(topNodesCmd)

}

package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/civo/civogo"
	k3dv4 "github.com/rancher/k3d/v4/pkg/client"
	k3dConfig "github.com/rancher/k3d/v4/pkg/config"
	"github.com/rancher/k3d/v4/pkg/config/v1alpha2"
	"github.com/rancher/k3d/v4/pkg/runtimes"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	clusterName = "my-cluster"
	localTpl    = `---
apiVersion: k3d.io/v1alpha2
kind: Simple
name: k3s-civo-default
servers: 1
agents: %d
image: docker.io/rancher/k3s:%s
ports:
  - port: 8080:80
    nodeFilters:
     - loadbalancer
  - port: 8443:443
    nodeFilters:
     - loadbalancer
options:
  k3d:
    wait: true
    timeout: "60s"
  kubeconfig:
    updateDefaultKubeconfig: true
    switchCurrentContext: true
`
)

func main() {
	fmt.Println("Jai Guru")

	//TODO retrieve it from ~/.civo.json
	apiKey := os.Getenv("CIVO_API_KEY")
	region := os.Getenv("CIVO_REGION")

	client, err := civogo.NewClient(apiKey, region)
	if err != nil {
		fmt.Printf("error with client %v", err)
		os.Exit(1)
	}

	k3sCluster, err := client.FindKubernetesCluster(clusterName)
	if err != nil && errors.Is(err, civogo.ZeroMatchesError) {
		fmt.Printf("No cluster found with name %s", clusterName)
		os.Exit(0)
	} else if err != nil {
		fmt.Printf("Error finding cluster \"%v\"", err)
		os.Exit(1)
	}

	var agentMemory int
	if k3sCluster != nil {
		// fmt.Printf("k3s ID :%s with nodes %d", k3sCluster.ID, k3sCluster.NumTargetNode)
		for _, inst := range k3sCluster.Instances {

			if inst.RAMMegabytes > agentMemory {
				agentMemory = inst.RAMMegabytes
			}
		}
	}

	//Start k3d cluster
	ctx := context.TODO()

	//TODO Unable to pull 1.20.0-k3s1 ??
	var k3sImageTag string
	if k3sCluster.KubernetesVersion == "1.20.0-k3s1" {
		k3sImageTag = "v1.20.7-k3s1"
	} else {
		k3sImageTag = "v" + k3sCluster.KubernetesVersion
	}
	str := fmt.Sprintf(localTpl, k3sCluster.NumTargetNode, k3sImageTag)
	var simpleConfig v1alpha2.SimpleConfig
	err = yaml.Unmarshal([]byte(str), &simpleConfig)
	if err != nil {
		fmt.Printf("Error starting cluster \"%v\"", err)
		os.Exit(1)
	}

	simpleConfig.Name = k3sCluster.Name

	cfg, err := k3dConfig.TransformSimpleToClusterConfig(ctx, runtimes.SelectedRuntime, simpleConfig)

	if err != nil {
		fmt.Printf("Error Processing config  \"%v\"", err)
		os.Exit(1)
	}

	if agentMemory != 0 {
		//TODO need to have more robust logic to set each agents memory
		// check if docker has that resources to support
		cfg.ClusterCreateOpts.AgentsMemory = fmt.Sprintf("%dMi", agentMemory)
	}

	err = k3dv4.ClusterRun(ctx, runtimes.SelectedRuntime, cfg)

	if err != nil {
		fmt.Printf("Error starting cluster \"%v\"", err)
		os.Exit(1)
	}
}

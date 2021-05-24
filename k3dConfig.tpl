---
apiVersion: k3d.io/v1alpha2
kind: Simple
name: k3s-civo-default
servers: 1
agents: {{ .NumTargetNodes }}
image: docker.io/rancher/k3s:{{ .K3sVersion }}
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
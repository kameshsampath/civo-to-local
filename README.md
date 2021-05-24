# Civo Remote to Local (POC)

A simple POC to create a remote [Civo](https://civo.io) cluster locally. Which can then be used to 

- work offline
- add apps
- scale 
- sync back when online

## Prerequisite

- [k3d](https://k3d.io/) 
- [Civo](https://civo.io) Account

## Running the PoC
```shell
 civo k3s create my-cluster --size=g3.k3s.small --nodes=3 --save --merge
```

Run the app:

```shell
export CIVO_API_KEY=<your civo api key>
export CIVO_REGION=<your civo region>
go run main.go 
```

Check `k3d` cluster locally

```shell
k3d cluster ls
```

```text
NAME         SERVERS   AGENTS   LOADBALANCER
my-cluster   1/1       3/3      true
```
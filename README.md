# gooom

Check if containers in Kubernetes Pods were OOMKilled

## Running

```
$ make dep build

$ ./gooom -h
Usage of ./gooom:
  -duration string
        (optional) duration before now to check (default "30m")
  -kubeconfig string
        (optional) absolute path to the kubeconfig file (default "${HOME}/.kube/config")
  -namespace string
        (optional) specific namespace to check (default checks all)
  -timeout string
        (optional) timeout (default "10s")
```

## Examples

Check all namespaces for OOMKilled containers in the last 30 minutes

```
$ ./gooom
2020/11/27 10:24:26 Checking for OOMKilled Containers in the last 30m0s
2020/11/27 10:24:26 Checking   2 pods in namespace: default
2020/11/27 10:24:26 Checking  21 pods in namespace: kafka
2020/11/27 10:24:26 Checking  40 pods in namespace: logging
2020/11/27 10:24:26 Checking  39 pods in namespace: prometheus
[...]
2020/11/27 10:24:32 Done
```

Check for OOMKilled containers in the 'prometheus' namespace in the last week

```
$ ./gooom -namespace prometheus -duration 168h
2020/11/27 10:23:57 Checking for OOMKilled Containers in the last 168h0m0s
2020/11/27 10:23:57 Checking  57 pods in namespace: prometheus
2020/11/27 10:23:57 Container in thanos-store-95d598ff7-qvzzq was OOMKilled at 2020-11-23 01:39:01 +0100 CET
2020/11/27 10:23:57 Done
```

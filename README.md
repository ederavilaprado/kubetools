# kubetools

Kubernetes helpers for troubleshooting and administration.


## Commands

Display allocated resources around all kubernetes minions...
```
$ kubetools top nodes
Allocated resources around all 80 minions:
  CPU Requests: 209351m (50%)
  CPU Limits: 193151m (49%)
  Memory Requests: 273217396Ki (19%)
  Memory Limits: 281708404Ki (20%)
```

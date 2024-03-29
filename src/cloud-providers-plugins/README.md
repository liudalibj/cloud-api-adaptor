# Example project about how to use the runtime cloud provider plugin feature with CAA and peerpod-ctrl

Follow [addnewprovider.md](../cloud-api-adaptor/docs/addnewprovider.md) we can add a new clour provider to CAA in complie time.

In this document we will add a new cloud provider via the go plugin that is load a new cloud provider `libvirt-ext` from CAA runtime.


### Step 1: Initialize and register the cloud provider manager for `libvirt-ext`

:information_source:[manager code](./libvirt-ext/manager.go)

### Step 2: Add provider specific code

:information_source:[provider code](./libvirt-ext/provider.go)

### Step 3: Build related docker images

- Build cloud provider plugin docker image
:information_source:[plugins docker image code](./Dockerfile.plugins)
```bash
registry=quay.io/<your-test-repo> && make image-plugins
```
- Include `libvirt-ext` plugin to CAA image
:information_source:[plugins docker image code](./Dockerfile.caa)
```bash
export registry=quay.io/<your-test-repo>
registry=quay.io/<your-test-repo> && make image-ppc
```
- Include `libvirt-ext` plugin to Peerpod-ctrl image
:information_source:[plugins docker image code](./Dockerfile.caa)
```bash
registry=quay.io/<your-test-repo> && make image-caa
```

### Step 4: Using the new CAA and Peerpod-ctrl in the peerpod env
- Config one test peerpod env via [document](../cloud-api-adaptor/libvirt/README.md)

- Check/update peer-pods-cm configmap
```bash
kubectl get cm -n confidential-containers-system   peer-pods-cm -o yaml
apiVersion: v1
data:
  CLOUD_CONFIG_VERIFY: "false"
  CLOUD_PROVIDER: libvirt-ext
  CLOUD_PROVIDER_PLUGIN_PATH: /cloud-providers
  DISABLECVM: "true"
  LIBVIRT_NET: default
  LIBVIRT_POOL: default
  LIBVIRT_URI: qemu+ssh://root@192.168.122.1/system?no_verify=1
  LIBVIRT_VOL_NAME: podvm-base.qcow2
```
> **Note** two items need be update:
> - `CLOUD_PROVIDER` to `libvirt-ext`
> - Add `CLOUD_PROVIDER_PLUGIN_PATH: /cloud-providers` if it's not there
- Update CAA image to the built out CAA and peerpod-ctrl images
```bash
kubectl set image -n confidential-containers-system   ds/cloud-api-adaptor-daemonset cloud-api-adaptor-con=${registry}/cloud-api-adaptor:dev-plugins-${COMMIT_ID}

kubectl set image -n confidential-containers-system deploy/peerpod-ctrl-controller-manage manager=${registry}/peerpod-ctrl:dev-plugins-${COMMIT_ID}
```
- Create peerpod to verify the update env
```bash
kubectl create -f - <<EOF
apiVersion: v1
kind: Pod
metadata:
  labels:
    run: busybox
  name: busybox
spec:
  containers:
  - image: quay.io/prometheus/busybox
    name: busybox
    resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Never
  runtimeClassName: kata-remote
EOF
```
- Check CAA log:
```bash

```



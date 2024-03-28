# Example project about how to use the runtime cloud provider plugin feature with CAA and peerpod-ctrl

Follow [addnewprovider.md](../cloud-api-adaptor/docs/addnewprovider.md) we can add a new clour provider to CAA in complie time.

In this document we will add a new cloud provider via the go plugin that is load a new cloud provider `libvirt-ext` from CAA runtime.

The demo `libvirt-ext` cloud provider is base on the built-in cloud provider `libvirt`, just make the `cloud-init-data` for new peerpod vm instance can be checked from CAA logs.

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
- Update CAA daemonset image and peerpod-ctrl deploy to the built out CAA and peerpod-ctrl images
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
kubectl logs -n confidential-containers-system ds/cloud-api-adaptor-daemonset
+ exec cloud-api-adaptor libvirt-ext -uri 'qemu+ssh://root@192.168.122.1/system?no_verify=1' -data-dir /opt/data-dir -pods-dir /run/peerpod/pods -network-name default -pool-name default -disable-cvm -socket /run/peerpod/hypervisor.sock
2024/03/30 03:48:33 [adaptor/cloud] Loading cloud providers from /cloud-providers
2024/03/30 03:48:33 [adaptor/cloud] Found cloud provider file /cloud-providers/libvirt-ext.so
cloud-api-adaptor version v0.9.0.alpha.1-dev
  commit: 840361d3832b7c96956a97f5e084f6637336e66e
  go: go1.21.8
cloud-api-adaptor: starting Cloud API Adaptor daemon for "libvirt-ext"
2024/03/30 03:48:33 [adaptor/cloud/libvirt] libvirt config: &libvirt.Config{URI:"qemu+ssh://root@192.168.122.1/system?no_verify=1", PoolName:"default", NetworkName:"default", DataDir:"/opt/data-dir", DisableCVM:true, VolName:"podvm-base.qcow2", LaunchSecurity:"", Firmware:"/usr/share/edk2/ovmf/OVMF_CODE.fd"}
2024/03/30 03:48:40 [adaptor/cloud/libvirt] Created libvirt connection
2024/03/30 03:48:40 [adaptor] server config: &adaptor.ServerConfig{TLSConfig:(*tlsutil.TLSConfig)(0xc000810000), SocketPath:"/run/peerpod/hypervisor.sock", CriSocketPath:"", PauseImage:"", PodsDir:"/run/peerpod/pods", ForwarderPort:"15150", ProxyTimeout:300000000000, AAKBCParams:"", EnableCloudConfigVerify:false}
2024/03/30 03:48:40 [util/k8sops] initialized PeerPodService
2024/03/30 03:48:40 [probe/probe] Using port: 8000
2024/03/30 03:48:40 [adaptor] server started
2024/03/30 03:48:55 [probe/probe] nodeName: peer-pods-worker-0
2024/03/30 03:48:55 [probe/probe] Selected pods count: 10
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: cc-operator-controller-manager-857f844f7d-q7zkc
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: cc-operator-daemon-install-86xs5
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: cc-operator-pre-install-daemon-9z9kd
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: cloud-api-adaptor-daemonset-j2n85
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: peerpod-ctrl-controller-manager-764bfd69b8-5b9rv
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: ingress-nginx-admission-create-zmljg
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: ingress-nginx-admission-patch-5mmjf
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: ingress-nginx-controller-7bf7bc78dc-f9l7s
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: kube-flannel-ds-pqff7
2024/03/30 03:48:55 [probe/probe] Ignored standard pod: kube-proxy-zjx75
2024/03/30 03:48:55 [probe/probe] All PeerPods standup. we do not check the PeerPods status any more.
2024/03/30 03:48:55 [podnetwork] routes on netns /var/run/netns/cni-54e35cd5-6fc3-8dfa-a8b9-2b1621fa51f0
2024/03/30 03:48:55 [podnetwork]     0.0.0.0/0 via 10.244.1.1 dev eth0
2024/03/30 03:48:55 [podnetwork]     10.244.0.0/16 via 10.244.1.1 dev eth0
2024/03/30 03:48:55 [adaptor/cloud] Credentials file is not in a valid Json format, ignored
2024/03/30 03:48:55 [adaptor/cloud] stored /run/peerpod/pods/cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4/daemon.json
2024/03/30 03:48:55 [adaptor/cloud] create a sandbox cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4 for pod busybox in namespace default (netns: /var/run/netns/cni-54e35cd5-6fc3-8dfa-a8b9-2b1621fa51f0)
2024/03/30 03:48:55 [adaptor/cloud/libvirt-ext] ===CreateInstance: userData from libvirt-ext: #cloud-config

write_files:
  - path: /run/peerpod/daemon.json
    content: |
      {
          "pod-network": {
              "podip": "10.244.1.27/24",
              "pod-hw-addr": "42:68:f7:5c:4b:2c",
              "interface": "eth0",
              "worker-node-ip": "192.168.122.25/24",
              "tunnel-type": "vxlan",
              "routes": [
                  {
                      "Dst": "0.0.0.0/0",
                      "GW": "10.244.1.1",
                      "Dev": "eth0"
                  },
                  {
                      "Dst": "10.244.0.0/16",
                      "GW": "10.244.1.1",
                      "Dev": "eth0"
                  }
              ],
              "mtu": 1450,
              "index": 0,
              "vxlan-port": 4789,
              "vxlan-id": 555000,
              "dedicated": false
          },
          "pod-namespace": "default",
          "pod-name": "busybox",
          "tls-server-key": "-----BEGIN PRIVATE KEY-----......-----END PRIVATE KEY-----\n",
          "tls-server-cert": "-----BEGIN CERTIFICATE-----......-----END CERTIFICATE-----\n",
          "tls-client-ca": "-----BEGIN CERTIFICATE-----......-----END CERTIFICATE-----\n"
      }
2024/03/30 03:48:55 [adaptor/cloud/libvirt] LaunchSecurityType: None
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Checking if instance (podvm-busybox-cabec4dc) exists
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Uploaded volume key /var/lib/libvirt/images/podvm-busybox-cabec4dc-root.qcow2
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Create cloudInit iso
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Uploading iso file: podvm-busybox-cabec4dc-cloudinit.iso
2024/03/30 03:48:55 [adaptor/cloud/libvirt] 45056 bytes uploaded
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Volume ID: /var/lib/libvirt/images/podvm-busybox-cabec4dc-cloudinit.iso
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Create XML for 'podvm-busybox-cabec4dc'
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Creating VM 'podvm-busybox-cabec4dc'
2024/03/30 03:48:55 [adaptor/cloud/libvirt] Starting VM 'podvm-busybox-cabec4dc'
2024/03/30 03:48:57 [adaptor/cloud/libvirt] VM id 344
2024/03/30 03:49:18 [adaptor/cloud/libvirt] Instance created successfully
2024/03/30 03:49:18 [adaptor/cloud/libvirt] created an instance podvm-busybox-cabec4dc for sandbox cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [util/k8sops] busybox is now owning a PeerPod object
2024/03/30 03:49:18 [adaptor/cloud] created an instance podvm-busybox-cabec4dc for sandbox cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [tunneler/vxlan] vxlan ppvxlan1 (remote 192.168.122.158:4789, id: 555000) created at /proc/1/task/1/ns/net
2024/03/30 03:49:18 [tunneler/vxlan] vxlan ppvxlan1 created at /proc/1/task/1/ns/net
2024/03/30 03:49:18 [tunneler/vxlan] vxlan ppvxlan1 is moved to /var/run/netns/cni-54e35cd5-6fc3-8dfa-a8b9-2b1621fa51f0
2024/03/30 03:49:18 [tunneler/vxlan] Add tc redirect filters between eth0 and vxlan1 on pod network namespace /var/run/netns/cni-54e35cd5-6fc3-8dfa-a8b9-2b1621fa51f0
2024/03/30 03:49:18 [adaptor/proxy] Listening on /run/peerpod/pods/cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4/agent.ttrpc
2024/03/30 03:49:18 [adaptor/proxy] failed to init cri client, the err: cri runtime endpoint is not specified, it is used to get the image name from image digest
2024/03/30 03:49:18 [adaptor/proxy] Trying to establish agent proxy connection to 192.168.122.158:15150
2024/03/30 03:49:18 [adaptor/proxy] established agent proxy connection to 192.168.122.158:15150
2024/03/30 03:49:18 [adaptor/cloud] agent proxy is ready
2024/03/30 03:49:18 [adaptor/proxy] CreateSandbox: hostname:busybox sandboxId:cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [adaptor/proxy]     storages:
2024/03/30 03:49:18 [adaptor/proxy]         mountpoint:/run/kata-containers/sandbox/shm source:shm fstype:tmpfs driver:ephemeral
2024/03/30 03:49:18 [adaptor/proxy] CreateContainer: containerID:cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [adaptor/proxy]     mounts:
2024/03/30 03:49:18 [adaptor/proxy]         destination:/proc source:proc type:proc
2024/03/30 03:49:18 [adaptor/proxy]         destination:/dev source:tmpfs type:tmpfs
2024/03/30 03:49:18 [adaptor/proxy]         destination:/dev/pts source:devpts type:devpts
2024/03/30 03:49:18 [adaptor/proxy]         destination:/dev/mqueue source:mqueue type:mqueue
2024/03/30 03:49:18 [adaptor/proxy]         destination:/sys source:sysfs type:sysfs
2024/03/30 03:49:18 [adaptor/proxy]         destination:/dev/shm source:/run/kata-containers/sandbox/shm type:bind
2024/03/30 03:49:18 [adaptor/proxy]         destination:/etc/resolv.conf source:/run/kata-containers/shared/containers/cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4-872c9aff001f163c-resolv.conf type:bind
2024/03/30 03:49:18 [adaptor/proxy]     annotations:
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-cpu-period: 100000
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-memory: 0
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-log-directory: /var/log/pods/default_busybox_a079ed8d-359d-4cf7-b24f-364927e848df
2024/03/30 03:49:18 [adaptor/proxy]         io.katacontainers.pkg.oci.container_type: pod_sandbox
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.container-type: sandbox
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-cpu-shares: 2
2024/03/30 03:49:18 [adaptor/proxy]         nerdctl/network-namespace: /var/run/netns/cni-54e35cd5-6fc3-8dfa-a8b9-2b1621fa51f0
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-id: cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-name: busybox
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-cpu-quota: 0
2024/03/30 03:49:18 [adaptor/proxy]         io.katacontainers.pkg.oci.bundle_path: /run/containerd/io.containerd.runtime.v2.task/k8s.io/cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-namespace: default
2024/03/30 03:49:18 [adaptor/proxy]         io.kubernetes.cri.sandbox-uid: a079ed8d-359d-4cf7-b24f-364927e848df
2024/03/30 03:49:18 [adaptor/proxy] getImageName: no pause image specified uses default pause image: registry.k8s.io/pause:3.7
2024/03/30 03:49:18 [adaptor/proxy] CreateContainer: calling PullImage for "registry.k8s.io/pause:3.7" before CreateContainer (cid: "cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4")
2024/03/30 03:49:19 [adaptor/proxy] CreateContainer: successfully pulled image "registry.k8s.io/pause:3.7"
2024/03/30 03:49:19 [adaptor/proxy] StartContainer: containerID:cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:24 [adaptor/proxy] CreateContainer: containerID:09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce
2024/03/30 03:49:24 [adaptor/proxy]     mounts:
2024/03/30 03:49:24 [adaptor/proxy]         destination:/proc source:proc type:proc
2024/03/30 03:49:24 [adaptor/proxy]         destination:/dev source:tmpfs type:tmpfs
2024/03/30 03:49:24 [adaptor/proxy]         destination:/dev/pts source:devpts type:devpts
2024/03/30 03:49:24 [adaptor/proxy]         destination:/dev/mqueue source:mqueue type:mqueue
2024/03/30 03:49:24 [adaptor/proxy]         destination:/sys source:sysfs type:sysfs
2024/03/30 03:49:24 [adaptor/proxy]         destination:/sys/fs/cgroup source:cgroup type:cgroup
2024/03/30 03:49:24 [adaptor/proxy]         destination:/etc/hosts source:/run/kata-containers/shared/containers/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce-585f91ac859be007-hosts type:bind
2024/03/30 03:49:24 [adaptor/proxy]         destination:/dev/termination-log source:/run/kata-containers/shared/containers/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce-672bdebb31331291-termination-log type:bind
2024/03/30 03:49:24 [adaptor/proxy]         destination:/etc/hostname source:/run/kata-containers/shared/containers/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce-5452912e1a3fdacb-hostname type:bind
2024/03/30 03:49:24 [adaptor/proxy]         destination:/etc/resolv.conf source:/run/kata-containers/shared/containers/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce-4118ab95f3e3f8fe-resolv.conf type:bind
2024/03/30 03:49:24 [adaptor/proxy]         destination:/dev/shm source:/run/kata-containers/sandbox/shm type:bind
2024/03/30 03:49:24 [adaptor/proxy]         destination:/var/run/secrets/kubernetes.io/serviceaccount source:/run/kata-containers/shared/containers/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce-764748773cab5307-serviceaccount type:bind
2024/03/30 03:49:24 [adaptor/proxy]     annotations:
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.sandbox-id: cabec4dcd78f3fb4fb6765b6e9822e583d442177b8546cf0d8b182d5776e45b4
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.sandbox-name: busybox
2024/03/30 03:49:24 [adaptor/proxy]         io.katacontainers.pkg.oci.bundle_path: /run/containerd/io.containerd.runtime.v2.task/k8s.io/09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.container-name: busybox
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.container-type: container
2024/03/30 03:49:24 [adaptor/proxy]         io.katacontainers.pkg.oci.container_type: pod_container
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.sandbox-uid: a079ed8d-359d-4cf7-b24f-364927e848df
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.sandbox-namespace: default
2024/03/30 03:49:24 [adaptor/proxy]         io.kubernetes.cri.image-name: quay.io/prometheus/busybox:latest
2024/03/30 03:49:24 [adaptor/proxy] getImageName: got image from annotations: quay.io/prometheus/busybox:latest
2024/03/30 03:49:24 [adaptor/proxy] CreateContainer: calling PullImage for "quay.io/prometheus/busybox:latest" before CreateContainer (cid: "09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce")
2024/03/30 03:49:26 [adaptor/proxy] CreateContainer: successfully pulled image "quay.io/prometheus/busybox:latest"
2024/03/30 03:49:26 [adaptor/proxy] StartContainer: containerID:09c07739c2fa8e2d54242b52597b0b7f42973c93c7ef4bddffce5c71189debce
```
- Check Peerpod-ctrl log:
```bash
kubectl logs -f -n confidential-containers-system   peerpod-ctrl-controller-manager-658f4d9ff8-qs5dj
2024-03-30T03:46:20Z	INFO	controller-runtime.metrics	Metrics server is starting to listen	{"addr": "127.0.0.1:8080"}
2024/03/30 03:46:20 [adaptor/cloud] Loading cloud providers from /cloud-providers
2024/03/30 03:46:20 [adaptor/cloud] Found cloud provider file /cloud-providers/libvirt-ext.so
2024/03/30 03:46:20 [adaptor/cloud/libvirt] libvirt config: &libvirt.Config{URI:"qemu+ssh://root@192.168.122.1/system?no_verify=1", PoolName:"default", NetworkName:"default", DataDir:"", DisableCVM:false, VolName:"podvm-base.qcow2", LaunchSecurity:"", Firmware:"/usr/share/edk2/ovmf/OVMF_CODE.fd"}
2024/03/30 03:46:20 [adaptor/cloud/libvirt] Created libvirt connection
2024-03-30T03:46:20Z	INFO	setup	starting manager
2024-03-30T03:46:20Z	INFO	Starting server	{"kind": "health probe", "addr": "[::]:8081"}
2024-03-30T03:46:20Z	INFO	Starting server	{"path": "/metrics", "kind": "metrics", "addr": "127.0.0.1:8080"}
I0330 03:46:21.015004       1 leaderelection.go:248] attempting to acquire leader lease confidential-containers-system/33f6c5d6.confidentialcontainers.org...
I0330 03:46:52.021599       1 leaderelection.go:258] successfully acquired lease confidential-containers-system/33f6c5d6.confidentialcontainers.org
2024-03-30T03:46:52Z	INFO	Starting EventSource	{"controller": "peerpod", "controllerGroup": "confidentialcontainers.org", "controllerKind": "PeerPod", "source": "kind source: *v1alpha1.PeerPod"}
2024-03-30T03:46:52Z	INFO	Starting Controller	{"controller": "peerpod", "controllerGroup": "confidentialcontainers.org", "controllerKind": "PeerPod"}
2024-03-30T03:46:52Z	DEBUG	events	peerpod-ctrl-controller-manager-658f4d9ff8-qs5dj_77f3fdad-419c-433f-9b58-3eff5cb2a756 became leader	{"type": "Normal", "object": {"kind":"Lease","namespace":"confidential-containers-system","name":"33f6c5d6.confidentialcontainers.org","uid":"8f291c13-e585-42e1-b8df-c13aa86b566b","apiVersion":"coordination.k8s.io/v1","resourceVersion":"915877"}, "reason": "LeaderElection"}
2024-03-30T03:46:52Z	INFO	Starting workers	{"controller": "peerpod", "controllerGroup": "confidentialcontainers.org", "controllerKind": "PeerPod", "worker count": 1}
```
## Troubleshooting
- "failed to map segment from shared object" from CAA/Peerpod-ctrl log
> - Please use update `CLOUD_PROVIDER_PLUGIN_PATH` to a path which have execute permissions, the plugin .so file need have execute permissions

- "plugin was built with a different version of package XXX" from CAA/Peerpod-ctrl log
> - Please check the go.mod of CAA and plugins project, the CAA and plugins should be built with same version of issue package XXX
> - Please make sure use same golang env to build CAA, Peerpod-ctrl and cloud-provider plugins

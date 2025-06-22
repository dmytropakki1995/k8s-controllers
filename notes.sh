mkdir -p ./kubebuilder/bin && \
    curl -L https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.30.0-linux-arm64.tar.gz -o kubebuilder-tools.tar.gz && \
    tar -C ./kubebuilder --strip-components=1 -zvxf kubebuilder-tools.tar.gz && \
    rm kubebuilder-tools.tar.gz

kubebuilder/bin/kubectl config set-credentials test-user --token=1234567890
kubebuilder/bin/kubectl config set-cluster test-env --server=https://127.0.0.1:6443 --insecure-skip-tls-verify
kubebuilder/bin/kubectl config set-context test-context --cluster=test-env --user=test-user --namespace=default 
kubebuilder/bin/kubectl config use-context test-context


echo "Starting etcd..."
kubebuilder/bin/etcd \
    --advertise-client-urls http://$HOST_IP:2379 \
    --listen-client-urls http://0.0.0.0:2379 \
    --data-dir ./etcd \
    --listen-peer-urls http://0.0.0.0:2380 \
    --initial-cluster default=http://$HOST_IP:2380 \
    --initial-advertise-peer-urls http://$HOST_IP:2380 \
    --initial-cluster-state new \
    --initial-cluster-token test-token &

echo "Starting kube-apiserver..."
sudo kubebuilder/bin/kube-apiserver \
    --etcd-servers=http://$HOST_IP:2379 \
    --service-cluster-ip-range=10.0.0.0/24 \
    --bind-address=0.0.0.0 \
    --secure-port=6443 \
    --advertise-address=$HOST_IP \
    --authorization-mode=AlwaysAllow \
    --token-auth-file=/tmp/token.csv \
    --enable-priority-and-fairness=false \
    --allow-privileged=true \
    --profiling=false \
    --storage-backend=etcd3 \
    --storage-media-type=application/json \
    --v=0 \
    --service-account-issuer=https://kubernetes.default.svc.cluster.local \
    --service-account-key-file=/tmp/sa.pub \
    --service-account-signing-key-file=/tmp/sa.key &


export PATH=$PATH:/opt/cni/bin:kubebuilder/bin
echo "Starting containerd..."
PATH=$PATH:/opt/cni/bin:/usr/sbin /opt/cni/bin/containerd -c /etc/containerd/config.toml &

echo "Starting kube-scheduler..."
sudo kubebuilder/bin/kube-scheduler \
    --kubeconfig=/root/.kube/config \
    --leader-elect=false \
    --v=2 \
    --bind-address=0.0.0.0 &

echo "Creating kubelet directories..."
mkdir -p /var/lib/kubelet
mkdir -p /etc/kubernetes/manifests
mkdir -p /var/log/kubernetes

cat << EOF | tee /var/lib/kubelet/config.yaml
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
authentication:
  anonymous:
    enabled: true
  webhook:
    enabled: true
  x509:
    clientCAFile: "/var/lib/kubelet/ca.crt"
authorization:
  mode: AlwaysAllow
clusterDomain: "cluster.local"
clusterDNS:
  - "10.0.0.10"
resolvConf: "/etc/resolv.conf"
runtimeRequestTimeout: "15m"
failSwapOn: false
seccompDefault: true
serverTLSBootstrap: true
containerRuntimeEndpoint: "unix:///run/containerd/containerd.sock"
staticPodPath: "/etc/kubernetes/manifests"
EOF


cp /root/.kube/config /var/lib/kubelet/kubeconfig
export KUBECONFIG=~/.kube/config
cp /tmp/sa.pub /tmp/ca.crt

echo "Starting kubelet..."
sudo PATH=$PATH:/opt/cni/bin:/usr/sbin kubebuilder/bin/kubelet \
    --kubeconfig=/var/lib/kubelet/kubeconfig \
    --config=/var/lib/kubelet/config.yaml \
    --root-dir=/var/lib/kubelet \
    --cert-dir=/var/lib/kubelet/pki \
    --hostname-override=$(hostname)\
    --pod-infra-container-image=registry.k8s.io/pause:3.10 \
    --node-ip=$HOST_IP \
    --cgroup-driver=cgroupfs \
    --max-pods=4  \
    --v=1 &

PATH=$PATH:/usr/sbin kubebuilder/bin/kubectl apply -f -<<EOF
apiVersion: v1
kind: Pod
metadata:
  name: test-pod-2
spec:
  containers:
    - name: test-container-nginx
      image: nginx:1.21
      securityContext:
        privileged: true
EOF
ERRO[2025-06-18T16:00:18.856656381Z] RunPodSandbox for &PodSandboxMetadata{Name:test-pod-2,Uid:77232677-eb34-4a56-aeb2-ed493b0ec59e,Namespace:default,Attempt:0,} failed, error  error="rpc error: code = InvalidArgument desc = failed to start sandbox \"116108124acdfc5b9012abcc32dce046139ee2d92583aa509f51abe6c1b084a2\": failed to get sandbox image \"registry.k8s.io/pause:3.10\": failed to pull image \"registry.k8s.io/pause:3.10\": failed to pull and unpack image \"registry.k8s.io/pause:3.10\": unable to initialize unpacker: no unpack platforms defined: invalid argument"
E0618 16:00:18.857288    3802 remote_runtime.go:193] "RunPodSandbox from runtime service failed" err="rpc error: code = InvalidArgument desc = failed to start sandbox \"116108124acdfc5b9012abcc32dce046139ee2d92583aa509f51abe6c1b084a2\": failed to get sandbox image \"registry.k8s.io/pause:3.10\": failed to pull image \"registry.k8s.io/pause:3.10\": failed to pull and unpack image \"registry.k8s.io/pause:3.10\": unable to initialize unpacker: no unpack platforms defined: invalid argument"



RunPodSandbox for &PodSandboxMetadata{Name:test-pod-2,Uid:88029ce9-dffa-4898-8f00-0401de0c1571,Namespace:default,Attempt:0,} failed, error  error="rpc error: code = InvalidArgument desc = failed to start sandbox \"5ec33c9be02afe4dcd6e2897d11391ab258183ee253486307a668e2b4ad2a37a\": failed to get sandbox image \"registry.k8s.io/pause:3.10\": failed to pull image \"registry.k8s.io/pause:3.10\": failed to pull and unpack image \"registry.k8s.io/pause:3.10\": unable to initialize unpacker: no unpack platforms defined: invalid argument"

cat <<EOF > config.toml
version = 3

[grpc]
  address = "/run/containerd/containerd.sock"

[plugins.'io.containerd.cri.v1.runtime']
  enable_selinux = false
  enable_unprivileged_ports = true
  enable_unprivileged_icmp = true
  device_ownership_from_security_context = false

[plugins.'io.containerd.cri.v1.images']
  snapshotter = "native"
  disable_snapshot_annotations = true

[plugins.'io.containerd.cri.v1.runtime'.cni]
  bin_dir = "/opt/cni/bin"
  conf_dir = "/etc/cni/net.d"

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runc]
  runtime_type = "io.containerd.runc.v2"

[plugins.'io.containerd.cri.v1.runtime'.containerd.runtimes.runc.options]
  SystemdCgroup = false

[plugins.'io.containerd.grpc.v1.cri'.containerd]
  default_platform = "linux/arm64"

EOF
sudo mv config.toml /etc/containerd/config.toml


sudo ctr images pull registry.k8s.io/pause:3.10

kubebuilder/bin/kubectl delete pod test-pod-2
kubebuilder/bin/kubectl get pods

# Detect and set the system architecture
ARCH=$(uname -m | sed -e 's/^x86_64$/amd64/' -e 's/^aarch64$/arm64/')

mkdir -p ./kubebuilder/bin && \
    curl -L https://storage.googleapis.com/kubebuilder-tools/kubebuilder-tools-1.30.0-linux-${ARCH}.tar.gz -o kubebuilder-tools.tar.gz && \
    tar -C ./kubebuilder --strip-components=1 -zvxf kubebuilder-tools.tar.gz && \
    rm kubebuilder-tools.tar.gz


wget https://github.com/containerd/containerd/releases/download/v2.0.5/containerd-static-2.0.5-linux-arm64.tar.gz
sudo apt-get install iptables

sudo /opt/cni/bin/ctr -n k8s.io c ls
sudo /opt/cni/bin/ctr -n k8s.io tasks exec -t --exec-id m 744d76d02a33de4f4724d5749bfdde9ea19f6298fc727d1cccb490bf41ec86e3 sh
# Kubernetes Controller (Golang)

This simple Go program uses the official Kubernetes Go client (`client-go`) to create and delete a Deployment resource by reading a manifest file (`deployment.yaml`).


## 🧰 Requirements

- Go 1.24 or later
- Access to a Kubernetes cluster Control Plane
- A valid `kubeconfig` file (default: `~/.kube/config`)
- A Deployment manifest file in YAML format



## 🛠 Installation and Usage
```sh
git switch feature/step6-list-deployments
go mod init && go mod tidy

# Create deployment
go run main.go --kubeconfig ./config --manifest-path ./deployment.yaml create

# Delete deployment
go run main.go --kubeconfig ./config --deployment my-deployment delete
```

**What it does:**
- Read a Kubernetes Deployment manifest from a YAML file
- Create the Deployment in your cluster using the Kubernetes API
- Supports kubeconfig from `~/.kube/config`
- Delete deployment with the specified name




## License

MIT License. See [LICENSE](LICENSE) for details.
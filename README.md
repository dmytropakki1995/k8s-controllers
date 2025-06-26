# Kubernetes Controller (Golang)

This simple Go program uses the official Kubernetes Go client (`client-go`) to monitor events in the kubernetes cluster using informer (`k8s.io/client-go/informers`).


## 🧰 Requirements

- Go 1.24 or later
- Access to a Kubernetes cluster Control Plane
- A valid `kubeconfig` file (default: `~/.kube/config`)
- A Deployment manifest file in YAML format



## 🛠 Installation and Usage
```sh
git switch feature/step6-list-deployments
go mod init && go mod tidy

# Start FastHTTP server to monitor k8s events
go run main.go --log-level trace --kubeconfig ./config --resource-types pod,deployment --namespace default server

# Delete deployment
go run main.go --kubeconfig ./config --deployment my-deployment delete
```

**What it does:**
- Connects to the Kubernetes cluster using the provided kubeconfig file or in-cluster config.
- Watches for k8s resource events (add, update, delete) in the specified namespace for the specified resource types in the `--resource-types` flag


## Project Structure
- `cmd/` — Contains your CLI commands.
- `cmd/server.go` - fasthttp server
- `pkg/informer` - informer implementation
- `pkg/testutil` - envtest kit
- `main.go` — Entry point for your application.
- `Makefile` — Build automation tasks.
- `.env` - Default values for command flags


## License
MIT License. See [LICENSE](LICENSE) for details.

# make sure if minikube is installed
minikube start

# switch to docker daemon of minikube
eval $(minikube docker-env)

# Build docker image
docker build -t alpacked/test-task-build .

# switch back to local docker daemon
eval $(minikube docker-env -u)

# Deply application
kubectl apply -f config.yaml


########
# ------ Make some tests
kubectl port-forward svc/debug-task-service 8081:8081

# Upload file to the server
curl -X POST -F "file=@test.txt" http://localhost:8081/upload

# Download file from the server
curl http://localhost:8081/uploads/test.txt
# Download file from the server that does not exist
curl http://localhost:8081/uploads/test_copy.txt

# Make sure that persistent volumes work as expected
kubectl get pods && kubectl delete pod <pod_id>

# Check if the file still exists and has the same content
kubectl exec -it debug-task-deployment-69c9595ccc-nx6kf -c app -- cat /app/uploads/test.txt

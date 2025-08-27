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

# Should be used to allocate a virtual IP for the LoadBalancer
minikube tunnel


########
# ------ Make some tests
# kubectl port-forward svc/nginx-service 8081:80

# Upload file to the server
curl -X POST -F "file=@test.txt" http://localhost/upload
# Upload one more time to check if file already exists
curl -X POST -F "file=@test.txt" http://localhost/upload

# Download file from the server
curl http://localhost/uploads/test.txt
# Download file from the server that does not exist
curl http://localhost/uploads/test_copy.txt

# Make sure that persistent volumes work as expected
kubectl get pods && kubectl delete pod debug-task-deployment-<pod_id>

# Check if the file still exists and has the same content
kubectl exec -it debug-task-deployment-<pod_id> -c app -- cat /app/uploads/test.txt

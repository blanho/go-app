set -e

RESOURCE_GROUP="azure-go-app-prod-rg"
AKS_CLUSTER="azure-go-app-aks"
ACR_NAME="azuregoappacr"
APP_NAME="azure-go-app"
NAMESPACE="production"

echo "Logging in to Azure..."
az login --service-principal --username $SP_CLIENT_ID --password $SP_CLIENT_SECRET --tenant $SP_TENANT_ID

echo "Setting subscription..."
az account set --subscription $SUBSCRIPTION_ID

echo "Getting AKS credentials..."
az aks get-credentials --resource-group $RESOURCE_GROUP --name $AKS_CLUSTER --admin

echo "Creating namespace if it doesn't exist..."
kubectl create namespace $NAMESPACE --dry-run=client -o yaml | kubectl apply -f -

echo "Building and pushing Docker image..."
VERSION=$(git rev-parse --short HEAD)
docker build -t $ACR_NAME.azurecr.io/$APP_NAME:$VERSION -t $ACR_NAME.azurecr.io/$APP_NAME:latest .
az acr login --name $ACR_NAME
docker push $ACR_NAME.azurecr.io/$APP_NAME:$VERSION
docker push $ACR_NAME.azurecr.io/$APP_NAME:latest

echo "Updating deployment manifest..."
sed -i "s|image: .*|image: $ACR_NAME.azurecr.io/$APP_NAME:$VERSION|" ./deploy/kubernetes/deployment.yaml

echo "Applying Kubernetes manifests..."
kubectl apply -f ./deploy/kubernetes/

echo "Waiting for deployment to complete..."
kubectl rollout status deployment/$APP_NAME -n $NAMESPACE

echo "Deployment completed successfully!"
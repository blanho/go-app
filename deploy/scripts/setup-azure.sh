set -e

RESOURCE_GROUP="azure-go-app-prod-rg"
LOCATION="eastus"
APP_NAME="azure-go-app"
ENVIRONMENT="prod"

echo "Logging in to Azure..."
az login --service-principal --username $SP_CLIENT_ID --password $SP_CLIENT_SECRET --tenant $SP_TENANT_ID


echo "Setting subscription..."
az account set --subscription $SUBSCRIPTION_ID


echo "Creating resource group if it doesn't exist..."
az group create --name $RESOURCE_GROUP --location $LOCATION


echo "Initializing Terraform..."
cd ./deploy/terraform


STORAGE_ACCOUNT="tfstate$RANDOM"
CONTAINER_NAME="tfstate"

echo "Creating storage account for Terraform state..."
az storage account create --resource-group $RESOURCE_GROUP --name $STORAGE_ACCOUNT --sku Standard_LRS --encryption-services blob

echo "Creating storage container..."
az storage container create --name $CONTAINER_NAME --account-name $STORAGE_ACCOUNT

echo "Getting storage account key..."
ACCOUNT_KEY=$(az storage account keys list --resource-group $RESOURCE_GROUP --account-name $STORAGE_ACCOUNT --query '[0].value' -o tsv)

echo "Configuring Terraform backend..."
cat > backend.tfvars <<EOF
resource_group_name  = "$RESOURCE_GROUP"
storage_account_name = "$STORAGE_ACCOUNT"
container_name       = "$CONTAINER_NAME"
key                  = "${APP_NAME}-${ENVIRONMENT}.terraform.tfstate"
EOF

terraform init -backend-config=backend.tfvars

echo "Applying Terraform configuration..."
terraform apply -var="app_name=$APP_NAME" -var="environment=$ENVIRONMENT" -var="location=$LOCATION" -auto-approve

echo "Azure infrastructure setup completed successfully!"
## Prerequisites

- Go 1.22+
- Docker
- Azure CLI
- Terraform
- kubectl
- An Azure subscription

## Deployment

### Infrastructure Setup

1. Deploy Azure infrastructure
cd deploy/scripts
./setup-azure.sh

2. Deploy the application
./deploy.sh

### Manual Deployment

1. Apply Terraform configuration
cd deploy/terraform
terraform init
terraform apply

2. Deploy to Kubernetes
cd deploy/kubernetes
kubectl apply -f .

### Key Components

- **Application Insights**: Telemetry and monitoring
- **Azure Key Vault**: Secure secret management
- **Azure Kubernetes Service**: Scalable container orchestration
- **Azure SQL Database**: Persistent data storage
- **Azure Cache for Redis**: Caching layer
- **Azure Service Bus**: Messaging and event-driven communication

## Configuration

Configuration is managed through environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| PORT | Server port | 8080 |
| ENVIRONMENT | Deployment environment | development |
| LOG_LEVEL | Logging level | info |
| DATABASE_URL | SQL Server connection string | |
| APPLICATION_INSIGHTS_KEY | App Insights instrumentation key | |
| KEY_VAULT_NAME | Azure Key Vault name | |
| MAX_CONNECTIONS | Maximum DB connections | 100 |
| SHUTDOWN_TIMEOUT | Graceful shutdown timeout (seconds) | 10 |

## License
blanho

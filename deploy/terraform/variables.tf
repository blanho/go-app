# deploy/terraform/variables.tf
variable "app_name" {
  description = "The name of the application"
  type        = string
  default     = "azure-go-app"
}

variable "location" {
  description = "The Azure region to deploy resources to"
  type        = string
  default     = "East US"
}

variable "environment" {
  description = "The environment (dev, test, prod)"
  type        = string
  default     = "prod"
}

variable "kubernetes_version" {
  description = "The version of Kubernetes"
  type        = string
  default     = "1.27.3"
}

variable "node_count" {
  description = "The initial number of nodes in the AKS cluster"
  type        = number
  default     = 3
}

variable "min_node_count" {
  description = "The minimum number of nodes in the AKS cluster"
  type        = number
  default     = 3
}

variable "max_node_count" {
  description = "The maximum number of nodes in the AKS cluster"
  type        = number
  default     = 10
}

variable "node_size" {
  description = "The size of the AKS nodes"
  type        = string
  default     = "Standard_D2s_v3"
}
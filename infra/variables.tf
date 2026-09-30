variable "aws_region" {
  description = "AWS region for all infrastructure."
  type        = string
  default     = "us-east-1"
}

variable "project_name" {
  description = "Name used to identify project resources."
  type        = string
  default     = "ecom-app"
}

variable "environment" {
  description = "Deployment environment name."
  type        = string
  default     = "dev"
}

variable "kubernetes_version" {
  description = "EKS Kubernetes version."
  type        = string
  default     = "1.31"
}

variable "node_instance_types" {
  description = "EC2 instance types for the managed EKS node group."
  type        = list(string)
  default     = ["t3.medium"]
}

variable "node_desired_size" {
  description = "Initial number of worker nodes."
  type        = number
  default     = 2
}

variable "node_min_size" {
  description = "Minimum number of worker nodes."
  type        = number
  default     = 2
}

variable "node_max_size" {
  description = "Maximum number of worker nodes."
  type        = number
  default     = 4
}

variable "cluster_endpoint_public_access_cidrs" {
  description = "CIDR blocks allowed to reach the public EKS API endpoint. Defaults to open for demo convenience - restrict to your IP before anything beyond a short-lived lab."
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "ecr_force_delete" {
  description = "Allow `terraform destroy` to delete ECR repositories even if they still have images. Pass at runtime: -var=\"ecr_force_delete=true\"."
  type        = bool
  default     = false
}

variable "services" {
  description = "Application services, each gets its own ECR repository (frontend/middleware/backend 3-tier app)."
  type        = list(string)
  default     = ["frontend", "middleware", "backend"]
}

variable "github_repo" {
  description = "GitHub repo allowed to assume the CI role via OIDC, in `org/repo` form."
  type        = string
  default     = "RanjitThakur678/argoCD-k8s-demo"
}

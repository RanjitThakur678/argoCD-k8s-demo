output "cluster_name" {
  description = "Name of the EKS cluster."
  value       = module.eks.cluster_name
}

output "aws_region" {
  description = "AWS region used by this deployment."
  value       = var.aws_region
}

output "cluster_endpoint" {
  description = "EKS API endpoint."
  value       = module.eks.cluster_endpoint
}

output "ecr_repository_urls" {
  description = "Map of service name -> ECR repository URL."
  value       = { for k, r in aws_ecr_repository.app : k => r.repository_url }
}

output "configure_kubectl" {
  description = "Command to configure kubectl for this cluster."
  value       = "aws eks update-kubeconfig --region ${var.aws_region} --name ${module.eks.cluster_name}"
}

output "lb_controller_role_arn" {
  description = "IAM role ARN for the AWS Load Balancer Controller's service account (IRSA)."
  value       = aws_iam_role.lb_controller.arn
}

output "argocd_image_updater_role_arn" {
  description = "IAM role ARN for Argo CD Image Updater's service account (IRSA) - read-only ECR access."
  value       = aws_iam_role.argocd_image_updater.arn
}

output "github_actions_role_arn" {
  description = "IAM role ARN GitHub Actions assumes via OIDC to push images to ECR."
  value       = aws_iam_role.github_actions.arn
}

output "vpc_id" {
  description = "VPC ID, needed by the AWS Load Balancer Controller Helm install."
  value       = module.vpc.vpc_id
}

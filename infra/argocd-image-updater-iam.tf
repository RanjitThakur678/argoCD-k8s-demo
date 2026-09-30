# IRSA role for Argo CD Image Updater (namespace `argocd`, service account `argocd-image-updater` -
# the Helm chart's default name). Read-only ECR access across all 3 service repos so it can detect
# new image tags pushed by CI and write the winning tag back to Git for Argo CD to sync.

data "aws_iam_policy_document" "argocd_image_updater_assume" {
  statement {
    effect  = "Allow"
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.eks.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(module.eks.cluster_oidc_issuer_url, "https://", "")}:sub"
      values   = ["system:serviceaccount:argocd:argocd-image-updater"]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(module.eks.cluster_oidc_issuer_url, "https://", "")}:aud"
      values   = ["sts.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "argocd_image_updater" {
  name               = "${local.name}-argocd-image-updater"
  assume_role_policy = data.aws_iam_policy_document.argocd_image_updater_assume.json
}

data "aws_iam_policy_document" "argocd_image_updater" {
  statement {
    effect    = "Allow"
    actions   = ["ecr:GetAuthorizationToken"]
    resources = ["*"]
  }

  statement {
    effect = "Allow"
    actions = [
      "ecr:DescribeRepositories",
      "ecr:ListImages",
      "ecr:DescribeImages",
      "ecr:BatchGetImage",
      "ecr:GetDownloadUrlForLayer",
    ]
    resources = [for r in aws_ecr_repository.app : r.arn]
  }
}

resource "aws_iam_role_policy" "argocd_image_updater" {
  name   = "${local.name}-argocd-image-updater-policy"
  role   = aws_iam_role.argocd_image_updater.id
  policy = data.aws_iam_policy_document.argocd_image_updater.json
}

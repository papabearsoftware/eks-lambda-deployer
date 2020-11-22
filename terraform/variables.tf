variable "github_org" {
  type        = string
  description = "GitHub account name"
}

variable "github_repo" {
  type = string
}

variable "github_branch" {
  type    = string
  default = "main"
}

variable "github_token" {
  type        = string
  description = "GitHub personal access token"
}

variable "region" {
  type = string
}
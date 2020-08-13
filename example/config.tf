provider "aws" {
  region = "ap-northeast-1"
}

terraform {
  backend "local" {
    path = ".terraform/local.tfstate"
  }
}

variable "vpc_id" {
  type = string
}

variable "main_az_a" {
  type = string
}

variable "main_az_c" {
  type = string
}

data "aws_caller_identity" "current" {}

data "aws_ecs_cluster" "default" {
    cluster_name = "default"
}

data "aws_ecs_task_definition" "hello-world" {
  task_definition = "hello-world"
}

data "aws_vpc" "main" {
  id = var.vpc_id
}

data "aws_subnet_ids" "public" {
  vpc_id = var.vpc_id

  filter {
    name   = "tag:Name"
    values = ["main-az-*"]
  }
}

data "aws_security_groups" "internal" {
  filter {
    name   = "tag:Name"
    values = ["main-internal"]
  }
}

data "aws_iam_role" "ecsEventsRole" {
  name = "ecsEventsRole"
}

data "aws_iam_role" "ecsTaskExecutionRole" {
  name = "ecsTaskExecutionRole"
}

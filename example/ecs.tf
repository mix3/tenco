resource "aws_ecs_task_definition" "hello-world" {
  family                   = "hello-world"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  cpu                      = 256
  memory                   = 512
  execution_role_arn       = data.aws_iam_role.ecsTaskExecutionRole.arn
  container_definitions    = <<TASK_DEFINITION
[
    {
        "name": "hello-world",
        "image": "alpine:latest",
        "command": ["echo", "hello-world"],
        "logConfiguration": {
            "logDriver": "awslogs",
            "options": {
                "awslogs-group": "/ecs/hello-world",
                "awslogs-region": "ap-northeast-1",
                "awslogs-stream-prefix": "ecs"
            }
        }
   }
]
TASK_DEFINITION
}

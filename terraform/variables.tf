# input variable can be changed without changeing main.terraform 

variable "container_image" {
    description = "Docker image name and tag to deploy"
    type = string
    default = "go-webapp:latest" #default value
}

variable "node_port"{
    description = "The nodePort to expose the service externally (30000 - 32767)"
    type= number

    # validation -> to ensure port is in a valid range
    validation {
        condition = var.node_port >= 30000 && var.node_port <= 32767
        error_message = "NodePort must be between 30000 and 32767"
    }
    default = 30000
}

variable "replicas" {
    description = "Number of pod replicas to run"
    type = number
    default = 2
}

variable "go_env" {
    description = "Environemnt for Go application"
    type = string
    default = "development"
}


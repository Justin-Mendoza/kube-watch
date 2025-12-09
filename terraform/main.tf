# Terraform BLOCK 0 terraform config

terraform {
    # min terraform version
    required_version = ">= 1.0"

    # provider requiremnts - what providers we are asking for
    required_providers {
        kubernetes = {
        # from official hashicorp registry
        source = "hashicorp/kubernetes"
        # versuin constraints ~> means any version 2.x.x
        version = "~> 2.23"
        }   
    }
}

# PROVIDER CONFIGUARTION how terraform connects to kubernetes
provider "kubernetes" {
    # uses local kube config (made by docker/minikube)
    config_path = "~/.kube/config"
}

# RESOURCES

# NAMESPACEAS -grouping of resources
resource "kubernetes_namespace" "app"{
    metadata{
        name = "local-go-app" #namespace name
    }
}

resource "kubernetes_deployment" "webapp"{
    # metadata abt deployment
    metadata {
        name = "go-webapp"
        namespace = kubernetes_namespace.app.metadata[0].name
        # namespace reference

    # labels - used for selecting resources
        labels = {
            app = "go-webapp"
            project = "k8s-go-infra"
        }
    }

    spec {
        replicas = var.replicas
    # selector - tells deploymeny what to manage
    selector {
        match_labels = {
            app = "go-webapp"
        }
    }

        template{
            metadata{
                labels = {
                    app = "go-webapp"
                }
            }

            spec {
                container {
                    # my docker image
                    image = var.container_image
                    name = "webapp-container"

                    port {
                        container_port = 8080
                    }

                    resources{
                        limits = {
                            cpu = "500m"
                            memory = "256Mi"
                        }
                        requests = {
                            
                        cpu = "250m"
                        memory = "128Mi"
                        }
                    }

                    # health check for monitoring
                    liveness_probe {
                        http_get {
                            path = "/"
                            port = 8080
                        }
                    # start 10 seconds adfter container start & check health every 10 seconds
                    initial_delay_seconds = 10
                    period_seconds = 10
                    }

                }
            }

        }

    }

}

resource "kubernetes_service" "webapp"{
    metadata {
        name = "go-webapp-service"
        namespace = kubernetes_namespace.app.metadata[0].name
    }

    spec{
        # selector - whichpods this service routes to
        selector = {
            app = kubernetes_deployment.webapp.metadata[0].labels.app
        }
        # port config
        port {
            port = 8080
            target_port = 8080
            node_port = var.node_port
        }

        type = "NodePort" 
        # service accessible to outside clustesr
    }
}
# values that terraform can displau after applying

output "namespace" { 
    description = "Namespace where application is deployed"
    value = kubernetes_namespace.app.metadata[0].name
}

output "deployment_name" { 
    description = "Name of Kubernetes deploymeny"
    value = kubernetes_deployment.webapp.metadata[0].name
}

output "service_name"{
    description = "Name of Kubernetes Service"
    value = kubernetes_service.webapp.metadata[0].name
}

output "node_port" {
    description = "NodePort assigned to the service"
    value = kubernetes_service.webapp.spec[0].port[0].node_port

    # sensitive = false => safe to show on console
    sensitive = false
}

output "service_url" {
    description = "URL to access service"
    value = "http://localhost:${kubernetes_service.webapp.spec[0].port[0].node_port}"
}

output "kubectl_commands"{
    description = "Useful kubectl commands for debugging"
    value = <<EOT
To check pods:    kubectl get pods -n ${kubernetes_namespace.app.metadata[0].name}
To check service: kubectl get svc -n ${kubernetes_namespace.app.metadata[0].name}
To view logs:     kubectl logs -n ${kubernetes_namespace.app.metadata[0].name} -l app=go-webapp
To delete all:    kubectl delete namespace ${kubernetes_namespace.app.metadata[0].name}
EOT
}
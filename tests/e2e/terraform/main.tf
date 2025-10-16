terraform {
  required_providers {
    spot = {
      source  = "rackerlabs/spot"
      version = ">= 0.2.0"
    }
  }
}

provider "spot" {
  token = var.rackspace_token
}

resource "spot_cluster" "e2e" {
  name           = "grapevine-e2e"
  region         = "us-central1"
  version        = "1.24"
  node_pools {
    node_count = 1
    name       = "worker-pool"
    size       = "spot-v1-small-2"
  }
}

output "kubeconfig" {
  value     = spot_cluster.e2e.kubeconfig
  sensitive = true
}
terraform {
  required_providers {
    scaleway = {
      source = "scaleway/scaleway"
    }
  }
  required_version = ">= 0.13"
}

resource "scaleway_k8s_cluster" "cluster" {
  name    = "waypoint-cluster"
  version = "1.26.0"
  cni     = "cilium"
  delete_additional_resources = true
}

resource "scaleway_k8s_pool" "pool" {
	cluster_id = scaleway_k8s_cluster.cluster.id
	name = "main-pool"
	node_type = "DEV1-M"
	size = 1
}

resource "scaleway_container_namespace" "namespace" {
	region = "nl-ams"
}

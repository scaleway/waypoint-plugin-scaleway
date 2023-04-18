This demo will guide you testing waypoint on scaleway by deploying Waypoint on a kubernetes cluster and deploying a demo application onto Scaleway Serverless Container using Waypoint

## Used resources

- Scaleway Kubernetes Kapsule (free control-plane, the only node used is paid)
- Scaleway Container Registry (Pushed image storage)
- Scaleway Serverless Containers (Deployed container)

## Requirements

- `docker`
- `jq`
- `terraform` [download](https://developer.hashicorp.com/terraform/downloads)
- `waypoint` [download](https://developer.hashicorp.com/waypoint/downloads)
- `waypoint-scaleway-plugin` [repo](https://github.com/scaleway/waypoint-plugin-scaleway)

## Usage

- Setup scaleway config or a complete environment with all required variables ([doc](../docs/scw-config.md))
- `docker login rg.fr-par.scw.cloud -u nologin -p ${SCW_SECRET_KEY}` Login to scaleway registry if not already logged in
- `terraform init` Init terraform project
- `terraform apply` Create your kubernetes cluster and container namespace
- `./setup.sh` Install waypoint server and complete the `waypoint.hcl` with your container namespace informations
- `waypoint ui -authenticate` Connect to waypoint ui
- (Optional) If on a server, open Web UI Address given by setup script and authenticate manually
- `waypoint init` Init waypoint project
- `waypoint up` Build and deploy your app on scaleway's containers

## Cleanup

- `terraform destroy`

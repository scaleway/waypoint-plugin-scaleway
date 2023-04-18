set -e

KUBECONFIG_FILE=`terraform show -json | jq '.values.root_module.resources[] | select(.address == "scaleway_k8s_cluster.cluster") | .values.kubeconfig[0].config_file' -r`

echo "${KUBECONFIG_FILE}" > ./kubeconfig

echo "Installing waypoint server..."

KUBECONFIG=./kubeconfig waypoint install -platform=kubernetes -accept-tos

CONTAINER_NAMESPACE_ID_REGIONAL=`terraform show -json | jq '.values.root_module.resources[] | select(.address == "scaleway_container_namespace.namespace") | .values.id' -r`
CONTAINER_NAMESPACE_ID="${CONTAINER_NAMESPACE_ID_REGIONAL#*/}"

REGISTRY_ENDPOINT=`terraform show -json | jq '.values.root_module.resources[] | select(.address == "scaleway_container_namespace.namespace") | .values.registry_endpoint' -r`

echo "Completing waypoint.hcl..."

sed -i.bak "s/{container-namespace}/${CONTAINER_NAMESPACE_ID}/g" waypoint.hcl
sed -i.bak "s#{registry-endpoint}#${REGISTRY_ENDPOINT}#g" waypoint.hcl

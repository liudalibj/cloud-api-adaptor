#!/usr/bin/env bash

set -o errexit
set -o pipefail
set -o nounset

script_dir=$(dirname "$(readlink -f "$0")")

registry="${registry:-quay.io/liudalibj}"
plugins_name="cloud-provider-plugins"
caa_name="cloud-api-adaptor"
peerpod_ctrl_name="peerpod-ctrl"

commit=$(git rev-parse HEAD)
[[ -n "$(git status --porcelain --untracked-files=no)" ]] && commit+='-dirty'

dev_tags=${DEV_TAGS:-"${commit}"}
plugins_tags="latest,dev-${dev_tags}"
extend_tags="dev-plugins-${dev_tags}"

supported_arches=${ARCHES:-"linux/amd64,linux/s390x"}

# Get a list of comma-separated tags (e.g. latest,dev-5d0da3dc9764), return
# the tag string (e.g "-t ${registry}/${name}:latest -t ${registry}/${name}:dev-5d0da3dc9764")
#
function get_tag_string() {
	local name="$1"
	local tags="$2"
	local tag_string=""

	for tag in ${tags/,/ };do
		tag_string+=" -t ${registry}/${name}:${tag}"
	done

	echo "$tag_string"
}

function build_cloud_provider_plugins_image() {
	pushd "${script_dir}/../"
	local tag_string
	tag_string="$(get_tag_string "$plugins_name" "$plugins_tags")"

	docker buildx build --platform "${supported_arches}" \
		--build-arg COMMIT="${commit}" \
		-f cloud-providers-plugins/Dockerfile.plugins \
		${tag_string} \
		--push \
		.
	popd
}

function build_caa_with_plugins_image() {
	pushd "${script_dir}/../"
	local tag_string
	tag_string="$(get_tag_string "$caa_name" "$extend_tags")"

	docker buildx build --platform "${supported_arches}" \
		-f cloud-providers-plugins/Dockerfile.caa \
		${tag_string} \
		--push \
		.
	popd
}

function build_peerpod_ctrl_with_plugins_image() {
	pushd "${script_dir}/../"
	local tag_string
	tag_string="$(get_tag_string "$peerpod_ctrl_name" "$extend_tags")"

	docker buildx build --platform "${supported_arches}" \
		-f cloud-providers-plugins/Dockerfile.peerpod_ctrl \
		${tag_string} \
		--push \
		.
	popd
}

# Get the options
while getopts ":pcl" option; do
    case $option in
        p) # build cloud provider plugins
            build_cloud_provider_plugins_image
            exit;;
        c) # caa + cloud provider plugins
            build_caa_with_plugins_image
            exit;;
        l) # peerpod-ctrl + cloud provider plugins
            build_peerpod_ctrl_with_plugins_image
            exit;;
        \?) # Invalid option
            echo "Error: Invalid option"
   esac
done


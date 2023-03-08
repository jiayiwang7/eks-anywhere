#!/usr/bin/env bash

set -e
set -o pipefail
set -o nounset
set -x

buckect_dst=$1
bucket_dst_path=${buckect_dst#"s3://"}
bucket_dst_name=${bucket_dst_path%%/*}
bucket_dst_dir=${bucket_dst_path#"$bucket_dst_name"}
base_public_files_url="https://$bucket_dst_name.s3.amazonaws.com$bucket_dst_dir"
releases_s3_url=$buckect_dst/releases.yaml
releases_url=$base_public_files_url/releases.yaml

og_dev_bundle_url="https://dev-release-prod-pdx.s3.us-west-2.amazonaws.com/bundle-release.yaml"
version=$(git -C ./ describe --tag)

latest_bundle_path="latest_bundle.yaml"
curl --silent -L $og_dev_bundle_url > $latest_bundle_path

bundle_s3_url=$buckect_dst/bundles/$version/bundle.yaml
bundles_public_url=$base_public_files_url/bundles/$version/bundle.yaml
aws s3 cp --acl public-read $latest_bundle_path $bundle_s3_url
rm $latest_bundle_path

my_eksa_bin_dir=bin/my-eks-a
my_eksa_bin=$my_eksa_bin_dir/eksctl-anywhere
my_eksa_tar="eksctl-anywhere-$version-linux-amd64.tar.gz"
my_eksa_tar_s3_url=$buckect_dst/releases/$version/$my_eksa_tar
my_eksa_tar_url=$base_public_files_url/releases/$version/$my_eksa_tar

make eks-a-binary GIT_VERSION=$version RELEASE_MANIFEST_URL=$releases_url OUTPUT_FILE=$my_eksa_bin
tar -czvf $my_eksa_tar -C $my_eksa_bin_dir eksctl-anywhere 

aws s3 cp --acl public-read $my_eksa_tar $my_eksa_tar_s3_url
rm $my_eksa_tar 

current_date=$(date "+%F %T.%9N %z")

new_release_declaration='{
  "bundleManifestUrl": "'$bundles_public_url'",
  "date": "'$current_date'",
  "eksABinary": {
    "linux": {
      "arch": [
        "amd64"
      ],
      "description": "EKS Anywhere CLI",
      "name": "eksctl-anywhere-linux",
      "os": "linux",
      "uri": "'$my_eksa_tar_url'"
    }
  },
  "gitTag": "'$version'",
  "number": 1,
  "version": "'$version'"
}'

# download current releases file and update it
releases_tmp=my-releases.yaml
curl --silent -L $releases_url > $releases_tmp

yq ".spec.releases += $new_release_declaration" -i $releases_tmp
yq '.spec.latestVersion = "'$version'"' -i $releases_tmp

aws s3 cp --acl public-read $releases_tmp $releases_s3_url
rm $releases_tmp

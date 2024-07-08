#!/bin/bash
get_yaml_value() {
	local key_path=$1
	local yaml_file="../SAVE_FILES/config.yaml"
	local value=$(python3 -c "
import sys, yaml
with open('$yaml_file', 'r') as file:
	config = yaml.safe_load(file)
keys = '$key_path'.split('.')
value = config
try:
	for key in keys:
		value = value[key]
	print(value)
except KeyError:
	print('Key not found')
")
	echo "$value"
}
GO_VERSION=$(get_yaml_value "go.version")
ARCH=$(dpkg --print-architecture)
case $ARCH in
	armhf)
		GOLANG_ARCH_NAME="linux-armv6l"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	arm64)
		GOLANG_ARCH_NAME="linux-arm64"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	amd64)
		GOLANG_ARCH_NAME="linux-amd64"
		ARCHIVE_NAME=$(echo "go$GO_VERSION.$GOLANG_ARCH_NAME.tar.gz")
		ARCHIVE_URL=$(echo "https://golang.org/dl/$ARCHIVE_NAME")
		;;
	*)
		echo "unknown arch"
		#GOLANG_ARCH_NAME="windows-amd64.zip"
		;;
esac
echo "Downloading $ARCHIVE_URL"
wget $ARCHIVE_URL --progress=bar -O go.tar.gz
echo "Extracting To : /usr/local/go"
mkdir -p ~/go
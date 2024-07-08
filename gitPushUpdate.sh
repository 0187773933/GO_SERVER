#!/bin/bash
function is_int() { return $(test "$@" -eq "$@" > /dev/null 2>&1); }
ssh-add -D
git init
git config --global --unset user.name
git config --global --unset user.email
git config user.name "0187773933"
git config user.email "collincerbus@student.olympic.edu"
ssh-add -k /Users/morpheous/.ssh/githubWinStitch

get_yaml_value() {
	local key_path=$1
	local yaml_file="./SAVE_FILES/config.yaml"
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

GIT_SSH_URL=$(get_yaml_value "git.ssh_url")
echo "GIT_SSH_URL: $GIT_SSH_URL"

LastCommit=$(git log -1 --pretty="%B" | xargs)
# https://stackoverflow.com/a/3626205
if $(is_int "${LastCommit}");
	 then
	 NextCommitNumber=$((LastCommit+1))
else
	echo "Not an integer Resetting"
	NextCommitNumber=1
fi
git add .
git tag -l | xargs git tag -d
if [ -n "$1" ]; then
	git commit -m "$1"
	git tag v1.1.$1
else
	git commit -m "$NextCommitNumber"
	git tag v1.1.$NextCommitNumber
fi
git remote add origin "$GIT_SSH_URL"

git push origin --tags
git push origin master
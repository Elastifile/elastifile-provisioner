#!/bin/bash
set -e

usage() {
    echo -e "Usage:"
    echo -e "\t-u docker user"
    echo -e "\t-p docker password"
    echo -e "\t-e docker email"
    echo -e "\t-b git branch"
    echo -e "\t-h print this help message"
}

while getopts ":hu:p:e:b:" opt; do
  case $opt in
    u)
      docker_user=$OPTARG
      ;;
    p)
      docker_password=$OPTARG
      ;;
    e)
      docker_email=$OPTARG
      ;;
    b)
      branch=$OPTARG
      ;;
    h)
      usage
      exit 1
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      ;;
  esac
done

if [ ! "$docker_user" ] || [ ! "$docker_password" ] || [ ! "$docker_email" ] || [ ! "$branch" ]
then
    echo "Missing args (must have: user + password + email + valid git branch)"
    usage
    exit 1
fi

GOROOT=/tmp/go1.9.2
GOPATH=/tmp/gopath
rm -rf $GOPATH
mkdir -p $GOROOT $GOPATH

install_file="go1.9.2.linux-amd64.tar.gz"
install_url="https://redirector.gvt1.com/edgedl/go/$install_file"

if ! [ "$(command -v $GOROOT/bin/go)" ]; then
  echo "## INSTALL GO ##"
  wget $install_url
  tar -C $GOROOT -xzf $install_file --strip 1
fi

export GOPATH=$GOPATH
export GOROOT=$GOROOT
export PATH=$GOROOT/bin:$PATH

provisioner_url="git@github.com/Elastifile/elastifile-provisioner"
provisioner_in_gopath=$GOPATH/src/$provisioner_url

go get $provisioner_url
cd $provisioner_in_gopath

git checkout $branch
if [ $branch != "master" ]; then
    git tag "development"
fi


docker login -u $docker_user -p $docker_password -e $docker_email
make clean
make push

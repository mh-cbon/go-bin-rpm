#!/bin/sh -e

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
#
# to use it
# curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/create-pkg.sh \
# | GH=mh-cbon/gh-api-cli sh -xe

if [ "${GH}" = "mh-cbon/go-bin-rpm" ]; then
  git pull origin master
  git checkout -b master
fi

if ["${GH_TOKEN}" = ""]; then
  echo "GH_TOKEN is not properly set. Check your travis file."
  exit 1
fi

rm -fr docker.sh
set +x
cat <<EOT > docker.sh

if type "dnf" > /dev/null; then
  sudo dnf install wget curl git -y
else
  sudo yum install wget curl git -y
fi

export GH_TOKEN="${GH_TOKEN}"
export GH="${GH}"

export GOINSTALL="/go"
export GOROOT=\${GOINSTALL}/go/
export PATH=\$PATH:\$GOROOT/bin

# install go, specific to vagrant
if type "wget" > /dev/null; then
  wget $getgo | sh -xe
fi
if type "curl" > /dev/null; then
  curl -L $getgo | sh -xe
fi

echo "\$PATH"
go version
go env

export GOPATH=/gopath/
export PATH=\$PATH:/gopath/bin


yes | go get -u github.com/mh-cbon/go-bin-rpm/go-bin-rpm-utils
go-bin-rpm-utils create-packages -push -repo=$GH
EOT
set -x

docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./docker.sh ${TRAVIS_TAG}"
sudo chown travis:travis ${REPO}-*.rpm

rm -fr docker.sh

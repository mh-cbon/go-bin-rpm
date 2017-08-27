#!/bin/sh -e
set -x

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
# it creates rpm file over a docker.
#
# to use it
# curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/create-pkg.sh \
# | GH=mh-cbon/gh-api-cli sh -xe

if [ "${GH}" = "mh-cbon/go-bin-rpm" ]; then
  git reset HEAD --hard
  git pull origin --force || echo "ok"
  git checkout -b master || echo "ok"
  git checkout master || echo "ok"
  curl https://glide.sh/get | sh
  glide install
fi

getgo="https://raw.githubusercontent.com/mh-cbon/latest/master/get-go.sh?d=`date +%F_%T`"

rm -fr docker.sh
TRAVIS_BUILD_DIR="/gopath/src/github.com/${GH}"
cat <<EOT > docker.sh
set -x
set -e
if type "dnf" > /dev/null; then
  if type "sudo" > /dev/null; then
    sudo dnf install wget curl git -y
  else
    dnf install wget curl git -y
  fi
else
  if type "sudo" > /dev/null; then
    sudo yum install wget curl git -y
  else
    yum install wget curl git -y
  fi
fi

set +x
export GH_TOKEN="${GH_TOKEN}"
set -x

export TRAVIS_TAG="${TRAVIS_TAG}"
export TRAVIS_BUILD_DIR="${TRAVIS_BUILD_DIR}"
export GH="${GH}"
export EMAIL="${EMAIL}"
export MYEMAIL="${MYEMAIL}"
export TRAVIS="${TRAVIS}"
export CI="${CI}"

export GOINSTALL="/go"
export GOROOT=\${GOINSTALL}/go/
export PATH=\$PATH:\$GOROOT/bin

echo "GH \$GH"
echo "getgo $getgo"

# install go, specific to vagrant
if type "wget" > /dev/null; then
  wget --quiet -O - $getgo | sh -xe
fi
if type "curl" > /dev/null; then
  curl -s -L $getgo | sh -xe
fi

echo "PATH \$PATH"
go version
go env

export GOPATH=/gopath/
export PATH=\$PATH:/\$GOPATH/bin

go get -u github.com/mh-cbon/go-bin-rpm/go-bin-rpm-utils

set +x
if [ -z "$GH_TOKEN" ]; then
  set -x
  go-bin-rpm-utils create-packages -keep -repo=$GH
else
  set -x
  go-bin-rpm-utils create-packages -push -keep -repo=$GH
fi
ls -al .
EOT
set -x
docker run -v $PWD/:${TRAVIS_BUILD_DIR} fedora /bin/sh -c "cd ${TRAVIS_BUILD_DIR} && sh ./docker.sh"
sudo chown travis:travis ./*-*.rpm

ls -al .

rm -fr docker.sh

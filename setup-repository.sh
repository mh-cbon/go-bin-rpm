#!/bin/sh -e
set -x

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
# it creates an rpm repo
#
# to use it
# curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/setup-repository.sh \
# | GH=mh-cbon/gh-api-cli EMAIL=mh-cbon@users.noreply.github.com sh -xe

# GH=$1
# EMAIL=$2

if ["${GH}" = ""]; then
  echo "GH is not properly set. Check your travis file."
  exit 1
fi

if [ "${GH}" = "mh-cbon/go-bin-rpm" ]; then
  git reset HEAD --hard
  git pull origin --force || echo "ok"
  git reset HEAD --hard
  ls -al
  git status
  git checkout -b master || echo "ok"
  git checkout master || echo "ok"
  curl https://glide.sh/get | sh
  glide install
fi


getgo="https://raw.githubusercontent.com/mh-cbon/latest/master/get-go.sh?d=`date +%F_%T`"
TRAVIS_BUILD_DIR="/gopath/src/github.com/${GH}"

rm -fr docker.sh
cat <<EOT > docker.sh

set -x
if type "dnf" > /dev/null; then
  if type "sudo" > /dev/null; then
    sudo dnf install wget curl git -y --quiet
  else
    dnf install wget curl git -y --quiet
  fi
else
  if type "sudo" > /dev/null; then
    sudo yum install wget curl git -y --quiet
  else
    yum install wget curl git -y --quiet
  fi
fi

set +x
export GH_TOKEN="${GH_TOKEN}"
set -x

export TRAVIS_BUILD_DIR="${TRAVIS_BUILD_DIR}"
export TRAVIS_TAG="${TRAVIS_TAG}"
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
export PATH=\$PATH:/gopath/bin

go get -u github.com/mh-cbon/go-bin-rpm/go-bin-rpm-utils
go-bin-rpm-utils setup-repository -out="\`pwd\`/rpm" -push -repo=$GH
EOT

set -x
docker run -v $PWD/:${TRAVIS_BUILD_DIR} fedora /bin/sh -c "cd ${TRAVIS_BUILD_DIR} && sh ./docker.sh"

rm -fr docker.sh

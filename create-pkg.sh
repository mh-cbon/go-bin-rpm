#!/bin/sh -e
set -x

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
  curl https://glide.sh/get | sh
  glide install
fi

getgo="https://raw.githubusercontent.com/mh-cbon/latest/master/get-go.sh?d=`date +%F_%T`"

rm -fr docker.sh

cat <<EOT > docker.sh
set -x
if type "dnf" > /dev/null; then
  if type "sudo" > /dev/null; then
    sudo dnf install wget curl git -y
  else
    dnf install wget curl git -y
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
  wget -q $getgo | sh -xe
fi
if type "curl" > /dev/null; then
  curl -s -L $getgo | sh -xe
fi

echo "PATH \$PATH"
go version
go env

export GOPATH=/gopath/
export PATH=\$PATH:/\$GOPATH/bin

mkdir -p \$GOPATH/src/github.com/${GH}/
cp -R /docker/* \$GOPATH/src/github.com/${GH}/*

ls -al \$GOPATH/src/github.com/${GH}/*

go get -u github.com/mh-cbon/go-bin-rpm/go-bin-rpm-utils
go-bin-rpm-utils create-packages -push -repo=$GH
EOT
set -x

docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./docker.sh"
# sudo chown travis:travis *-*.rpm

rm -fr docker.sh

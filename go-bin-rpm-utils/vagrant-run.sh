#!/bin/sh -e

set -e
set -x

sudo yum install git -y

export GOINSTALL="/go"
export GOROOT=${GOINSTALL}/go/
export PATH=$PATH:$GOROOT/bin

getgo="https://raw.githubusercontent.com/mh-cbon/latest/master/get-go.sh?d=`date +%F_%T`"
# install go, specific to vagrant
if type "wget" > /dev/null; then
  wget $getgo -O - | sh -xe
fi
if type "curl" > /dev/null; then
  curl -s -L $getgo | sh -xe
fi

echo "$PATH"
go version
go env

export GOPATH=/gopath/
export PATH=$PATH:/gopath/bin

sudo chown -R vagrant:vagrant -R $GOPATH
mkdir -p ${GOPATH}/bin

[ -d "$GOPATH" ] || echo "$GOPATH does not exists, do you run vagrant ?"
[ -d "$GOPATH" ] || exit 1;


set +x
# everything here will be replicated into the CI build file (.travis.yml)
export GH_TOKEN="$GH_TOKEN"
NAME="go-bin-rpm"
export REPO="mh-cbon/$NAME"
export EMAIL="mh-cbon@users.noreply.github.com"

# set env specific to vagrant
export VERSION="LAST"
export BUILD_DIR="$GOPATH/src/github.com/$REPO/pkg-build"

# set env specific to travis
# export TRAVIS_TAG="0.0.1-beta999"
# export TRAVIS_BUILD_DIR="$GOPATH/src/github.com/$REPO/pkg-build"

set -x
# setup glide
if type "glide" > /dev/null; then
  echo "glide already installed"
  glide -v
else
  curl https://glide.sh/get | sh
fi

cd $GOPATH/src/github.com/$REPO
git checkout master
git reset HEAD --hard
[ -d "$GOPATH/src/github.com/$REPO/vendor" ] || glide install
go install

# build the binaries
BINBUILD_DIR="$GOPATH/src/github.com/$REPO/build"
rm -fr "$BINBUILD_DIR"
mkdir -p "$BINBUILD_DIR/{386,amd64}"

PKGBUILD_DIR="$GOPATH/src/github.com/$REPO/rpm"

# build the packages
set +x
echo ""
echo "# =================================================="
echo "# =================================================="
echo ""
set -x
f="-X main.VERSION=${VERSION}"
GOOS=linux GOARCH=386 go build --ldflags "$f" -o "$BINBUILD_DIR/386/go-bin-rpm" $k
GOOS=linux GOARCH=amd64 go build --ldflags "$f" -o "$BINBUILD_DIR/amd64/go-bin-rpm" $k
go run /vagrant/*go create-packages -push -repo=$REPO

set +x
echo ""
echo "# =================================================="
echo "# =================================================="
echo ""
set -x
go run /vagrant/*go setup-repository -out="${PKGBUILD_DIR}" -push -repo=$REPO

set +x
echo ""
echo "# =================================================="
echo "# =================================================="
echo "      All Done!"
set -x

#

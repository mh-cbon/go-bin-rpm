#!/bin/sh -e

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
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
  git pull origin master
  git checkout -b master
  curl https://glide.sh/get | sh
  glide install
fi


getgo="https://raw.githubusercontent.com/mh-cbon/latest/master/get-go.sh?d=`date +%F_%T`"

rm -fr docker.sh
set +x
cat <<EOT > docker.sh
export GH_TOKEN="${GH_TOKEN}"
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
  wget $getgo | sh -xe
fi
if type "curl" > /dev/null; then
  curl -L $getgo | sh -xe
fi

echo "PATH \$PATH"
go version
go env

export GOPATH=/gopath/
export PATH=\$PATH:/gopath/bin

go get -u github.com/mh-cbon/go-bin-rpm/go-bin-rpm-utils
go-bin-rpm-utils setup-repository -out="`pwd`/rpm" -push -repo=$GH
EOT
set -x

docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./docker.sh"

rm -fr docker.sh

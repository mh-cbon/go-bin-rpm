#!/bin/sh -e

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
#
# to use it
# curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/create-pkg.sh \
# | GH=mh-cbon/gh-api-cli sh -xe

REPO=`echo ${GH} | cut -d '/' -f 2`
USER=`echo ${GH} | cut -d '/' -f 1`


cat <<EOT > docker.sh
dnf install rpm-build -y

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-bin-rpm sh -xe

cd /docker
if ["${REPO}" == "go-bin-rpm"]; then
  VERBOSE=* ./build/amd64/go-bin-rpm generate -a 386 --version ${TRAVIS_TAG} -b pkg-build/386/ -o ${REPO}-386.rpm
  VERBOSE=* ./build/amd64/go-bin-rpm generate -a amd64 --version ${TRAVIS_TAG} -b pkg-build/amd64/ -o ${REPO}-amd64.rpm
else
  VERBOSE=* go-bin-rpm generate -a 386 --version ${TRAVIS_TAG} -b pkg-build/386/ -o ${REPO}-386.rpm
  VERBOSE=* go-bin-rpm generate -a amd64 --version ${TRAVIS_TAG} -b pkg-build/amd64/ -o ${REPO}-amd64.rpm
fi
rm -fr pkg-build/*
EOT

rm -fr pkg-build/*
docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./docker.sh ${TRAVIS_TAG}"
sudo chown travis:travis ${REPO}-*.rpm

rm -fr docker.sh

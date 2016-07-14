
dnf install rpm-build -y
cd /docker
TRAVIS_TAG=0.0.1
VERBOSE=* ./go-bin-rpm generate -a 386 --version ${TRAVIS_TAG} -b pkg-build/386/ -o go-bin-rpm-386.rpm
VERBOSE=* ./go-bin-rpm generate -a amd64 --version ${TRAVIS_TAG} -b pkg-build/amd64/ -o go-bin-rpm-amd64.rpm

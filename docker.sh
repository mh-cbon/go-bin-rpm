
dnf install rpm-build -y
cd /docker
TAG=$1
if [[ -z ${TAG} ]]; then TAG="0.0.0"; fi
VERBOSE=* ./go-bin-rpm generate -a 386 --version ${TAG} -b pkg-build/386/ -o go-bin-rpm-386.rpm
VERBOSE=* ./go-bin-rpm generate -a amd64 --version ${TAG} -b pkg-build/amd64/ -o go-bin-rpm-amd64.rpm

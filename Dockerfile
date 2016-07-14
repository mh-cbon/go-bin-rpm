FROM fedora

ARG TAG=0.0.0

RUN dnf install rpm-build -y
# WORKDIR /docker
RUN VERBOSE=* ./go-bin-rpm generate -a 386 --version ${TAG} -b pkg-build/386/ -o go-bin-rpm-386.rpm
RUN VERBOSE=* ./go-bin-rpm generate -a amd64 --version ${TAG} -b pkg-build/amd64/ -o go-bin-rpm-amd64.rpm

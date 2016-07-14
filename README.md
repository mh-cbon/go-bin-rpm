# go-bin-rpm

Create binary rpm package with ease

Using a `json` files to declare rules, it then performs necessary operations
to invoke `rpmbuild` and build the package.

# Install

__deb/rpm__

```sh
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| sh -x  mh-cbon/go-bin-rpm '${REPO}-${ARCH}${EXT}'
```

__others__

```sh
mkdir -p $GOPATH/src/github.com/mh-cbon
cd $GOPATH/src/github.com/mh-cbon
git clone https://github.com/mh-cbon/go-bin-rpm.git
cd go-bin-rpm
glide install
go install
```

# Requirements

A centos/fedora/redhat system, vagrant, travis, docker, whatever.

# Workflow overview

To create a binary package you need to

- build your application binaries
- invoke `go-bin-rpm` to generate the package

# Usage

```sh
NAME:
   go-bin-rpm - Generate a binary rpm package

USAGE:
   go-bin-rpm <cmd> <options>

VERSION:
   0.0.0

COMMANDS:
     generate-spec  Generate the SPEC file
     generate       Generate the package
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

#### generate-spec

```sh
NAME:
   go-bin-rpm generate-spec - Generate the SPEC file

USAGE:
   go-bin-rpm generate-spec [command options] [arguments...]

OPTIONS:
   --file value, -f value  Path to the rpm.json file (default: "rpm.json")
   -a value, --arch value  Target architecture of the build
   --version value         Target version of the build
```

#### generate

```sh
NAME:
   go-bin-rpm generate - Generate the package

USAGE:
   go-bin-rpm generate [command options] [arguments...]

OPTIONS:
   --file value, -f value        Path to the rpm.json file (default: "rpm.json")
   -b value, --build-area value  Path to the build area (default: "pkg-build")
   -a value, --arch value        Target architecture of the build
   -o value, --output value      Output package to this path
   --version value               Target version of the build
```

#### test

```sh
NAME:
   go-bin-rpm test - Test the package json file

USAGE:
   go-bin-rpm test [command options] [arguments...]

OPTIONS:
   --file value, -f value  Path to the rpm.json file (default: "rpm.json")
```

# Installing generated package

__TLDR__

```sh
# install
sudo rpm -ivh pkg.rpm
# upgrade
sudo rpm -Uvh pkg.rpm
# remove
sudo rpm -evv nx pkg.rpm
```

# Json file

For a real world example including service, shortcuts, env, see [this](demo/rpm.json)

For a casual example to provide a simple binary, see [this](rpm.json)

# Vagrant recipe

Please check the demo app [here](demo/)

# Travis recipe

- get a github repo
- get a travis account
- connect your github account to travis and register your repo
- install travis client `gem install --user travis`
- run `travis setup releases`
- personalize the `.travis.yml`

```yml
sudo: required

services:
  - docker

language: go
go:
  - tip

before_install:
  - sudo apt-get -qq update
  - docker pull fedora

install:
  - glide install

before_deploy:
  - mkdir -p build/{386,amd64}
  - GOOS=linux GOARCH=386 go build -o build/386/program main.go
  - GOOS=linux GOARCH=amd64 go build -o build/amd64/program main.go
  - docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./docker.sh ${TRAVIS_TAG} program"
  - sudo chown travis:travis program-{386,amd64}.rpm

deploy:
  provider: releases
  api_key:
    secure: ... your own here
  file:
    - program-386.deb
    - program-amd64.deb
  skip_cleanup: true
  on:
    tags: true
```

```sh
# docker.sh
#!/bin/sh -xe

dnf install rpm-build -y

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-bin-rpm ASSET='${REPO}-${ARCH}${EXT}' sh -xe

cd /docker
TAG=$1
NAME=$2

if [[ -z ${TAG} ]]; then TAG="0.0.0"; fi

VERBOSE=* go-bin-rpm generate -a 386 --version ${TAG} -b pkg-build/386/ -o ${NAME}-386.rpm
VERBOSE=* go-bin-rpm generate -a amd64 --version ${TAG} -b pkg-build/amd64/ -o ${NAME}-amd64.rpm
```

# useful rpm commands

```sh
# check dependencies before install
rpm -qpR pkg.rpm
# show info of a package before install
rpm -qip pkg.rpm
# install with no dependencies
rpm -ivh --nodeps pkg.rpm
# show info of installed package
rpm -qi pkg
# check installed package
rpm -q pkg
# list files of installed package
rpm -ql pkg
```

# Readings of interest

- https://fedoraproject.org/wiki/Packaging:RPMMacros
- http://www.rpm.org/max-rpm/s1-rpm-build-creating-spec-file.html
- http://www.rpm.org/max-rpm/s1-rpm-inside-files-list-directives.html
- http://www.rpm.org/max-rpm/s1-rpm-inside-scripts.html
- http://www.rpm.org/max-rpm-snapshot/s1-rpm-depend-manual-dependencies.html
- https://fedoraproject.org/wiki/PackagingDrafts:SystemdClarification#Packaging
- https://fedoraproject.org/wiki/Packaging:Scriptlets?rd=Packaging:ScriptletSnippets#Systemd
- https://fedoraproject.org/wiki/Packaging:Guidelines#BuildRequires_2
- http://wiki.networksecuritytoolkit.org/nstwiki/index.php/RPM_Quick_Reference#Secret_.25pretrans_and_.25posttrans_RPM_Scriptlets
- https://fedoraproject.org/wiki/Packaging:Scriptlets?rd=Packaging:ScriptletSnippets#desktop-database
- https://fedoraproject.org/wiki/Archive:PackagingDrafts/DesktopVerify?rd=PackagingDrafts/DesktopVerify
- https://fedoraproject.org/wiki/Archive:PackagingDrafts/DesktopFiles?rd=PackagingDrafts/DesktopFiles

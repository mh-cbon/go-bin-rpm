# go-bin-rpm

Create binary rpm package with ease

Using a `json` files to declare rules, it then performs necessary operations to invoke `rpmbuild` and build the package.

__wip__

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
git clone https://github.com/mh-cbon/go-bin-deb.git
cd go-bin-deb
glide install
go install
```

# Requirements

A centos/fedora/redhat system, vagrant, travis, docker, whatever.

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
language: go
go:
  - tip
before_install:
  - sudo apt-get -qq update
  - sudo apt-get install build-essential lintian -y
  - curl https://glide.sh/get | sh
  - wget -q -O - --no-check-certificate https://raw.githubusercontent.com/mh-cbon/go-bin-deb/master/install.sh | sh
  - wget -q -O - --no-check-certificate https://raw.githubusercontent.com/mh-cbon/changelog/master/install.sh | sh
install:
  - glide install
before_deploy:
  - mkdir -p build/{386,amd64}
  - GOOS=linux GOARCH=386 go build -o build/386/program main.go
  - GOOS=linux GOARCH=amd64 go build -o build/amd64/program main.go
  - go-bin-deb generate -a 386 --version ${TRAVIS_TAG} -w pkg-build-386/ -o ${TRAVIS_BUILD_DIR}/program-386.deb
  - go-bin-deb generate -a amd64 --version ${TRAVIS_TAG} -w pkg-build-amd64/ -o ${TRAVIS_BUILD_DIR}/program-amd64.deb
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
# list iles of installed package
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

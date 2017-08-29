# go-bin-rpm

[![travis Status](https://travis-ci.org/mh-cbon/go-bin-rpm.svg?branch=master)](https://travis-ci.org/mh-cbon/go-bin-rpm) [![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/go-bin-rpm)](https://goreportcard.com/report/github.com/mh-cbon/go-bin-rpm) [![GoDoc](https://godoc.org/github.com/mh-cbon/go-bin-rpm?status.svg)](http://godoc.org/github.com/mh-cbon/go-bin-rpm) [![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)

Create binary rpm package with ease


Using a `json` files to declare rules, it then performs necessary operations
to invoke `rpmbuild` and build the package.

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

See [the demo](demo/).

# TOC
- [Install](#install)
  - [Glide](#glide)
  - [linux rpm/deb repository](#linux-rpmdeb-repository)
  - [linux rpm/deb standalone package](#linux-rpmdeb-standalone-package)
- [Usage](#usage)
  - [Requirements](#requirements)
  - [Workflow overview](#workflow-overview)
  - [Json file](#json-file)
- [JSON tokens](#json-tokens)
- [CLI](#cli)
  - [go-bin-rpm -help](#go-bin-rpm--help)
  - [go-bin-rpm generate-spec -help](#go-bin-rpm-generate-spec--help)
  - [go-bin-rpm generate -help](#go-bin-rpm-generate--help)
  - [go-bin-rpm test -help](#go-bin-rpm-test--help)
- [Recipes](#recipes)
  - [Installing generated package](#installing-generated-package)
  - [Vagrant recipe](#vagrant-recipe)
  - [Travis recipe](#travis-recipe)
  - [useful rpm commands](#useful-rpm-commands)
  - [Readings of interest](#readings-of-interest)
  - [Release the project](#release-the-project)
- [History](#history)

# Install

Check the [release page](https://github.com/mh-cbon/go-bin-rpm/releases)!

#### Glide
```sh
mkdir -p $GOPATH/src/github.com/mh-cbon/go-bin-rpm
cd $GOPATH/src/github.com/mh-cbon/go-bin-rpm
git clone https://github.com/mh-cbon/go-bin-rpm.git .
glide install
go install
```

#### linux rpm/deb repository
```sh
wget -O - https://raw.githubusercontent.com/mh-cbon/latest/master/bintray.sh \
| GH=mh-cbon/go-bin-rpm sh -xe
# or
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/bintray.sh \
| GH=mh-cbon/go-bin-rpm sh -xe
```

#### linux rpm/deb standalone package
```sh
curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-bin-rpm sh -xe
# or
wget -q -O - --no-check-certificate \
https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/go-bin-rpm sh -xe
```

# Usage

### Requirements

A centos/fedora/redhat system, vagrant, travis, docker, whatever.

### Workflow overview

To create a binary package you need to

- build your application binaries
- invoke `go-bin-rpm` to generate the package
- create rpm repositories on `travis` hosted on `gh-pages` using this [script](setup-repository.sh)

### Json file

For a real world example including service, shortcuts, env, see [this](demo/rpm.json)

For a casual example to provide a simple binary, see [this](rpm.json)

## JSON tokens

Several tokens are provided to consume into the JSON file.

|name|description|example|
| --- | --- | -- |
|__!version!__|The `version` number provided in the command line, or in the JSON file.|1.0.2|
|__!arch!__|The `architecture` short name provided in the command line or in the JSON file.|amd64|
|__!name!__|The `name` of the project provided in the JSON file.|hello|

# CLI
#### go-bin-rpm -help
```sh
NAME:
   go-bin-rpm - Generate a binary rpm package

USAGE:
   go-bin-rpm <cmd> <options>
   
VERSION:
   1.0.0
   
COMMANDS:
     generate-spec  Generate the SPEC file
     generate       Generate the package
     test           Test the package json file
     help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
#### go-bin-rpm generate-spec -help
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
#### go-bin-rpm generate -help
```sh
NAME:
   go-bin-rpm generate - Generate the package

USAGE:
   go-bin-rpm generate [command options] [arguments...]

OPTIONS:
   --file value, -f value        Path to the rpm.json file (default: "rpm.json")
   -b value, --build-area value  Path to the build area (default: "pkg-build")
   -a value, --arch value        Target architecture of the build
   -o value, --output value      File path to the resulting rpm file
   --version value               Target version of the build
```
#### go-bin-rpm test -help
```sh
NAME:
   go-bin-rpm test - Test the package json file

USAGE:
   go-bin-rpm test [command options] [arguments...]

OPTIONS:
   --file value, -f value  Path to the rpm.json file (default: "rpm.json")
```

# Recipes

### Installing generated package

__TLDR__

```sh
# install
sudo rpm -ivh pkg.rpm
# upgrade
sudo rpm -Uvh pkg.rpm
# remove
sudo rpm -evv nx pkg.rpm
```

### Vagrant recipe

Please check the demo app [here](demo/)

### Travis recipe

- get a github repo
- get a travis account
- connect your github account to travis and register your repo
- install travis client `gem install --user travis`
- run `travis encrypt --add -r YOUR_USERNAME/dummy GH_TOKEN=xxxx`
- run `travis setup releases`
- personalize the `.travis.yml`

```yml
  sudo: required
  services:
  - docker
  language: go
  go:
  - 1.8
  - 1.9
  env:
    matrix:
    - OKARCH=amd64 OSARCH=amd64
    - OKARCH=386 OSARCH=i386
    global:
    - VERSION=${TRAVIS_TAG}
    - GH_USER=${TRAVIS_REPO_SLUG%/*}
    - GH_APP=${TRAVIS_REPO_SLUG#*/}
    - JFROG_CLI_OFFER_CONFIG=false
    - secure: iTmECM/2XzLk4/bMm4xCYcVyKJj6kV815wLlFyCMzcBHoQmOjYxWu6s36sQ9SecPu/fjE67+rJ6opPpeMz9SRTa34+L64vA0gRA61TynC3NUOIkvzs0y956fhxVGPF5knLnKGeWIHHObfxwo8ovvY+utWI4DqYwfWoIzGZ7Etc1Pp/zY4rIFYiNpRHH5zpvKG3+wjDx6ciUmZASM4ZFEY9OY5LlnOKxaxTFz71i0SZJO18fNTX1NrWEf55pU35CzVpk1rZ1V7cKt1i7k9MqWWwG5RXIzjLWjTx2aFJOchTj2kV7MM8FPk3FcQIiJ0JCPqkaVMNWTa3f++wH85fk2r4ErlVGvtG6xuHAy7uM4GBPX2nmUnEXIlg8qVHgYqUxg18qi92osJIVv3WZgNrL42dRQMHDcEf0b0NMEifpsPES65CnVQPMEZlN+eXzgVMMKS9s/Vf5B8HziDxkTy/crghqMQtff2HZ1eyH1gmLMA9HXQIFWCsAIEfrrgP5doNIfsHJaWP31agSyrHhGMI8y77kVj9K/cxL0XVq9mKlXtnI7S5OKgDD1Dn/PuUPkfrPFQ0wxCfXJ+FyqSyckUd7rTlihJKOlmEvHif3y7BK+qCbj/G7JveQGLKHt/fiIbo00PBvalLSSJDQEsFyO+cKLsR8uVyX/BgaMed6qt2RrzYE=
  before_install:
  - sudo apt-get -qq update
  - mkdir -p ${GOPATH}/bin
  - cd ~
  - curl https://glide.sh/get | sh
  install:
  - cd $GOPATH/src/github.com/${TRAVIS_REPO_SLUG}
  - glide install
  - go install
  script: echo "pass"
  before_deploy:
  - docker pull fedora
  - mkdir -p build/$OSARCH
  - GOOS=linux go build --ldflags "-X main.VERSION=$VERSION" -o build/$OSARCH/$GH_APP
    main.go
  - cp $GOPATH/bin/go-bin-rpm tmp
  - |
    docker run -v $PWD:/mnt/travis fedora /bin/sh -c "cd /mnt/travis && (curl -s -L https://bintray.com/bincrafters/public-rpm/rpm > /etc/yum.repos.d/w.repo) && dnf install changelog rpm-build -y --quiet && ./tmp generate --file rpm.json -a $OSARCH --version $VERSION -o $GH_APP-$OSARCH-$VERSION.rpm"
  - rm -f ./tmp
  - cp $GH_APP-$OSARCH-$VERSION.rpm $GH_APP-$OKARCH.rpm
  - curl -fL https://getcli.jfrog.io | sh
  - ./jfrog bt pc --key=$BTKEY --user=$GH_USER --licenses=MIT --vcs-url=https://github.com/$GH_USER/rpm
    $GH_USER/rpm/$GH_APP || echo "package already exists"
  - ./jfrog bt upload --override=true --key $BTKEY --publish=true --deb=unstable/main/$OSARCH
    $GH_APP-$OSARCH-$VERSION.rpm $GH_USER/rpm/$GH_APP/$VERSION pool/g/$GH_APP/
  deploy:
    provider: releases
    api_key:
      secure: CY2nebPdr2CSCZW34QCtlw/IdbaHl5T77xPFlmvXB2Z+0SnO0RTW7JvFMa2mDYxa6ibZ6dR2br9YwdgJYnqV+PnXCizvZ5KPqpHxE31ta4s1IokZr+v9J+deGvUdk60oF5mxkqcGgAtScEGC5ZVJ/0EqAn64o4+H3fOQfA1pYTpzUBL/c9yUNqAFLFDVXz1sd7eSccPwf1uthdhndybMgatogfQuUBmm3vNJYYheAF8XCimBmrsIkPed+OKfhkDqUCTdgSTOQWvv0Uf8ib5VUH0w+UV8Wx69/KNKVhp/f7Nhf6GCKT1AKh/fQxjpRaWdkQLsn7nqPVuF0dHYV/mtdo4EP0FDj+2a3LvtGpEst90Mo0SRzauhqCQqCopyOf3JKkKPqTyMRDKAzYWAymjeLGaPda4wOxNROWV7yBuXNTTUmU2GDPUMULnLA7v+0ml6wd3gGCOMU5It8Iynkuxts8ATlpa0qels3memQITfhkTdR3CFT2mr/frkDiVOtqnp6BJoQIjhSMXoMRfnSpnNOszsiLNa9pM+hNG3HeZN0MQ+gTlRgqmTSitvllr751oUhgNzjv35FDxaywFwKlqtaJfX9UVCLxcBTvDcP4ZKHJRbgFOmffv2mnKi1S8K26LUkuLZDKvCZgrw8iM1KjvPX/GP9tXaxgLrfsfQOcOGGGs=
    file_glob: true
    file:
    - $GH_APP-$OKARCH.rpm
    skip_cleanup: true
    overwrite: true
    true:
      tags: true
```

### useful rpm commands

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

### Readings of interest

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

### Release the project

```sh
gump patch -d # check
gump patch # bump
```

# History

[CHANGELOG](CHANGELOG.md)

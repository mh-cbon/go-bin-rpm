# {{.Name}}

{{template "badge/travis" .}}{{template "badge/goreport" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

Using a `json` files to declare rules, it then performs necessary operations
to invoke `rpmbuild` and build the package.

This tool is part of the [go-github-release workflow](https://github.com/mh-cbon/go-github-release)

See [the demo](demo/).

# {{toc 5}}

# Install

{{template "gh/releases" .}}

#### Glide
{{template "glide/install" .}}

#### linux rpm/deb repository
{{template "linux/gh_src_repo" .}}

#### linux rpm/deb standalone package
{{template "linux/gh_pkg" .}}

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

# CLI

{{exec "go-bin-rpm" "-help" | color "sh"}}

{{exec "go-bin-rpm" "generate-spec" "-help" | color "sh"}}

{{exec "go-bin-rpm" "generate" "-help" | color "sh"}}

{{exec "go-bin-rpm" "test" "-help" | color "sh"}}

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
  - tip

env:
  global:
    - MYAPP=dummy
    - MYEMAIL=some@email.com
    - secure: GH_TOKEN

before_install:
  - sudo apt-get -qq update
  - mkdir -p ${GOPATH}/bin

install:
  - cd $GOPATH/src/github.com/mh-cbon/$MYAPP
  - go install

script: echo "pass"

before_deploy:
  - docker pull fedora
  - mkdir -p build/{386,amd64}
  - GOOS=linux GOARCH=386 go build --ldflags "-X main.VERSION=${TRAVIS_TAG}" -o build/386/$MYAPP main.go
  - GOOS=linux GOARCH=amd64 go build --ldflags "-X main.VERSION=${TRAVIS_TAG}" -o build/amd64/$MYAPP main.go
  - curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/create-pkg.sh | GH=YOUR_USERNAME/$MYAPP sh -xe

after_deploy:
  - curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/setup-repository.sh | GH=YOUR_USERNAME/$MYAPP EMAIL=$MYEMAIL sh -xe

deploy:
  provider: releases
  api_key:
    secure: GH_TOKEN xxxx
  file_glob: true
  file:
    - $MYAPP-386.rpm
    - $MYAPP-amd64.rpm
  skip_cleanup: true
  on:
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

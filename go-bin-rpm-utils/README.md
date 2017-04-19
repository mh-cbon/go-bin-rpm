# go-bin-rpm-utils

[![MIT License](http://img.shields.io/badge/License-MIT-yellow.svg)](../LICENSE)

go-bin-rpm-utils is a cli tool to generate rpm package and repos.


# Usage

```sh
export GH_TOKEN=`gh-api-cli get-auth -n release`

vagrant up

vagrant rsync && vagrant ssh -c "export GH_TOKEN=$GH_TOKEN; sh /vagrant/vagrant-run.sh"

vagrant destroy -f
```

# See also

https://github.com/mh-cbon/go-bin-rpm#recipes

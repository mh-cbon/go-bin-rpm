---
License: MIT
LicenseFile: ../LICENSE
LicenseColor: yellow
---
# {{.Name}}

{{template "license/shields" .}}

{{pkgdoc}}

# Usage

```sh
export GH_TOKEN=`gh-api-cli get-auth -n release`

vagrant up

vagrant rsync && vagrant ssh -c "export GH_TOKEN=$GH_TOKEN; sh /vagrant/vagrant-run.sh"

vagrant ssh -c 'curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/go-bin-rpm sh -xe'
vagrant ssh -c 'go-bin-rpm -v'
vagrant ssh -c 'which go-bin-rpm'
vagrant ssh -c 'sudo yum remove go-bin-rpm -y'
vagrant ssh -c 'curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/source.sh | GH=mh-cbon/go-bin-rpm sh -xe'
vagrant ssh -c 'go-bin-rpm -v'
vagrant ssh -c 'which go-bin-rpm'
vagrant ssh -c 'sudo yum remove go-bin-rpm -y'

vagrant destroy -f
```

# See also

https://github.com/mh-cbon/go-bin-rpm#recipes

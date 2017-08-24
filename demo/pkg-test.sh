set -e
set -x

curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh \
| GH=mh-cbon/changelog sh -xe || echo "ok"

cd /vagrant/
rm -fr /vagrant/pkg-build
rm -f *.rpm


VERBOSE=* ./go-bin-rpm generate-spec -a 386 --version 0.0.1

VERBOSE=* ./go-bin-rpm generate -a 386 --version 0.0.1 -b pkg-build/386/ -o hello-386.rpm
ls -alh

VERBOSE=* ./go-bin-rpm generate -a amd64 --version 0.0.1 -b pkg-build/amd64/ -o hello-amd64.rpm
ls -alh

mkdir tomate
VERBOSE=* ./go-bin-rpm generate -a amd64 --version 0.0.1 -b pkg-build/amd64/ -o tomate/hello-amd64.rpm
ls -alh
ls -alh tomate

# remove the package
sudo rpm -e hello || echo "ok"

sudo systemctl daemon-reload
systemctl status hello || echo "ok"

# install test
sudo rpm -ivh --test /vagrant/hello-amd64.rpm
# install
sudo rpm -ivh /vagrant/hello-amd64.rpm

which hello
echo \$some
systemctl status hello
curl http://localhost:8080/

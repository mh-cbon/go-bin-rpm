set -e
set -x

rm -fr build && mkdir -p build/{386,amd64}
rm -fr pkg-build && mkdir -p pkg-build/{386,amd64}

GOOS=linux GOARCH=386 go build -o build/386/hello hello.go
GOOS=linux GOARCH=amd64 go build -o build/amd64/hello hello.go

go build -o go-bin-rpm ../main.go

vagrant up cli
vagrant ssh -c 'sudo yum install rpm-build desktop-file-utils -y' cli

vagrant rsync cli
vagrant ssh -c 'sh /vagrant/pkg-test.sh' cli

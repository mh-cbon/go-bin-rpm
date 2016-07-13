# hello package - a demo

`hello` is a program which serves a web server on port 8080.

The package installs
- the program bin and its assets
- a service unit file for systemd
- a desktop link to open the hello homepage
- and environment variable

# run it

```sh
rm -fr build && mkdir -p build/{386,amd64}
rm -fr build && mkdir -p pkg-build/{386,amd64}
GOOS=linux GOARCH=386 go build -o build/386/hello hello.go
GOOS=linux GOARCH=amd64 go build -o build/amd64/hello hello.go
./go-bin-rpm generate-spec -a 386 -v 0.0.1

vagrant up cli
vagrant ssh -c 'sudo yum install rpm-build desktop-file-utils -y' cli

go build -o go-bin-rpm ../main.go && vagrant rsync
vagrant ssh -c 'rm -fr /vagrant/pkg-build' cli

vagrant ssh -c 'rm -f hello-386.rpm' cli
vagrant ssh -c 'cd /vagrant/ && VERBOSE=* ./go-bin-rpm generate -a 386 --version 0.0.1 -b pkg-build/386/ -o hello-386.rpm' cli
vagrant ssh -c 'cd /vagrant/ && ls -alh' cli

vagrant ssh -c 'rm -f hello-amd64.rpm' cli
vagrant ssh -c 'cd /vagrant/ && VERBOSE=* ./go-bin-rpm generate -a amd64 --version 0.0.1 -b pkg-build/amd64/ -o hello-amd64.rpm' cli
vagrant ssh -c 'cd /vagrant/ && ls -alh' cli


vagrant ssh -c 'sudo rpm -e hello' cli
vagrant ssh -c 'sudo rpm -ivh --test /vagrant/hello-amd64.rpm' cli
vagrant ssh -c 'sudo rpm -ivh /vagrant/hello-amd64.rpm' cli
vagrant ssh -c 'hello' cli # should bug as the service is started
vagrant ssh -c "echo \$some" cli
vagrant ssh -c "systemctl status hello" cli
vagrant ssh -c "curl http://localhost:8080/" cli
```


prebump=
  666 git fetch --tags origin master
  666 git pull origin master

preversion=
  philea -s "666 go vet %s" "666 go-fmt-fail %s"
  666 go run main.go -v
  666 changelog finalize --version !newversion!
  666 commit -q -m "changelog: !newversion!" -f change.log
  666 rm-glob -r build/
  666 rm-glob -r assets/
  666 build-them-all build main.go -o "build/&os-&arch/&pkg" --ldflags "-X main.VERSION=!newversion!"

postversion=
  666 changelog md -o CHANGELOG.md --vars='{"name":"go-bin-rpm"}'
  666 commit -q -m "changelog: !newversion!" -f CHANGELOG.md
  666 go install --ldflags "-X main.VERSION=!newversion!"
  666 emd gen -out README.md
  666 commit -q -m "README: !newversion!" -f README.md
  666 git push
  666 git push --tags
  philea -s -S -p "build/windows*/*" "666 archive create -f -o=assets/%dname.zip -C=build/%dname/ ."
  philea -s -S -e windows*/** -p "build/*/**" "666 archive create -f -o=assets/%dname.tar.gz -C=build/%dname/ ."
  666 gh-api-cli create-release -n release -o mh-cbon -r go-bin-rpm \
    --ver !newversion! -c "changelog ghrelease --version !newversion!" \
    --draft !isprerelease!
  666 gh-api-cli upload-release-asset -n release --glob "assets/*" -o mh-cbon -r go-bin-rpm --ver !newversion!
  666 rm-glob -r build/
  666 rm-glob -r assets/

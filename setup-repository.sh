#!/bin/sh -e

# this is an helper
# to use into your travis file
# it is limited to amd64/386 arch
#
# to use it
# curl -L https://raw.githubusercontent.com/mh-cbon/go-bin-rpm/master/setup-deb-repository.sh \
# | GH=mh-cbon/gh-api-cli EMAIL=mh-cbon@users.noreply.github.com sh -xe

# GH=$1
# EMAIL=$2

REPO=`echo ${GH} | cut -d '/' -f 2`
USER=`echo ${GH} | cut -d '/' -f 1`

# clean up build.
rm -fr ${REPO}-*.rpm
rm -fr ${REPO}-*.deb

sudo apt-get install build-essential -y

if type "gh-api-cli" > /dev/null; then
  echo "gh-api-cli already installed"
else
  curl -L https://raw.githubusercontent.com/mh-cbon/latest/master/install.sh | GH=mh-cbon/gh-api-cli sh -xe
fi

git fetch --all
git config user.name "${USER}"
git config user.email "${EMAIL}"
if [ `git branch --list gh-pages `]
then
  git checkout gh-pages
else
  git checkout -b gh-pages
  find . -maxdepth 1 -mindepth 1 -not -name .git -exec rm -rf {} \;
  git commit -am "clean up"
fi

rm -fr rpm
mkdir -p rpm/{i386,x86_64}
gh-api-cli dl-assets -o ${USER} -r ${REPO} --out rpm/i386/%r-%v_%a.%e -g "*386*rpm"
gh-api-cli dl-assets -o ${USER} -r ${REPO} --out rpm/x86_64/%r-%v_%a.%e -g "*amd64*rpm"

cat <<EOT > createrepo.sh
yum install createrepo -y
cd /docker/rpm/i386
createrepo .
cd /docker/rpm/x86_64
createrepo .
EOT
docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./createrepo.sh"


cat <<EOT > gen-repo-file.sh
DESC=`rpm -qip *.rpm | grep Summary | cut -d ':' -f2 | cut -d ' ' -f2- | tail -n 1`
cat <<EOTin > rpm/${REPO}.repo
[${REPO}]
name=${DESC}
baseurl=https://${USER}.github.io/${REPO}/rpm/\$basearch/
enabled=1
skip_if_unavailable=1
gpgcheck=0
EOTin
EOT
docker run -v $PWD:/docker fedora /bin/sh -c "cd /docker && sh ./gen-repo-file.sh"

rm -f gen-repo-file.sh
rm -f createrepo.sh


git add -A
git commit -m "Created rpm repository"

git status
git branch

set +x # disable debug output because that would display the token in clear text..
echo "git push --force --quiet https://GH_TOKEN@github.com/${GH}.git gh-pages"
git push --force --quiet "https://${GH_TOKEN}@github.com/${GH}.git" gh-pages \
 2>&1 | sed -re "s/${GH_TOKEN}/GH_TOKEN/g"

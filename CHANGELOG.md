# Changelog - go-bin-rpm

### 0.0.16-beta

__Changes__

- travis upload token
- utils: updated utility scripts
- spec generation: the spec generation does not fail anymore if the build area does not exists, this helps to preview the spec file

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 20 Apr 2017 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.15...0.0.16-beta#diff)
______________

### 0.0.15

__Changes__

- travis: fix missing assets
- packaging: trick the build script to be able to build go-bin-rpm itself

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.14...0.0.15#diff)
______________

### 0.0.14

__Changes__

- travis: update build file, add deb/rpm repositories setup
- source repository: use token when downloading assets
- source repository: fix description reading of rpm package
- source repository: fix base arch escaping in the creation of the repo file
- support: add centralized script to generate packages and repository

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 30 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.13...0.0.14#diff)
______________

### 0.0.13

__Changes__

- rpmbuild: fix src file when copying rpmbuild file result

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.12...0.0.13#diff)
______________

### 0.0.12

__Changes__

- fix int to string conversion

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.11...0.0.12#diff)
______________

### 0.0.11

__Changes__

- spec file: fix version field, it must not contain prerelease or metadata

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.10...0.0.11#diff)
______________

### 0.0.10

__Changes__

- rpm: fix changelog command to generate rpm changelog
- changelog: fix changelog format

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.9...0.0.10#diff)
______________

### 0.0.9

__Changes__

- Spec file generator: Version field must not contain prerelease information,
  if the version contains a prerelease information,
  the value is now recorded into Release field.
- glide: add semver dependency
- packaging: add deb package support
- travis: updated changelog url installer
- README: install section

__Contributors__

- mh-cbon

Released by mh-cbon, Sat 23 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.8...0.0.9#diff)
______________

### 0.0.8

__Changes__

- docker: Fix docker.sh example file
- README: fix dockersh example

__Contributors__

- mh-cbon

Released by mh-cbon, Fri 15 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.7...0.0.8#diff)
______________

### 0.0.7

__Changes__

- travis: update build script
- README: update travis section

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.6...0.0.7#diff)
______________

### 0.0.6

__Changes__

- travis: update build script

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.5...0.0.6#diff)
______________

### 0.0.5

__Changes__

- travis: fix deploy section

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.4...0.0.5#diff)
______________

### 0.0.4

__Changes__

- travis: fix deploy section

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.3...0.0.4#diff)
______________

### 0.0.3

__Changes__

- travis: fix deploy section

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.2...0.0.3#diff)
______________

### 0.0.2

__Changes__

- release: fix package name

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/0.0.1...0.0.2#diff)
______________

### 0.0.1

__Changes__

- Initial release

__Contributors__

- mh-cbon

Released by mh-cbon, Thu 14 Jul 2016 -
[see the diff](https://github.com/mh-cbon/go-bin-rpm/compare/62d0a83570d6d0f310d675e2a760e7552a40c63a...0.0.1#diff)
______________



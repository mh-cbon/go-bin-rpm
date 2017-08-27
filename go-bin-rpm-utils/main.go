// go-bin-rpm-utils is a cli tool to generate rpm package and repos.
package main

import (
	"flag"
	"os"
)

func main() {

	flag.Parse()
	action := flag.Arg(0)

	// basic arg parsing
	var reposlug string
	var email string
	var ghToken string
	var version string
	var archs string
	var out string

	flag.StringVar(&reposlug, "repo", "", "The repo slug such USER/REPO.")
	flag.StringVar(&ghToken, "ghToken", "", "The ghToken to write on your repository.")
	flag.StringVar(&email, "email", "", "Your gh email.")
	flag.StringVar(&version, "version", "", "The package version.")
	flag.StringVar(&archs, "archs", "386,amd64", "The archs to build.")
	flag.StringVar(&out, "out", "", "The out build directory.")
	push := flag.Bool("push", false, "Push the new assets")
	keep := flag.Bool("keep", false, "Keep the new assets")
	flag.CommandLine.Parse(os.Args[2:])

	// os.Env fallback
	email = readEnv(email, "EMAIL", "MYEMAIL")
	reposlug = readEnv(reposlug, "REPO")
	ghToken = readEnv(ghToken, "GH_TOKEN")

	// ci fallback
	// todo: make use of pre defined ci env
	if isTravis() {
		version = readEnv(version, "TRAVIS_TAG")
		out = readEnv(out, "TRAVIS_BUILD_DIR")
	}
	if isVagrant() {
		version = readEnv(version, "VERSION")
		out = readEnv(out, "BUILD_DIR")
	}

	// integrity check
	requireArg(reposlug, "repo", "REPO")
	if *push {
		requireArg(ghToken, "ghToken", "GH_TOKEN")
	}
	// requireArg(email, "email", "EMAIL", "MYEMAIL")

	if isTravis() {
		requireArg(version, "version", "TRAVIS_TAG")
		requireArg(out, "out", "TRAVIS_BUILD_DIR")

	} else if isVagrant() {
		requireArg(version, "version", "VERSION")
		requireArg(out, "out", "BUILD_DIR")

	} else {
		panic("nop, no such ci system...")
	}

	// execute some common setup, in case.
	alwaysHide[ghToken] = "$GH_TOKEN"

	// removeAll(out)
	mkdirAll(out)

	if version == "LAST" {
		version = latestGhRelease(reposlug)
	}

	// execute the action
	if action == "create-packages" {
		CreatePackage(reposlug, ghToken, email, version, archs, out, *push, *keep)

	} else if action == "setup-repository" {
		SetupRepo(reposlug, ghToken, email, version, archs, out, *push, *keep)
	}
}

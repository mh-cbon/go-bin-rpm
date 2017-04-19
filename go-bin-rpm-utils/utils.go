package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/mh-cbon/go-bin-rpm/stringexec"
)

func pushAssetsGh(version, ghToken, outbuild, glob string) {
	if tryexec(`gh-api-cli -v`) != nil {
		exec(`latest -repo=%v`, "mh-cbon/gh-api-cli")
	}
	exec(`gh-api-cli rm-assets --guess --ver %v -t %v --glob %q`, version, ghToken, glob)
	exec(`gh-api-cli upload-release-asset --guess --ver %v -t %v --glob %q`, version, ghToken, outbuild+"/"+glob)
	exec(`rm -f %v`, outbuild+"/"+glob)
}

func setupGitRepo(repoPath, reposlug, user, email string) {
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		mkdirAll(repoPath)
		chdir(repoPath)
		exec(`git clone https://github.com/%v.git .`, reposlug)
		exec(`git config user.name %v`, user)
		exec(`git config user.email %v`, email)
	}
}

func resetGit(repoPath string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	chdir(repoPath)
	exec(`git reset HEAD --hard`)
	exec(`git clean -ffx`)
	exec(`git clean -ffd`)
	exec(`git clean -ffX`)
	chdir(cwd)
}

func getBranchGit(repoPath, reposlug, branch, origin string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	chdir(repoPath)
	tryexec(`git remote rm %v`, origin)
	tryexec(`git remote add %v https://github.com/%v.git`, origin, reposlug)
	tryexec("yes | git fetch %v", origin)

	if tryexec(`git checkout %v`, branch) != nil {
		exec(`git checkout -b %v`, branch)
		// if new, clean it
		exec(`git rm --cached $(git ls-files)`)
		exec(`git status`)
		exec(`git commit -am clean`)
	}
	chdir(cwd)
}

func commitPushGit(repoPath, ghToken, reposlug, branch, message string) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	chdir(repoPath)
	exec(`git add -A`)
	exec(`git commit -m %q`, message)
	u := fmt.Sprintf(`https://%v@github.com/%v.git`, ghToken, reposlug)
	exec(`git push --force --quiet %q %v`, u, branch)
	chdir(cwd)
}

func requireArg(val, n string, env ...string) {
	if val == "" {
		log.Printf("missing argument -%v or env %q\n", n, env)
		os.Exit(1)
	}
}

func readEnv(c string, k ...string) string {
	if c == "" {
		for _, kk := range k {
			c = os.Getenv(kk)
			if c != "" {
				break
			}
		}
	}
	return c
}

func mkdirAll(f string) error {
	fmt.Println("mkdirAll", f)
	return os.MkdirAll(f, os.ModePerm)
}
func removeAll(f string) error {
	if f == "" {
		panic("nop")
	}
	if f == "." {
		panic("nop .")
	}
	fmt.Println("removeAll", f)
	return tryexec("rm -fr %v", f)
}
func chdir(f string) error {
	fmt.Println("Chdir", f)
	return os.Chdir(f)
}

func isTravis() bool {
	return strings.ToLower(os.Getenv("CI")) == "true" &&
		strings.ToLower(os.Getenv("TRAVIS")) == "true"
}

func isVagrant() bool {
	_, s := os.Stat("/vagrant/")
	return !os.IsNotExist(s)
}

func latestGhRelease(repo string) string {
	ret := ""
	u := fmt.Sprintf(`https://api.github.com/repos/%v/releases/latest`, repo)
	fmt.Println("latestGhRelease", u)
	r := getURL(u)
	k := map[string]interface{}{}
	json.Unmarshal(r, &k)

	if x, ok := k["tag_name"]; ok {
		ret = x.(string)
	} else {
		panic("latest version not found")
	}
	fmt.Println("latestGhRelease", ret)
	return ret
}

var alwaysHide = map[string]string{}

func clean(s string) string {
	for search, replace := range alwaysHide {
		s = strings.Replace(s, search, replace, -1)
	}
	return s
}

func tryexec(w string, params ...interface{}) error {
	w = fmt.Sprintf(w, params...)
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("+", clean(w))
	cmd, err := stringexec.Command(cwd, w)
	if err != nil {
		return err
	}
	// out, err := cmd.CombinedOutput()
	// sout := string(out)
	// fmt.Println(clean(sout))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func exec(w string, params ...interface{}) {
	if err := tryexec(w, params...); err != nil {
		panic(err)
	}
}

func gettryexec(w string, params ...interface{}) ([]byte, error) {
	w = fmt.Sprintf(w, params...)
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	fmt.Println("+", clean(w))
	cmd, err := stringexec.Command(cwd, w)
	if err != nil {
		return nil, err
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	// out, err := cmd.CombinedOutput()
	// sout := string(out)
	// fmt.Println(clean(sout))
	return cmd.CombinedOutput()
}

func getexec(w string, params ...interface{}) string {
	o, err := gettryexec(w, params...)
	if err != nil {
		panic(err)
	}
	return string(o)
}

func writeFile(f, content string) {
	fmt.Println("writeFile", f)
	err := ioutil.WriteFile(f, []byte(content), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func getURL(u string) []byte {
	response, err := http.Get(u)
	fmt.Println("getURL", u)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	var ret bytes.Buffer
	_, err = io.Copy(&ret, response.Body)
	if err != nil {
		panic(err)
	}
	return ret.Bytes()
}

func dlURL(u, to string) bool {
	fmt.Println("dlURL ", u)
	fmt.Println("to ", to)
	response, err := http.Get(u)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	f, err := os.Create(to)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = io.Copy(f, response.Body)
	if err != nil {
		panic(err)
	}
	return true
}

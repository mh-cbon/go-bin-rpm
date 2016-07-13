package stringexec

import (
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
)

// Return a new exec.Cmd object for the given command string
func Command(cwd string, cmd string) (*exec.Cmd, error) {
	if runtime.GOOS == "windows" {
		return ExecStringWindows(cwd, cmd)
	}
	return ExecStringFriendlyUnix(cwd, cmd)
}

func ExecStringWindows(cwd string, cmd string) (*exec.Cmd, error) {
	dir, err := ioutil.TempDir("", "stringexec")
	if err != nil {
		return nil, err
	}
	err = ioutil.WriteFile(dir+"/some.bat", []byte(cmd), 0766)
	if err != nil {
		return nil, err
	}

	oCmd := exec.Command("cmd", []string{"/C", dir + "/some.bat"}...)
	oCmd.Dir = cwd
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	// defer os.Remove(tmpfile.Name()) // clean up // not sure how to clean it :x
	return oCmd, nil
}

func ExecStringFriendlyUnix(cwd string, cmd string) (*exec.Cmd, error) {
	oCmd := exec.Command("sh", []string{"-c", cmd}...)
	oCmd.Dir = cwd
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	return oCmd, nil
}

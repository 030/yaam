package yaamtest

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/030/yaam/internal/pkg/project"
)

var npmrc = `registry=` + project.Url + `/npm/3rdparty-npm/
always-auth=true
_auth=aGVsbG86d29ybGQ=`

func NpmConfig() (int, error) {
	if err := os.RemoveAll(filepath.Join(npmHomeDemoProject, "node_modules")); err != nil {
		return 1, err
	}
	packageLockJson := filepath.Join(npmHomeDemoProject, "package-lock.json")
	if _, err := os.Stat(packageLockJson); err == nil {
		if err := os.Remove(packageLockJson); err != nil {
			return 1, err
		}
	}

	npmrcWithCacheLocation := npmrc + `
cache=` + testDirNpm + `/cache` + time.Now().Format("20060102150405111") + ``
	if err := os.WriteFile(filepath.Join(npmHomeDemoProject, ".npmrc"), []byte(npmrcWithCacheLocation), 0600); err != nil {
		return 1, err
	}

	cmd := exec.Command("bash", "-c", "npm cache clean --force && npm i")
	cmd.Dir = npmHomeDemoProject
	co, err := cmd.CombinedOutput()
	if err != nil {
		return cmd.ProcessState.ExitCode(), fmt.Errorf(cmdExitErrMsg, string(co), err)
	}

	return 0, nil
}

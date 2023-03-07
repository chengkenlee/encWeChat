/*
 *@author ChengKen
 *@date   15/02/2023 12:06
 */
package conn

import (
	"bytes"
	"enc/util"
	"os/exec"
	"runtime"
	"strings"
)

func Runshell(command string) string {
	var cmd *exec.Cmd
	var outbuf, errbuf bytes.Buffer

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "echo yarn app -kill xxxx")
	}
	if runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/bash", "-c", command)
	}
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	err := cmd.Start()
	if err != nil {
		util.Logger.Warn(err.Error())
		return ""
	}
	err = cmd.Wait()
	if err != nil {
		util.Logger.Warn(err.Error())
		return ""
	}
	return strings.ReplaceAll(outbuf.String(), "\n", "")
}

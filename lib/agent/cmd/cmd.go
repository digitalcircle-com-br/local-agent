package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/digitalcircle-com-br/local-agent/lib/agent/config"
	"github.com/digitalcircle-com-br/local-agent/lib/common"
	"github.com/skratchdot/open-golang/open"
)

func openCmd(in []string) {
	open.Run(in[1])
}
func execCmd(in []string) {
	var cmd *exec.Cmd
	if len(in) < 2 {

		cmd = exec.Command(in[0])
	} else {
		cmd = exec.Command(in[0], in[1:]...)
	}

	go func() {
		cmd.Stdout = log.Writer()
		cmd.Stderr = log.Writer()
		err := cmd.Start()
		if err != nil {
			log.Printf("Error running cmd %s: %s", strings.Join(in, " "), err.Error())
		}
	}()
}

func replaceVar(in string, c *common.CmdReq) string {
	ret := in
	ret = strings.ReplaceAll(ret, "${USER}", config.Cfg.User)
	ret = strings.ReplaceAll(ret, "${REQID}", c.ReqID)
	for k, v := range config.Cfg.Vars {
		ret = strings.ReplaceAll(ret, fmt.Sprintf("${%s}", k), v)
	}
	for i := range c.Params {
		ret = strings.ReplaceAll(ret, fmt.Sprintf("${P%v}", i), c.Params[i])
	}
	return ret
}
func replaceVars(in []string, c *common.CmdReq) []string {
	ret := make([]string, len(in))
	for i := range in {
		ret[i] = replaceVar(in[i], c)
	}
	return ret
}

func Exec(c *common.CmdReq) {
	cmd, ok := config.Cfg.Cmds[c.Cmd]
	if !ok {
		return
	}
	cmd = replaceVars(cmd, c)

	log.Printf("Calling cmd: [%s]", strings.Join(cmd, ","))

	if cmd[0] == "@open" {
		openCmd(cmd)
	} else {
		execCmd(cmd)
	}
}

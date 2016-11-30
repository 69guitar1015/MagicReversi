package main

import (
	"strings"

	do "gopkg.in/godo.v2"
)

// EdisonHost is host name of edison
var EdisonHost = "edison.local"

// EdisonUser is user name of edison
var EdisonUser = "root"

// AppName is Application Name
var AppName = "MagicReversi"

var options = do.M{
	"app_name":    AppName,
	"edison_host": EdisonUser + "@" + EdisonHost,
	"edison_user": EdisonUser,
}

func tasks(p *do.Project) {

	p.Task("build", nil, func(c *do.Context) {
		c.Run("GOOS=linux GOARCH=386 go build")
	})

	p.Task("copy_to_edison", nil, func(c *do.Context) {
		c.Run("scp {{.app_name}} {{.edison_host}}:/home/{{.edison_user}}", options)
	})

	p.Task("killall_process", nil, func(c *do.Context) {
		// c.Run("ssh {{.edison_host}} killall -q {{.app_name}}", options)
		ps := c.RunOutput("ssh {{.edison_host}} 'ps | grep {{.app_name}}'", options)

		rows := strings.Split(ps, "\n")

		ids := []string{}
		for _, row := range rows {
			if strings.Contains(row, AppName) && !strings.Contains(row, "grep") {
				ids = append(ids, strings.Split(strings.Trim(row, " \n"), " ")[0])
			}
		}

		if len(ids) != 0 {
			addOptions := options
			addOptions["ps_ids"] = strings.Join(ids, " ")

			c.Run("ssh {{.edison_host}} 'kill {{.ps_ids}}'", addOptions)
		}
	})

	p.Task("exec_process", nil, func(c *do.Context) {
		c.Run("ssh {{.edison_host}} 'nohup /home/{{.edison_user}}/{{.app_name}} > /home/{{.edison_user}}/output.log 2>&1 &'", options)
	})

	defaultTask := do.S{
		"build",
		"killall_process",
		"copy_to_edison",
		"exec_process",
	}

	p.Task("default", defaultTask, nil).Src("**/*.go")
}

func main() {
	do.Godo(tasks)
}

package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/xujihui1985/mydocker/cgroup"
	"github.com/xujihui1985/mydocker/cgroup/subsystems"
	"github.com/xujihui1985/mydocker/container"
	"log"
	"os"
	"strings"
)

const (
	usage = `mydocker is a container runtime implementation`
)

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		runCommand,
		initCommand,
	}

	app.Before = func(context *cli.Context) error {
		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

var runCommand = cli.Command{
	Name: "run",
	Usage: "create a container with namespace and cgroups",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: "ti",
			Usage: "enable tty",
		},
		cli.StringFlag{
			Name: "m",
			Usage: "memory limit",
		},
	},

	Action: func(ctx *cli.Context) error {
		if len(ctx.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		var cmdArray []string
		for _, arg := range ctx.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := ctx.Bool("ti")
		resCfg := subsystems.ResourceConfig{
			MemoryLimit: ctx.String("m"),
		}
		Run(tty, cmdArray, &resCfg)
		return nil
	},
}

var initCommand = cli.Command{
	Name: "init",
	Usage: "init container process",

	Action: func(ctx *cli.Context) error {
		cmd := ctx.Args().Get(0)
		log.Printf("command %s\n", cmd)
		return container.RunContainerInitProcess()
	},
}

func Run(tty bool, commandArr []string, cfg *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if err := parent.Start(); err != nil {
		log.Fatal(err)
	}

	cgroupManager := cgroup.NewManager("mydocker-cgroup")
	defer cgroupManager.Destroy()

	cgroupManager.Set(cfg)
	cgroupManager.Apply(parent.Process.Pid)

	sendInitCommand(commandArr, writePipe)
	parent.Wait()
	os.Exit(-1)
}

func sendInitCommand(cmdArray []string, writePipe *os.File) {
	defer writePipe.Close()
	cmd := strings.Join(cmdArray, " ")
	writePipe.WriteString(cmd)
}

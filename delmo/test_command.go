package delmo

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/mitchellh/cli"
)

type TestCommand struct {
	Ui cli.Ui
}

func (t *TestCommand) Help() string {
	helpText := `
Usage: delmo test [options]

  Run a test :-)
`
	return strings.TrimSpace(helpText)
}

func (t *TestCommand) Run(args []string) int {
	flags := flag.FlagSet{
		Usage: func() { t.Help() },
	}

	var path string
	flags.StringVar(&path, "f", "delmo.yml", "")
	if err := flags.Parse(args); err != nil {
		return 1
	}

	suite, err := Load(path)
	if err != nil {
		t.Ui.Error(fmt.Sprintf("%v", err))
		return 1
	}
	dockerCompose, err := NewDockerCompose(suite.ComposeFile, "delmo")
	if err != nil {
		t.Ui.Error(fmt.Sprintf("%s", err))
		return 1
	}
	t.Ui.Info(fmt.Sprintf("Starting System %s", suite.System.Name))
	dockerCompose.Start()
	t.Ui.Info("Waiting 2 seconds")
	time.Sleep(2 * time.Second)
	t.Ui.Info("Stopping System")
	dockerCompose.Stop()
	t.Ui.Info("Reading output")
	out, err := ioutil.ReadAll(dockerCompose.Output())
	t.Ui.Info(string(out))
	return 0

}

func (t *TestCommand) Synopsis() string {
	return "Run some tests"
}

package cliplus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type DashCmd struct {
	Name           string
	Dash           string
	MainResolver   MainResolver
	NoArgsMainFunc func()
}

func NewDashCmd(name string, resolver MainResolver) *DashCmd {
	return &DashCmd{
		Name:           name,
		Dash:           "-",
		MainResolver:   resolver,
		NoArgsMainFunc: defaultNoArgsMainFunc,
	}
}

func (cmd *DashCmd) run(name string, arg ...string) error {
	prefix := cmd.Name + cmd.Dash
	if strings.HasPrefix(name, prefix) {
		noPrefixName := strings.TrimPrefix(name, prefix)
		main, err := cmd.MainResolver.ResolveMain(noPrefixName)
		if err == nil && main != nil {
			main()
			return nil
		}
		lpResolver := NewLookPathMainResolver()
		main, err = lpResolver.ResolveMain(name)
		if err == nil && main != nil {
			main()
			return nil
		}
		return fmt.Errorf("not found: %s", noPrefixName)
	}
	if len(arg) == 0 {
		cmd.NoArgsMainFunc()
		return nil
	}
	newArgs := make([]string, len(arg))
	newArgs[0] = prefix + arg[0]
	for i := 1; i < len(arg); i++ {
		newArgs[i] = arg[i]
	}
	os.Args = newArgs
	cmd.Main()
	return nil
}

func (cmd *DashCmd) Main() {
	if len(os.Args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "invalid state: $0 is not specified")
		os.Exit(1)
	}
	name := filepath.Base(os.Args[0])
	err := cmd.run(name, os.Args[1:]...)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(127)
	}
}

func (cmd *DashCmd) GenerateMain() func() {
	return func() { cmd.Main() }
}

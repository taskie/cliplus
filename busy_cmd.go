package cliplus

import (
	"fmt"
	"os"
	"path/filepath"
)

var defaultNoArgsMainFunc = func() {
	_, _ = fmt.Fprintln(os.Stderr, "please specify subcommand")
	os.Exit(1)
}

type BusyCmd struct {
	Name           string
	MainResolver   MainResolver
	NoArgsMainFunc func()
}

func NewBusyCmd(name string, resolver MainResolver) *BusyCmd {
	return &BusyCmd{
		Name:           name,
		MainResolver:   resolver,
		NoArgsMainFunc: defaultNoArgsMainFunc,
	}
}

func (cmd *BusyCmd) run(name string, arg ...string) error {
	if name == cmd.Name {
		if len(arg) == 0 {
			cmd.NoArgsMainFunc()
			return nil
		}
		os.Args = arg
		cmd.Main()
		return nil
	}
	main, err := cmd.MainResolver.ResolveMain(name)
	if err != nil {
		return err
	}
	if main == nil {
		return fmt.Errorf("not found: %s", name)
	}
	main()
	return nil
}

func (cmd *BusyCmd) Main() {
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

func (cmd *BusyCmd) GenerateMain() func() {
	return func() { cmd.Main() }
}

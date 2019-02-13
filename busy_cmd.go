package cliplus

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type BusyCmd struct {
	ResolverMap map[string]func()
	Resolver    func(name string) func()
}

func NewBusyCmd(selfName string) *BusyCmd {
	resolverMap := make(map[string]func())
	busyCmd := &BusyCmd{
		ResolverMap: resolverMap,
	}
	resolverMap[selfName] = func() { busyCmd.Main() }
	return busyCmd
}

func (busyCmd *BusyCmd) Resolve(name string) func() {
	if busyCmd.Resolver != nil {
		return busyCmd.Resolver(name)
	}
	return busyCmd.ResolverMap[name]
}

func (busyCmd *BusyCmd) Run(name string) error {
	cmd := busyCmd.Resolve(name)
	if cmd == nil {
		return fmt.Errorf("command not found: %s", name)
	}
	cmd()
	return nil
}

func (busyCmd *BusyCmd) Main() {
	name := filepath.Base(os.Args[0])
	os.Args = os.Args[1:]
	if len(os.Args) == 0 {
		log.Println("please specify subcommand")
		os.Exit(1)
	}
	err := busyCmd.Run(name)
	if err != nil {
		log.Println(err)
		os.Exit(127)
	}
}

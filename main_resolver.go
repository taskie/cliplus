package cliplus

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type MainResolver interface {
	ResolveMain(name string) (func(), error)
}

type MapMainResolver struct {
	Map map[string]func()
}

func NewMapMainResolver(resolverMap map[string]func()) *MapMainResolver {
	return &MapMainResolver{
		Map: resolverMap,
	}
}

func (resolver *MapMainResolver) ResolveMain(name string) (func(), error) {
	if v, ok := resolver.Map[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("not found: %s", name)
}

type LookPathMainResolver struct{}

func NewLookPathMainResolver() *LookPathMainResolver {
	return new(LookPathMainResolver)
}

func runCommand(name string, arg ...string) {
	fmt.Printf("%s %v\n", name, arg)
	cmd := exec.Command(name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				os.Exit(status.ExitStatus())
			}
		} else {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	}
}

func (resolver *LookPathMainResolver) ResolveMain(name string) (func(), error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return nil, err
	}
	return func() {
		runCommand(path, os.Args[1:]...)
	}, nil
}

type MultiMainResolver struct {
	MainResolvers []MainResolver
}

func NewMultiMainResolver(resolvers []MainResolver) *MultiMainResolver {
	return &MultiMainResolver{
		MainResolvers: resolvers,
	}
}

func (resolver *MultiMainResolver) ResolveMain(name string) (func(), error) {
	for _, res := range resolver.MainResolvers {
		main, err := res.ResolveMain(name)
		if err != nil {
			return main, nil
		}
	}
	return nil, fmt.Errorf("not found: %s", name)
}

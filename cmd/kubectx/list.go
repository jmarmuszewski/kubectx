package main

import (
	"fmt"
	"io"

	"facette.io/natsort"
	"github.com/fatih/color"

	"github.com/ahmetb/kubectx/cmd/kubectx/kubeconfig"
)

type context struct {
	Name string `yaml:"name"`
}

type kubeconfigContents struct {
	APIVersion     string    `yaml:"apiVersion"`
	CurrentContext string    `yaml:"current-context"`
	Contexts       []context `yaml:"contexts"`
}

// ListOp describes listing contexts.
type ListOp struct{}

func (_ ListOp) Run(stdout, _ io.Writer) error {
	kc := new(kubeconfig.Kubeconfig).WithLoader(defaultLoader)
	defer kc.Close()
	rootNode, err := kc.ParseRaw()
	if err != nil {
		return err
	}

	ctxs := kubeconfig.ContextNames(rootNode)
	natsort.Sort(ctxs)

	// TODO support KUBECTX_CURRENT_FGCOLOR
	// TODO support KUBECTX_CURRENT_BGCOLOR
	cur :=  kubeconfig.GetCurrentContext(rootNode)
	for _, c := range ctxs {
		s := c
		if c == cur {
			s = color.New(color.FgGreen, color.Bold).Sprint(c)
		}
		fmt.Fprintf(stdout, "%s\n", s)
	}
	return nil
}

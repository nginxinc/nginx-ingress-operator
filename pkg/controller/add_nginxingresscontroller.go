package controller

import (
	"github.com/nginxinc/nginx-ingress-operator/pkg/controller/nginxingresscontroller"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, nginxingresscontroller.Add)
}

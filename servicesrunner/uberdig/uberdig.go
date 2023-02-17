package uberdig

import (
	"go.uber.org/dig"
)

type (
	ContainerCall func() (*dig.Container, error)
	Invoke        func(container *dig.Container) error
	InvokeError   func(container *dig.Container, err error)
)

func Run(containerCall ContainerCall, invoke Invoke, onError InvokeError) error {
	container, err := containerCall()
	if err != nil {
		return err
	}
	if err := invoke(container); err != nil {
		onError(container, err)
	}
	return nil
}

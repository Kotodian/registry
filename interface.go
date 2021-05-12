package registry

import "errors"

var (
	ErrNodesNil = errors.New("service doesn't have any node")
)

type Service interface {
	ID() string
	ListNodes() []Node
	GetNode(id string) Node
	AddNode(Node) error
	DeleteNode(Node) error
}

type Node interface {
	ID() string
}

type Registry interface {
	ListServices() []Service
	GetService(id string) (Service, error)
	AddService(service Service) error
	NotifyService()
	DeleteService(service Service) error
}

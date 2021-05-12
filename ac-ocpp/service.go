package ac_ocpp

import (
	"errors"
	"registry"
)

type AcOCPP struct {
	hostname string
	nodes    []*ChargeStation
}

func (a *AcOCPP) ID() string {
	return a.hostname
}

func (a *AcOCPP) ListNodes() []registry.Node {
	if len(a.nodes) == 0 {
		return nil
	}
	nodes := make([]registry.Node, 0)
	for _, node := range a.nodes {
		nodes = append(nodes, node)
	}
	return nodes
}

func (a *AcOCPP) GetNode(id string) registry.Node {
	if len(a.nodes) == 0 {
		return nil
	}
	for _, node := range a.nodes {
		if node.sn == id {
			return node
		}
	}
	return nil
}

func (a *AcOCPP) AddNode(node registry.Node) error {
	chargeStation, ok := node.(*ChargeStation)
	if !ok {
		return errors.New("this node is not a charge station")
	}

	a.nodes = append(a.nodes, chargeStation)
	return nil
}

func (a *AcOCPP) DeleteNode(node registry.Node) error {
	_, ok := node.(*ChargeStation)
	if !ok {
		return errors.New("this node is not a charge station")
	}
	if len(a.nodes) == 0 {
		return registry.ErrNodesNil
	}
	nodes := make([]*ChargeStation, 0)
	for _, station := range a.nodes {
		if station.sn == node.ID() {
			continue
		}
		nodes = append(nodes, station)
	}
	a.nodes = nodes
	return nil
}

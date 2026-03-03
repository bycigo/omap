package omap

import (
	"errors"

	"go.yaml.in/yaml/v3"
)

// MarshalYAML implements the yaml.Marshaler interface for Map.
func (m *Map[K, V]) MarshalYAML() (any, error) {
	kvNodes := make([]*yaml.Node, 0, len(m.kv)*2)
	for k, v := range m.All() {
		keyNode := &yaml.Node{}
		if err := keyNode.Encode(k); err != nil {
			return nil, err
		}
		valueNode := &yaml.Node{}
		if err := valueNode.Encode(v); err != nil {
			return nil, err
		}
		kvNodes = append(kvNodes, keyNode, valueNode)
	}

	mapNode := &yaml.Node{
		Kind:    yaml.MappingNode,
		Tag:     "!!map",
		Content: kvNodes,
	}

	return mapNode, nil
}

// UnmarshalYAML implements the yaml.Unmarshaler interface for Map.
func (m *Map[K, V]) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.MappingNode {
		return errors.New("expected a mapping node")
	}

	if len(n.Content)%2 != 0 {
		return errors.New("mapping node has odd number of content nodes")
	}

	for i := 0; i < len(n.Content); i += 2 {
		var key K
		if err := n.Content[i].Decode(&key); err != nil {
			return err
		}
		var value V
		if err := n.Content[i+1].Decode(&value); err != nil {
			return err
		}
		m.Set(key, value)
	}

	return nil
}

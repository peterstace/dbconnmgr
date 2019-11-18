package main

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

type ConnDetails struct {
	Username string
	Password string
	Endpoint string
	Database string
	Labels   map[string]string
}

type Node struct {
	ConnDetails
	SumChildren     []Node
	ProductChildren []Node
}

func ParseTree(r io.Reader) (Node, error) {
	var m map[interface{}]interface{}
	if err := yaml.NewDecoder(r).Decode(&m); err != nil {
		return Node{}, fmt.Errorf("decoding yaml: %v", err)
	}
	return newNode(m)
}

func newNode(m map[interface{}]interface{}) (Node, error) {
	var n Node
	for k, v := range m {
		_, ok := k.(string)
		if !ok {
			return Node{}, fmt.Errorf("key with non-string type: %v", k)
		}
		k := k.(string)

		vAsString, vIsString := v.(string)
		vAsChildren, vIsChildren := v.([]interface{})

		switch k {
		case "username":
			if !vIsString {
				return Node{}, fmt.Errorf("username is non-string: %T", v)
			}
			n.Username = vAsString
		case "password":
			if !vIsString {
				return Node{}, fmt.Errorf("password is non-string: %T", v)
			}
			n.Password = vAsString
		case "endpoint":
			if !vIsString {
				return Node{}, fmt.Errorf("endpoint is non-string: %T", v)
			}
			n.Endpoint = vAsString
		case "database":
			if !vIsString {
				return Node{}, fmt.Errorf("database is non-string: %T", v)
			}
			n.Database = vAsString
		case "sum", "product":
			if !vIsChildren {
				return Node{}, fmt.Errorf("sum or product is non-array: %T", v)
			}
			if len(n.SumChildren) != 0 || len(n.ProductChildren) != 0 {
				return Node{}, errors.New("multiple sets of children defined")
			}
			for _, c := range vAsChildren {
				childAsMap, childIsMap := c.(map[interface{}]interface{})
				if !childIsMap {
					return Node{}, fmt.Errorf("child is non-map: %T", c)
				}
				child, err := newNode(childAsMap)
				if err != nil {
					return Node{}, fmt.Errorf("could not build sub-tree: %v", err)
				}
				switch k {
				case "sum":
					n.SumChildren = append(n.SumChildren, child)
				case "product":
					n.ProductChildren = append(n.ProductChildren, child)
				}
			}
		default:
			if !vIsString {
				return Node{}, fmt.Errorf("%s label is non-string: %T", k, v)
			}
			if n.Labels == nil {
				n.Labels = make(map[string]string)
			}
			n.Labels[k] = vAsString
		}
	}
	return n, nil
}

// TODO: no error?
func (n Node) Flatten() ([]ConnDetails, error) {
	return []ConnDetails{n.ConnDetails}, nil
}

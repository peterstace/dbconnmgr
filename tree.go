package main

import (
	"errors"
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

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

func (n Node) Flatten() (ConnList, error) {
	switch {
	case len(n.ProductChildren) > 0 && len(n.SumChildren) > 0:
		return nil, errors.New("multiple sets of children defined")
	case len(n.ProductChildren) == 0 && len(n.SumChildren) == 0:
		return []ConnDetails{n.ConnDetails}, nil
	case len(n.ProductChildren) > 0:
		var multiplicands [][]ConnDetails
		for _, child := range n.ProductChildren {
			childDetails, err := child.Flatten()
			if err != nil {
				return nil, fmt.Errorf("could not flatten child: %v", err)
			}
			multiplicands = append(multiplicands, childDetails)
		}
		product := []ConnDetails{n.ConnDetails}
		for _, multiplicand := range multiplicands {
			var newProduct []ConnDetails
			for i := range product {
				for j := range multiplicand {
					merged, err := Merge(product[i], multiplicand[j])
					if err != nil {
						return nil, fmt.Errorf("could not merge children: %v", err)
					}
					newProduct = append(newProduct, merged)
				}
			}
			product = newProduct
		}
		return product, nil
	case len(n.SumChildren) > 0:
		var allDetails []ConnDetails
		for _, child := range n.SumChildren {
			childDetails, err := child.Flatten()
			if err != nil {
				return nil, fmt.Errorf("could not flatten child: %v", err)
			}
			allDetails = append(allDetails, childDetails...)
		}
		for i := range allDetails {
			merged, err := Merge(allDetails[i], n.ConnDetails)
			if err != nil {
				return nil, fmt.Errorf("could not merge connection details: %v", err)
			}
			allDetails[i] = merged
		}
		return allDetails, nil
	default:
		panic(false)
	}
}

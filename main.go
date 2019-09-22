package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

func main() {
	var m map[interface{}]interface{}
	if err := yaml.NewDecoder(os.Stdin).Decode(&m); err != nil {
		log.Fatal(err)
	}

	results, err := depthFirstWalk(m, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	choices, err := buildChoices(results)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	for i, choice := range choices {
		fmt.Printf("%d: %s\n", i, choice)
	}
}

type keyval struct {
	key, val string
}

func depthFirstWalk(node map[interface{}]interface{}, kvs []keyval) ([][]keyval, error) {
	var children []interface{}
	for key, val := range node {
		if key == "children" {
			var ok bool
			children, ok = val.([]interface{})
			if !ok {
				return nil, errors.New("found children key not of type []inteface{}")
			}
		} else {
			keyStr, ok := key.(string)
			if !ok {
				return nil, errors.New("found key not of type string")
			}
			valStr, ok := val.(string)
			if !ok {
				return nil, errors.New("found value not of type string")
			}
			kvs = append(kvs, keyval{keyStr, valStr})
		}
	}

	if len(children) == 0 {
		kvsCopy := make([]keyval, len(kvs))
		copy(kvsCopy, kvs)
		return [][]keyval{kvsCopy}, nil
	}

	var results [][]keyval
	for _, child := range children {
		childMap, ok := child.(map[interface{}]interface{})
		if !ok {
			return nil, fmt.Errorf("found child not of type map[string]interface{}: %T", child)
		}
		result, err := depthFirstWalk(childMap, kvs)
		if err != nil {
			return nil, err
		}
		for _, subResult := range result {
			results = append(results, subResult)
		}
	}
	return results, nil
}

type choice struct {
	connStr string
	tags    map[string]string
}

func buildChoices(choicesKeyVals [][]keyval) ([]choice, error) {
	choices := make([]choice, 0, len(choicesKeyVals))
	for _, choiceKeyVals := range choicesKeyVals {
		m := make(map[string]string)
		for _, kv := range choiceKeyVals {
			m[kv.key] = kv.val

		}

		username, ok := m["username"]
		if !ok {
			return nil, errors.New("username missing")
		}
		password, ok := m["password"]
		if !ok {
			return nil, errors.New("password missing")
		}
		endpoint, ok := m["endpoint"]
		if !ok {
			return nil, errors.New("endpoint missing")
		}
		database, ok := m["database"]
		if !ok {
			return nil, errors.New("database missing")
		}

		delete(m, "username")
		delete(m, "password")
		delete(m, "endpoint")
		delete(m, "database")

		choices = append(choices, choice{
			connStr: fmt.Sprintf("postgres://%s:%s@%s/%s", username, password, endpoint, database),
			tags:    m,
		})
	}
	return choices, nil
}

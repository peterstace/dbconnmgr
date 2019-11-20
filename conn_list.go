package main

import "fmt"

type ConnList []ConnDetails

func (c ConnList) Selections() []string {
	kvList := make([]map[string]string, len(c))
	for i, details := range c {
		kvList[i] = map[string]string{
			"uname": details.Username,
			"db":    details.Database,
		}
		for k, v := range details.Labels {
			kvList[i][k] = v
		}
	}

	out := make([]string, len(c))
	for i := range c {
		out[i] = fmt.Sprintf("%v", kvList[i])
	}
	return out
}

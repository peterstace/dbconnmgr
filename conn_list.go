package main

import "fmt"

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

type ConnList []ConnDetails

func (c ConnList) Selections() []string {
	kvList := make([]map[string]string, len(c))
	keyToMaxWidth := make(map[string]int)
	for i, details := range c {
		kvList[i] = map[string]string{
			"uname": details.Username,
			"db":    details.Database,
		}
		keyToMaxWidth["uname"] = max(keyToMaxWidth["uname"], len(details.Username))
		keyToMaxWidth["db"] = max(keyToMaxWidth["db"], len(details.Database))
		for k, v := range details.Labels {
			kvList[i][k] = v
			keyToMaxWidth[k] = max(keyToMaxWidth[k], len(v))
		}
	}

	type keywidth struct {
		key   string
		width int
	}
	keys := make([]keywidth, 0, len(keyToMaxWidth))
	for key, width := range keyToMaxWidth {
		keys = append(keys, keywidth{key, width})
	}

	// TODO: sort the keys array

	var format string
	for _, kw := range keys {
		format += fmt.Sprintf("%s=%%-%ds ", kw.key, kw.width)
	}

	out := make([]string, len(c))
	for i := range c {
		args := make([]interface{}, len(keys))
		for j, kw := range keys {
			args[j] = kvList[i][kw.key]
		}
		out[i] = fmt.Sprintf(format, args...)
	}
	return out
}

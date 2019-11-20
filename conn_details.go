package main

import "fmt"

type ConnDetails struct {
	Username string
	Password string
	Endpoint string
	Database string
	Labels   map[string]string
}

func (d ConnDetails) ConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s",
		d.Username,
		d.Password,
		d.Endpoint,
		d.Database,
	)
}

func Merge(x, y ConnDetails) (ConnDetails, error) {
	var m ConnDetails

	mergeStr := func(s1, s2 string) (string, error) {
		if s1 != "" && s2 != "" {
			return "", fmt.Errorf("duplicated: %s and %s", s1, s2)
		}
		if s1 == "" {
			return s2, nil
		}
		return s1, nil
	}

	var err error
	m.Username, err = mergeStr(x.Username, y.Username)
	if err != nil {
		return m, fmt.Errorf("merging usernames: %v", err)
	}
	m.Password, err = mergeStr(x.Password, y.Password)
	if err != nil {
		return m, fmt.Errorf("merging passwords: %v", err)
	}
	m.Endpoint, err = mergeStr(x.Endpoint, y.Endpoint)
	if err != nil {
		return m, fmt.Errorf("merging endpoints: %v", err)
	}
	m.Database, err = mergeStr(x.Database, y.Database)
	if err != nil {
		return m, fmt.Errorf("merging databases: %v", err)
	}
	for _, labels := range []struct {
		setA, setB map[string]string
	}{
		{x.Labels, y.Labels},
		{y.Labels, x.Labels},
	} {
		for k, v := range labels.setA {
			if m.Labels == nil {
				m.Labels = make(map[string]string)
			}
			m.Labels[k], err = mergeStr(v, labels.setB[k])
			if err != nil {
				return m, fmt.Errorf("merging label %s: %v", k, err)
			}
		}
	}
	return m, nil
}

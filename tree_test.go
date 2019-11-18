package main

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestParseTree(t *testing.T) {
	for i, tt := range []struct {
		input  string
		output Node
	}{
		{
			input: `
				username: u
				password: p
				endpoint: e
				database: d
				label1: foo
				label2: bar
				sum:
					- username: A
					- username: B
				`,
			output: Node{
				ConnDetails: ConnDetails{
					Username: "u",
					Password: "p",
					Endpoint: "e",
					Database: "d",
					Labels: map[string]string{
						"label1": "foo",
						"label2": "bar",
					},
				},
				SumChildren: []Node{
					Node{ConnDetails: ConnDetails{Username: "A"}},
					Node{ConnDetails: ConnDetails{Username: "B"}},
				},
			},
		},
		{
			input: `
				product:
					- username: A
					- username: B
				`,
			output: Node{
				ProductChildren: []Node{
					Node{ConnDetails: ConnDetails{Username: "A"}},
					Node{ConnDetails: ConnDetails{Username: "B"}},
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tt.input = strings.ReplaceAll(tt.input, "\t", " ")
			got, err := ParseTree(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("\ngot:  %v\nwant: %v\n", got, tt.output)
			}
		})
	}
}

func TestIllegalTree(t *testing.T) {
	_, err := ParseTree(strings.NewReader(strings.ReplaceAll(`
		username: u
		password: p
		sum:
			- username: u1
			  password: p1
		product:
			- username: u2
			  password: p2`, "\t", " ")))
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestFlatten(t *testing.T) {
	for i, tt := range []struct {
		input  string
		output []ConnDetails
	}{
		{
			input: `
				username: user
				password: pass
				endpoint: end
				database: db
				label1: foo
				label2: bar`,
			output: []ConnDetails{
				{
					Username: "user",
					Password: "pass",
					Endpoint: "end",
					Database: "db",
					Labels: map[string]string{
						"label1": "foo",
						"label2": "bar",
					},
				},
			},
		},
		{
			input: `
				endpoint: end
				sum:
					- username: user1
					  password: pass1
					  sum:
					  	- db: db1
					  	- db: db2
					- username: user2
					  password: pass2
					  sum:
					  	- db: db1
					  	- db: db2`,
			output: []ConnDetails{
				{
					Username: "user1",
					Password: "pass1",
					Endpoint: "end",
					Database: "db1",
				},
				{
					Username: "user1",
					Password: "pass1",
					Endpoint: "end",
					Database: "db2",
				},
				{
					Username: "user2",
					Password: "pass2",
					Endpoint: "end",
					Database: "db1",
				},
				{
					Username: "user2",
					Password: "pass2",
					Endpoint: "end",
					Database: "db2",
				},
			},
		},
		{
			input: `
				endpoint: end
				product:
				  - sum:
				    - username: user1
				      password: pass1
				    - username: user2
				      password: pass2
				  - sum:
				    - db: db1
				    - db: db2`,
			output: []ConnDetails{
				{
					Username: "user1",
					Password: "pass1",
					Endpoint: "end",
					Database: "db1",
				},
				{
					Username: "user1",
					Password: "pass1",
					Endpoint: "end",
					Database: "db2",
				},
				{
					Username: "user2",
					Password: "pass2",
					Endpoint: "end",
					Database: "db1",
				},
				{
					Username: "user2",
					Password: "pass2",
					Endpoint: "end",
					Database: "db2",
				},
			},
		},
	} {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			tt.input = strings.ReplaceAll(tt.input, "\t", " ")
			tree, err := ParseTree(strings.NewReader(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			got, err := tree.Flatten()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tt.output) {
				t.Errorf("\ngot:  %v\nwant: %v\n", got, tt.output)
			}
		})
	}
}
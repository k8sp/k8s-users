package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/topicai/candy"
)

//Policy include muliple rules
type Policy struct {
	Rules []Rule
}

//Rule {"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": #Spec}
type Rule struct {
	ApiVersion string `json:"apiVersion,omitempty:"`
	Kind       string `json:"kind,omitempty"`
	Spec       Spec   `json:"spec,omitempty"`
}

//Spec  {"user":<username>, "namespace": <namespace>, "resource": "*", "apiGroup": "*", "readonly": false }}
type Spec struct {
	// subject matching
	User  string `json:"user,omitempty"`
	Group string `json:"group,omitempty"`
	//resource-matching
	ApiGroup  string `json:"apiGroup,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Resource  string `json:"resource,omitempty"`
	// non-resource-mathing
	NonResourcePath string `json:"nonResourcePath,omitempty"`
	Readonly        bool   `json:"readonly,omitempty"`
}

//Exists whether user has already in polic file
func (p *Policy) Exists(u Users) bool {
	for _, r := range p.Rules {
		if r.Spec.User == u.Username {
			return true
		}
	}
	return false
}

//Update policies using param: usersBody
func (p *Policy) Update(u Users) {
	for i, r := range p.Rules {
		// user
		if r.Spec.User == u.Username {
			p.Rules[i] = newRule(u.Username, u.Namespace)
			// service account
		} else if r.Spec.User == getDefaultServiceAcccount(u.Namespace) {
			p.Rules[i] = newRule(getDefaultServiceAcccount(u.Namespace), u.Namespace)
		} else {
			continue
		}
	}
}

// Append new policy at the end of policies
func (p *Policy) Append(u Users) {
	p.Rules = append(p.Rules, newRule(u.Username, u.Namespace))
}

//LoadPoliciesfromJSONFile init policy struct with json file
func LoadPoliciesfromJSONFile(filename string) (Policy, error) {
	p := Policy{}
	input, e := ioutil.ReadFile(filename)
	candy.Must(e)

	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.EqualFold(strings.TrimSpace(line), "") {
			continue
		}
		var r Rule
		err := json.Unmarshal([]byte(line), &r)
		fmt.Println(r)
		candy.Must(err)
		p.Rules = append(p.Rules, r)
	}
	return p, nil
}

//DumpJSONFile write policy to a file formated json
func (p *Policy) DumpJSONFile(filename string) error {
	lines := []string{}
	for _, v := range p.Rules {
		b, err := json.Marshal(v)
		candy.Must(err)
		//fmt.Println(string(b))
		lines = append(lines, string(b))
	}
	e := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	return e
}

func newRule(username, namespace string) Rule {
	spec := Spec{
		User:            username,
		Group:           "system:authenticated",
		ApiGroup:        "*",
		Namespace:       namespace,
		Resource:        "*",
		NonResourcePath: "*",
		Readonly:        false,
	}
	r := Rule{
		ApiVersion: "abac.authorization.kubernetes.io/v1beta1",
		Kind:       "Policy",
		Spec:       spec,
	}
	return r
}

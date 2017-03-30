package users

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/topicai/candy"
)

type Policy struct {
	rules []Rule
}

//Rule {"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": #Spec}
type Rule struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       Spec   `json:"spec"`
}

//Spec  {"user":<username>, "namespace": <namespace>, "resource": "*", "apiGroup": "*", "readonly": false }}
type Spec struct {
	User            string `json:"user"`
	Namespace       string `json:"namespace"`
	Resource        string `json:"resource"`
	APIGroup        string `json:"apiGroup"`
	Readonly        bool   `json:"readonly"`
	NonResourcePath string `json:"nonResourcePath"`
}

func (p *Policy) Exists(user UsersBody) bool {
	for _, r := range p.rules {
		if r.Spec.User == user.Username {
			return true
		}
	}
	return false
}

//Update policies using param: usersBody
func (p *Policy) Update(user UsersBody) {
	for i, r := range p.rules {
		if r.Spec.User == user.Username {
			p.rules[i] = makeRule(user.Username, user.Namespace)
		} else if r.Spec.User == getDefaultServiceAcccount(user.Username) {
			p.rules[i] = makeRule(getDefaultServiceAcccount(user.Username), user.Namespace)
		} else {
			continue
		}
	}
}

// Append new policy at the end of policies
func (p *Policy) Append(user UsersBody) {
	p.rules = append(p.rules, makeRule(user.Username, user.Namespace))
	p.rules = append(p.rules, makeRule(user.Username, user.Namespace))
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
		candy.Must(err)
		p.rules = append(p.rules, r)
	}
	return p, nil
}

func (p *Policy) ToJsonFile(filename string) error {
	lines := []string{}
	for _, p := range p.rules {
		b, err := json.Marshal(p)
		candy.Must(err)
		lines = append(lines, string(b))
	}
	e := ioutil.WriteFile(filename, []byte(strings.Join(lines, "\n")), 0644)
	return e
}

func makeRule(username, namespace string) Rule {
	spec := Spec{
		User:            username,
		Namespace:       namespace,
		Resource:        "*",
		APIGroup:        "*",
		Readonly:        false,
		NonResourcePath: "*",
	}
	r := Rule{
		APIVersion: "abac.authorization.kubernetes.io/v1beta1",
		Kind:       "Policy",
		Spec:       spec,
	}
	return r
}

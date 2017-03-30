package users

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/topicai/candy"
)

//Policy include muliple rules
type Policy struct {
	rules []Rule
}

//Rule {"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": #Spec}
type Rule struct {
	apiVersion string
	kind       string
	spec       Spec
}

//Spec  {"user":<username>, "namespace": <namespace>, "resource": "*", "apiGroup": "*", "readonly": false }}
type Spec struct {
	user            string
	namespace       string
	resource        string
	apiGroup        string
	readonly        bool
	nonResourcePath string
}

func (p *Policy) Exists(user UsersBody) bool {
	for _, r := range p.rules {
		if r.spec.user == user.Username {
			return true
		}
	}
	return false
}

//Update policies using param: usersBody
func (p *Policy) Update(user UsersBody) {
	for i, r := range p.rules {
		if r.spec.user == user.Username {
			p.rules[i] = makeRule(user.Username, user.Namespace)
		} else if r.spec.user == getDefaultServiceAcccount(user.Namespace) {
			p.rules[i] = makeRule(getDefaultServiceAcccount(user.Namespace), user.Namespace)
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
		user:            username,
		namespace:       namespace,
		resource:        "*",
		apiGroup:        "*",
		readonly:        false,
		nonResourcePath: "*",
	}
	r := Rule{
		apiVersion: "abac.authorization.kubernetes.io/v1beta1",
		kind:       "Policy",
		spec:       spec,
	}
	return r
}

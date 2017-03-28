package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/topicai/candy"
)

//Policy : {"apiVersion": "abac.authorization.kubernetes.io/v1beta1", "kind": "Policy", "spec": #Spec}
type Policy struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Spec       Spec   `json:"spec"`
}

//Spec : {"user":<username>, "namespace": <namespace>, "resource": "*", "apiGroup": "*", "readonly": false }}
type Spec struct {
	User            string `json:"user"`
	Namespace       string `json:"namespace"`
	Resource        string `json:"resource"`
	APIGroup        string `json:"apiGroup"`
	Readonly        bool   `json:"readonly"`
	NonResourcePath string `json:"nonResourcePath"`
}

func getDefaultServiceAcccount(username string) string {
	return "system:serviceaccount:" + username + ":default"
}

func makePolicy(username, namespace string) Policy {
	spec := Spec{
		User:            username,
		Namespace:       namespace,
		Resource:        "*",
		APIGroup:        "*",
		Readonly:        false,
		NonResourcePath: "*",
	}
	policy := Policy{
		APIVersion: "abac.authorization.kubernetes.io/v1beta1",
		Kind:       "Policy",
		Spec:       spec,
	}
	return policy
}

func newPolicy(username, namespace string, policies []Policy) []Policy {
	newPolicies := []Policy{}
	for _, p := range policies {
		// ignore the policy rule if username already exists.
		if p.Spec.User == username || p.Spec.User == getDefaultServiceAcccount(username) {
			continue
		}
		newPolicies = append(newPolicies, p)
	}
	// append new policy
	newPolicies = append(newPolicies, makePolicy(username, namespace))
	newPolicies = append(newPolicies, makePolicy(getDefaultServiceAcccount(username), namespace))
	return newPolicies
}

// UpdatePolicyFile append user policy rule to the old file, if the user has
// already exists, replace it.
func UpdatePolicyFile(username, namespace, policyFile string) error {
	// TODO: Add file lock
	input, e := ioutil.ReadFile(policyFile)
	candy.Must(e)
	policies := []Policy{}
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		if strings.EqualFold(strings.TrimSpace(line), "") {
			continue
		}
		fmt.Println("iii" + line)
		var p Policy
		err := json.Unmarshal([]byte(line), &p)
		candy.Must(err)
		policies = append(policies, p)
	}
	fmt.Println("----")
	policies = newPolicy(username, namespace, policies)
	newLines := []string{}
	for _, p := range policies {
		b, err := json.Marshal(p)
		candy.Must(err)
		newLines = append(newLines, string(b))
	}
	e = ioutil.WriteFile(policyFile, []byte(strings.Join(newLines, "\n")), 0644)
	return e
}

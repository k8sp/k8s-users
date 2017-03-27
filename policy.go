package main

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/nightlyone/lockfile"
	"github.com/topicai/candy"
)

type Policy struct {
	apiVersion string
	kind       string
	spec       Spec
}

type Spec struct {
	user      string
	namespace string
	resource  string
	apiGroup  string
	readonly  bool
}

func getDefaultServiceAcccount(username string) string {
	return "system:serviceaccount:" + username + ":default"
}

func getUserPolicy(username, namespace string) Policy {
	spec := Spec{
		user:      username,
		namespace: namespace,
		resource:  "*",
		apiGroup:  "*",
		readonly:  false,
	}
	policy := Policy{
		apiVersion: "abac.authorization.kubernetes.io/v1beta1",
		kind:       "Policy",
		spec:       spec,
	}
	return policy
}

func makePolicy(username, namespace string) Policy {
	spec := Spec{
		user:      username,
		namespace: namespace,
		resource:  "*",
		apiGroup:  "*",
		readonly:  false,
	}
	policy := Policy{
		apiVersion: "abac.authorization.kubernetes.io/v1beta1",
		kind:       "Policy",
		spec:       spec,
	}
	return policy
}

func newPolicy(username, namespace string, policies []Policy) []Policy {
	newPolicies := []Policy{}
	for _, p := range policies {
		// ignore the policy rule if username already exists.
		if p.spec.user == username || p.spec.user == getDefaultServiceAcccount(username) {
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
	// file lock
	lock, e := lockfile.New(policyFile)
	candy.Must(e)

	e = lock.TryLock()
	candy.Must(e)

	defer lock.Unlock()

	input, e := ioutil.ReadFile(policyFile)
	candy.Must(e)

	policies := []Policy{}
	lines := strings.Split(string(input), "\n")
	for _, line := range lines {
		var p Policy
		err := json.Unmarshal([]byte(line), &p)
		candy.Must(err)
		policies = append(policies, p)
	}

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

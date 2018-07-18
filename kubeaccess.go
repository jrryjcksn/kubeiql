// Copyright 2018 Yipee.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

// Functions for retrieving Kubernetes information from a cluster

// Get a single resource instance from a namespace
func getK8sResource(kind, namespace, name string) map[string]interface{} {
	return fromJson(
		lookUpResource(kind, namespace, name)).(map[string]interface{})
}

func fromJson(val []byte) interface{} {
	var result interface{}

	if err := json.Unmarshal(val, &result); err != nil {
		panic(err)
	}

	return result
}

func lookUpResource(kind, namespace, name string) []byte {
	cmd := exec.Command("/usr/local/bin/kubectl", "get",
		"-o", "json", "--namespace", namespace, kind, name)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return bytes
}

// Get all resource instances of a specific kind
func getAllK8sObjsOfKind(
	ctx context.Context,
	kind string,
	test func(map[string]interface{}) bool) []resource {
	cmd := exec.Command("/usr/local/bin/kubectl", "get",
		"-o", "json", "--all-namespaces", kind)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	var results []resource
	arr := (fromJson(bytes).(map[string]interface{}))["items"].([]interface{})
	for _, res := range arr {
		if test(res.(map[string]interface{})) {
			results =
				append(results, mapToResource(ctx, res.(map[string]interface{})))
		}
	}
	if results == nil {
		results = make([]resource, 0)
	}
	return results
}

// Get all resource instances of a specific kind in a specific namespace
func getAllK8sObjsOfKindInNamespace(
	ctx context.Context,
	kind, ns string,
	test func(map[string]interface{}) bool) []resource {
	cmd := exec.Command("/usr/local/bin/kubectl", "get",
		"-o", "json", "--namespace", ns, kind)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	var results []resource
	arr := (fromJson(bytes).(map[string]interface{}))["items"].([]interface{})
	for _, res := range arr {
		if test(res.(map[string]interface{})) {
			results =
				append(results, mapToResource(ctx, res.(map[string]interface{})))
		}
	}
	if results == nil {
		results = make([]resource, 0)
	}
	return results
}

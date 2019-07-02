// Copyright 2019 Drone IO, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package yaml

import filepath "github.com/bmatcuk/doublestar"

// Conditions defines a group of conditions.
type Conditions struct {
	Action   Condition    `json:"action,omitempty"`
	Cron     Condition    `json:"cron,omitempty"`
	Ref      Condition    `json:"ref,omitempty"`
	Repo     Condition    `json:"repo,omitempty"`
	Instance Condition    `json:"instance,omitempty"`
	Target   Condition    `json:"target,omitempty"`
	Event    Condition    `json:"event,omitempty"`
	Branch   Condition    `json:"branch,omitempty"`
	Status   Condition    `json:"status,omitempty"`
	Paths    Condition    `json:"paths,omitempty"`
	Matrix   ConditionMap `json:"matrix,omitempty"`
}

// Condition defines a runtime condition.
type Condition struct {
	Include []string `yaml:"include,omitempty" json:"include,omitempty"`
	Exclude []string `yaml:"exclude,omitempty" json:"exclude,omitempty"`
}

type ConditionMap struct {
	Include map[string]string
	Exclude map[string]string
}

// Match returns true if the string matches the include
// patterns and does not match any of the exclude patterns.
func (c *Condition) Match(v string) bool {
	if c.Excludes(v) {
		return false
	}
	if c.Includes(v) {
		return true
	}
	if len(c.Include) == 0 {
		return true
	}
	return false
}

// Includes returns true if the string matches the include
// patterns.
func (c *Condition) Includes(v string) bool {
	for _, pattern := range c.Include {
		if ok, _ := filepath.Match(pattern, v); ok {
			return true
		}
	}
	return false
}

// Excludes returns true if the string matches the exclude
// patterns.
func (c *Condition) Excludes(v string) bool {
	for _, pattern := range c.Exclude {
		if ok, _ := filepath.Match(pattern, v); ok {
			return true
		}
	}
	return false
}

// UnmarshalYAML implements yml unmarshalling.
func (c *Condition) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var out1 string
	var out2 []string
	var out3 = struct {
		Include []string
		Exclude []string
	}{}

	err := unmarshal(&out1)
	if err == nil {
		c.Include = []string{out1}
		return nil
	}

	unmarshal(&out2)
	unmarshal(&out3)

	c.Exclude = out3.Exclude
	c.Include = append(
		out3.Include,
		out2...,
	)

	return nil
}

// UnmarshalYAML unmarshals the condition map.
func (c *ConditionMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	out1 := struct {
		Include map[string]string
		Exclude map[string]string
	}{
		Include: map[string]string{},
		Exclude: map[string]string{},
	}

	out2 := map[string]string{}

	unmarshal(&out1)
	unmarshal(&out2)

	c.Include = out1.Include
	c.Exclude = out1.Exclude
	for k, v := range out2 {
		c.Include[k] = v
	}
	return nil
}

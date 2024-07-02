package parser

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// SigmaRule represents the full structure of a Sigma rule.
type SigmaRule struct {
	Title          string    `yaml:"title"`
	ID             string    `yaml:"id,omitempty"`
	Related        []Related `yaml:"related,omitempty"`
	Status         string    `yaml:"status,omitempty"`
	Description    string    `yaml:"description,omitempty"`
	License        string    `yaml:"license,omitempty"`
	Author         string    `yaml:"author,omitempty"`
	References     []string  `yaml:"references,omitempty"`
	Date           string    `yaml:"date,omitempty"`
	Modified       string    `yaml:"modified,omitempty"`
	LogSource      LogSource `yaml:"logsource"`
	Detection      Detection `yaml:"detection"`
	Fields         []string  `yaml:"fields,omitempty"`
	FalsePositives []string  `yaml:"falsepositives,omitempty"`
	Level          string    `yaml:"level,omitempty"`
	Tags           []string  `yaml:"tags,omitempty"`
}

// Related represents related Sigma rules.
type Related struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
}

// LogSource represents the log source details of the Sigma rule.
type LogSource struct {
	Category string `yaml:"category,omitempty"`
	Product  string `yaml:"product,omitempty"`
	Service  string `yaml:"service,omitempty"`
}

// Detection represents the detection logic of the Sigma rule.
type Detection struct {
	Condition string                 `yaml:"condition"`
	Searches  map[string]interface{} `yaml:",inline"`
}

// ParseSigmaRule reads and parses a Sigma rule from a YAML file.
func ParseSigmaRule(filename string) (SigmaRule, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return SigmaRule{}, fmt.Errorf("failed to read file: %w", err)
	}

	var rule SigmaRule
	err = yaml.Unmarshal(data, &rule)
	if err != nil {
		return SigmaRule{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	log.Printf("Successfully parsed Sigma rule: %s", rule.Title)
	return rule, nil
}

package parser

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseSigmaRule(t *testing.T) {
	// Create a temporary Sigma rule file for testing
	tempFile, err := ioutil.TempFile("", "sigma_rule_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	content := `
title: Test Rule
logsource:
  category: test_category
detection:
  condition: selection
  selection:
    EventID: 1234
`
	if _, err := tempFile.Write([]byte(content)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tempFile.Close()

	// Test ParseSigmaRule function
	rule, err := ParseSigmaRule(tempFile.Name())
	if err != nil {
		t.Fatalf("ParseSigmaRule failed: %v", err)
	}

	if rule.Title != "Test Rule" {
		t.Errorf("Expected title 'Test Rule', got '%s'", rule.Title)
	}
	if rule.LogSource.Category != "test_category" {
		t.Errorf("Expected logsource category 'test_category', got '%s'", rule.LogSource.Category)
	}
	if rule.Detection.Condition != "selection" {
		t.Errorf("Expected detection condition 'selection', got '%s'", rule.Detection.Condition)
	}
}

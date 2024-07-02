package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// SigmaRule represents the full structure of a Sigma rule
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

// Related represents related Sigma rules
type Related struct {
	ID   string `yaml:"id"`
	Type string `yaml:"type"`
}

// LogSource represents the log source details of the Sigma rule
type LogSource struct {
	Category string `yaml:"category,omitempty"`
	Product  string `yaml:"product,omitempty"`
	Service  string `yaml:"service,omitempty"`
}

// Detection represents the detection logic of the Sigma rule
type Detection struct {
	Condition string                 `yaml:"condition"`
	Searches  map[string]interface{} `yaml:",inline"`
}

func main() {
	sigmaRulePath := "path_to_sigma_rule.yaml"
	rule, err := parseSigmaRule(sigmaRulePath)
	if err != nil {
		log.Fatalf("Error parsing Sigma rule: %v", err)
	}

	sqlQuery, err := generateSQL(rule)
	if err != nil {
		log.Fatalf("Error generating SQL query: %v", err)
	}

	fmt.Println(sqlQuery)
}

// parseSigmaRule reads and parses a Sigma rule from a YAML file
func parseSigmaRule(filename string) (SigmaRule, error) {
	// Read the file content
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return SigmaRule{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Unmarshal the YAML content into a SigmaRule struct
	var rule SigmaRule
	err = yaml.Unmarshal(data, &rule)
	if err != nil {
		return SigmaRule{}, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	log.Printf("Successfully parsed Sigma rule: %s", rule.Title)
	return rule, nil
}

// generateSQL generates an SQL query from a parsed Sigma rule
func generateSQL(rule SigmaRule) (string, error) {
	detection := rule.Detection

	// Build the WHERE clause
	whereClause, err := buildWhereClause(detection.Searches)
	if err != nil {
		return "", err
	}

	// Append the timeframe condition if specified
	if timeframe, ok := detection.Searches["timeframe"]; ok {
		whereClause = appendTimeframe(whereClause, timeframe)
	}

	// Build the HAVING clause based on the detection condition
	havingClause, err := buildHavingClause(detection.Condition)
	if err != nil {
		return "", err
	}

	// Construct the final SQL query
	sqlQuery := fmt.Sprintf(`
    SELECT SourceAddress, COUNT(*) 
    FROM events 
    WHERE %s 
    GROUP BY SourceAddress 
    HAVING COUNT(*) > %s
    `, whereClause, havingClause)

	log.Println("Successfully generated SQL query")
	return sqlQuery, nil
}

// buildWhereClause constructs the WHERE clause from detection searches
func buildWhereClause(searches map[string]interface{}) (string, error) {
	var whereClauses []string
	for key, value := range searches {
		conditions, err := generateConditions(key, value)
		if err != nil {
			return "", fmt.Errorf("failed to generate conditions for key '%s': %w", key, err)
		}
		whereClauses = append(whereClauses, conditions)
	}
	return strings.Join(whereClauses, " AND "), nil
}

// appendTimeframe appends a timeframe condition to the WHERE clause
func appendTimeframe(whereClause string, timeframe interface{}) string {
	return fmt.Sprintf("%s AND timestamp >= now() - interval '%s'", whereClause, timeframe)
}

// buildHavingClause constructs the HAVING clause from the detection condition
func buildHavingClause(condition string) (string, error) {
	if condition == "" {
		return "", errors.New("detection condition is empty")
	}

	if strings.Contains(condition, "count() by ") {
		return strings.Split(condition, "count() by ")[1], nil
	}
	return condition, nil
}

// generateConditions generates SQL conditions for a given key and value
func generateConditions(key string, value interface{}) (string, error) {
	var conditions []string
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			conditions = append(conditions, fmt.Sprintf("%s = '%v'", key, item))
		}
		return fmt.Sprintf("(%s)", strings.Join(conditions, " OR ")), nil
	case map[interface{}]interface{}:
		for subKey, subValue := range v {
			subConditions, err := generateConditions(fmt.Sprintf("%s.%s", key, subKey), subValue)
			if err != nil {
				return "", fmt.Errorf("failed to generate conditions for subkey '%s': %w", subKey, err)
			}
			conditions = append(conditions, subConditions)
		}
		return fmt.Sprintf("(%s)", strings.Join(conditions, " AND ")), nil
	default:
		return fmt.Sprintf("%s = '%v'", key, v), nil
	}
}

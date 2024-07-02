package sqlgenerator

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/seacurity/coriolis/parser"
)

// GenerateSQL generates an SQL query from a parsed Sigma rule.
func GenerateSQL(rule parser.SigmaRule) (string, error) {
	detection := rule.Detection

	whereClause, err := buildWhereClause(detection.Searches)
	if err != nil {
		return "", err
	}

	if timeframe, ok := detection.Searches["timeframe"]; ok {
		whereClause = appendTimeframe(whereClause, timeframe)
	}

	havingClause, err := buildHavingClause(detection.Condition)
	if err != nil {
		return "", err
	}

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

// buildWhereClause constructs the WHERE clause from detection searches.
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

// appendTimeframe appends a timeframe condition to the WHERE clause.
func appendTimeframe(whereClause string, timeframe interface{}) string {
	return fmt.Sprintf("%s AND timestamp >= now() - interval '%s'", whereClause, timeframe)
}

// buildHavingClause constructs the HAVING clause from the detection condition.
func buildHavingClause(condition string) (string, error) {
	if condition == "" {
		return "", errors.New("detection condition is empty")
	}

	if strings.Contains(condition, "count() by ") {
		return strings.Split(condition, "count() by ")[1], nil
	}
	return condition, nil
}

// generateConditions generates SQL conditions for a given key and value.
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

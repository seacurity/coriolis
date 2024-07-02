package sql_generator

import (
	"testing"

	"github.com/seacurity/coriolis/parser"
)

func TestGenerateSQL(t *testing.T) {
	rule := parser.SigmaRule{
		Title: "Test Rule",
		LogSource: parser.LogSource{
			Category: "test_category",
		},
		Detection: parser.Detection{
			Condition: "selection",
			Searches: map[string]interface{}{
				"selection": map[interface{}]interface{}{
					"EventID": 1234,
				},
			},
		},
	}

	sqlQuery, err := GenerateSQL(rule)
	if err != nil {
		t.Fatalf("GenerateSQL failed: %v", err)
	}

	expectedQuery := `
    SELECT SourceAddress, COUNT(*) 
    FROM events 
    WHERE (EventID = '1234') 
    GROUP BY SourceAddress 
    HAVING COUNT(*) > selection
    `
	if sqlQuery != expectedQuery {
		t.Errorf("Expected query '%s', got '%s'", expectedQuery, sqlQuery)
	}
}

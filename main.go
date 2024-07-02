package main

import (
	"fmt"
	"log"

	"github.com/seacurity/coriolis/parser"
	"github.com/seacurity/coriolis/scheduler"
	"github.com/seacurity/coriolis/sqlgenerator"
)

func main() {
	sigmaRulePath := "path_to_sigma_rule.yaml"
	rule, err := parser.ParseSigmaRule(sigmaRulePath)
	if err != nil {
		log.Fatalf("Error parsing Sigma rule: %v", err)
	}

	sqlQuery, err := sqlgenerator.GenerateSQL(rule)
	if err != nil {
		log.Fatalf("Error generating SQL query: %v", err)
	}

	fmt.Println(sqlQuery)

	// Start the scheduler to run the search every 5 minutes
	scheduler.StartScheduler(func() {
		rule, err := parser.ParseSigmaRule(sigmaRulePath)
		if err != nil {
			log.Printf("Error parsing Sigma rule: %v", err)
			return
		}

		sqlQuery, err := sqlgenerator.GenerateSQL(rule)
		if err != nil {
			log.Printf("Error generating SQL query: %v", err)
			return
		}

		fmt.Println(sqlQuery)
		// Add logic to execute the SQL query here
	})
}

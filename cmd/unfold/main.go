package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tanuvnair/unfold/internal/config"
	"github.com/tanuvnair/unfold/internal/matcher"
	"github.com/tanuvnair/unfold/internal/parser"
	_ "github.com/tanuvnair/unfold/internal/parser/kotak"
	"github.com/tanuvnair/unfold/internal/report"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: unfold <config.json> <statement.csv>")
	}
	configPath, statementPath := os.Args[1], os.Args[2]

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	bankParser, err := parser.Get(cfg.BankKey())
	if err != nil {
		log.Fatalf("select parser: %v", err)
	}

	file, err := os.Open(statementPath)
	if err != nil {
		log.Fatalf("open statement: %v", err)
	}
	defer file.Close()

	fmt.Printf("Unfolding %s account...\n", cfg.BankName)

	transactions, err := bankParser.Parse(file, cfg)
	if err != nil {
		log.Fatalf("parse statement: %v", err)
	}

	matched := matcher.Filter(transactions, cfg.NormalizedKeywords())

	rpt := report.Build(cfg.BankName, matched)
	outputPath := report.PathFor(statementPath)
	if err := report.Write(outputPath, rpt); err != nil {
		log.Fatalf("write output: %v", err)
	}

	fmt.Printf("Audit complete. Found %d autopay transactions.\n", len(matched))
	fmt.Printf("Report written to %s\n", outputPath)
}

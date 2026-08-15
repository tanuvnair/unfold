package main

import (
	"flag"
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
	bankFlag := flag.String("bank", "", "bank profile to use (name or key); optional when config has one profile")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: unfold [--bank name] <config.json> <statement.csv>\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}
	configPath, statementPath := flag.Arg(0), flag.Arg(1)

	file, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	cfg, err := file.Select(*bankFlag)
	if err != nil {
		log.Fatalf("select profile: %v", err)
	}

	bankParser, err := parser.Get(cfg.BankKey())
	if err != nil {
		log.Fatalf("select parser: %v", err)
	}

	f, err := os.Open(statementPath)
	if err != nil {
		log.Fatalf("open statement: %v", err)
	}
	defer f.Close()

	fmt.Printf("Unfolding %s account...\n", cfg.BankName)

	transactions, err := bankParser.Parse(f, cfg)
	if err != nil {
		log.Fatalf("parse statement: %v", err)
	}

	matched := matcher.Filter(transactions, cfg.NormalizedKeywords(), cfg.NormalizedExcludeKeywords())

	rpt := report.Build(cfg.BankName, matched)
	outputPath := report.PathFor(statementPath)
	if err := report.Write(outputPath, rpt); err != nil {
		log.Fatalf("write output: %v", err)
	}

	fmt.Printf("Audit complete. Found %d autopay transactions.\n", len(matched))
	fmt.Printf("Report written to %s\n", outputPath)
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/tanuvnair/unfold/internal/analyze"
	"github.com/tanuvnair/unfold/internal/config"
	_ "github.com/tanuvnair/unfold/internal/parser/kotak"
	"github.com/tanuvnair/unfold/internal/report"
)

func main() {
	bankFlag := flag.String("bank", "", "bank profile to use (name or key); optional when config has one profile")
	verboseFlag := flag.Bool("verbose", false, "print each matched transaction description")
	flag.BoolVar(verboseFlag, "v", false, "print each matched transaction description")
	diffFlag := flag.Bool("diff", false, "compare against the previous report and print only what changed")
	dryRunFlag := flag.Bool("dry-run", false, "compute the report without writing autopay_report.json")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: unfold [flags] <config.json> <statement.csv>\n")
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

	f, err := os.Open(statementPath)
	if err != nil {
		log.Fatalf("open statement: %v", err)
	}
	defer f.Close()

	fmt.Printf("Unfolding %s account...\n", cfg.BankName)

	rpt, err := analyze.Run(cfg, f)
	if err != nil {
		log.Fatalf("analyze: %v", err)
	}

	outputPath := report.PathFor(statementPath)

	if *diffFlag {
		printDiff(outputPath, rpt)
	}

	if *verboseFlag {
		for _, row := range rpt.Transactions {
			fmt.Println(report.DescriptionOf(row))
		}
	}

	if *dryRunFlag {
		fmt.Printf("Audit complete. Found %d autopay transactions.\n", rpt.TransactionCount)
		fmt.Printf("Dry run: skipped writing %s\n", outputPath)
		return
	}

	if err := report.Write(outputPath, rpt); err != nil {
		log.Fatalf("write output: %v", err)
	}

	fmt.Printf("Audit complete. Found %d autopay transactions.\n", rpt.TransactionCount)
	fmt.Printf("Report written to %s\n", outputPath)
}

func printDiff(outputPath string, next report.Report) {
	prev, err := report.Read(outputPath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("Diff: no previous report at %s; treating all %d as new\n", outputPath, next.TransactionCount)
			for _, row := range next.Transactions {
				fmt.Printf("+ %s\n", report.DescriptionOf(row))
			}
			return
		}
		log.Fatalf("read previous report: %v", err)
	}

	d := report.Diff(prev, next)
	fmt.Printf("Diff: +%d new, -%d removed, =%d unchanged\n", len(d.Added), len(d.Removed), d.Unchanged)
	for _, row := range d.Added {
		fmt.Printf("+ %s\n", report.DescriptionOf(row))
	}
	for _, row := range d.Removed {
		fmt.Printf("- %s\n", report.DescriptionOf(row))
	}
}

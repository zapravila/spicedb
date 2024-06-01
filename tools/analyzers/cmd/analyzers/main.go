package main

import (
	"github.com/zapravila/spicedb/tools/analyzers/closeafterusagecheck"
	"github.com/zapravila/spicedb/tools/analyzers/exprstatementcheck"
	"github.com/zapravila/spicedb/tools/analyzers/nilvaluecheck"
	"github.com/zapravila/spicedb/tools/analyzers/paniccheck"
	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(
		nilvaluecheck.Analyzer(),
		exprstatementcheck.Analyzer(),
		closeafterusagecheck.Analyzer(),
		paniccheck.Analyzer(),
		lendowncastcheck.Analyzer(),
	)
}

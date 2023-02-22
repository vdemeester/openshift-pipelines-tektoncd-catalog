package main

import (
	"flag"
	"os"
	"testing"
	"path/filepath"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress", // can define default values
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts) // godog v0.10.0 and earlier
	godog.BindCommandLineFlags("godog.", &opts)        // godog v0.11.0 and later
}

func InitializeTestSuite(s *godog.TestSuiteContext) {
		fmt.Printf("TestSuite: %+v\n", s)
}
func anImageShouldBePublishedAtHttpsregistryNocode(arg1 int) error {
        return nil
}

func iHaveALocalRegistryRunning() error {
        return nil
}

func iHaveCheckedOutHttpsgithubcomkelseyhightowernocodeInAWorkspace() error {
        return nil
}

func iUseTheBuildTaskWithTheFollowingParameters(ctx *godog.ScenarioContext) func(parameters *godog.Table) error {
		return func(parameters *godog.Table) error {
		params := map[string]string{}
		for _, row := range parameters.Rows[1:] {	
				params[row.Cells[0].Value] = row.Cells[1].Value
		}
        return fmt.Errorf("fooo is bar: %+v\n\n%+v", params, ctx)
		}
}

func theTaskRunShouldSuceed() error {
        return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
        ctx.Step(`^an image should be published at https:\/\/registry:(\d+)\/nocode$`, anImageShouldBePublishedAtHttpsregistryNocode)
        ctx.Step(`^I have a local registry running$`, iHaveALocalRegistryRunning)
        ctx.Step(`^I have checked out https:\/\/github\.com\/kelseyhightower\/nocode in a workspace$`, iHaveCheckedOutHttpsgithubcomkelseyhightowernocodeInAWorkspace)
        ctx.Step(`^I use the build Task with the following parameters:$`, iUseTheBuildTaskWithTheFollowingParameters(ctx))
        ctx.Step(`^the TaskRun should suceed$`, theTaskRunShouldSuceed)
}

func TestMain(m *testing.M) {
		flag.Parse()
		cwd, err := os.Getwd()
		if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(1)
		}
		paths, err := findFeatures(cwd)
		if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err)
				os.Exit(1)
		}
		fmt.Println(paths)
		opts.Paths = paths

	status := godog.TestSuite{
		Name:                 "godogs",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}.Run()

	// Optional: Run `testing` package's logic besides godog.
	if st := m.Run(); st > status {
		status = st
	}

	os.Exit(status)
}

func findFeatures(root string) ([]string, error) {
	var files []string
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && strings.HasSuffix(path, "features") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

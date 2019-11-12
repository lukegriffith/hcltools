package main

import (
	"flag"
	"fmt"
	"os"

	"bufio"
)

func main() {

	var repoFile string
	var checkoutDir string
	var repos []string

	moduleList := &ModuleList{}

	args(&repoFile, &checkoutDir)

	file, err := os.Open(repoFile)

	if err != nil {
		fmt.Println("Unable to open repoFile.")
		panic(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		repos = append(repos, scanner.Text())
	}

	for _, r := range repos {
		parseRepo(r, checkoutDir, moduleList)
	}

	moduleList.PrintModules()

}

func args(repoFile *string, checkoutDir *string) {

	flag.StringVar(repoFile, "repoFile", "", "path to a text file with list of git pull urls.")

	flag.StringVar(checkoutDir, "checkoutDir", "/tmp", "tmp path to check repos out to.")

	flag.Parse()

	if *repoFile == "" || *checkoutDir == "" {
		flag.Usage()
		os.Exit(0)
	}
}


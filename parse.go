package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	sep = "/"
)

func parseRepo(repo string, checkoutDir string, moduleList *ModuleList) error {

	log.Println("Testing repo", repo)

	directory := checkoutDir + "/tfexplicitver"

	repository, err := git.PlainClone(directory, false, &git.CloneOptions{
		URL: repo,
	})

	if err != nil {
		return err
	}

	defer os.RemoveAll(directory)
	defer os.Remove(directory)

	branches := []string{"master", "devel", "staging", "production", "ccc.local", "c3.zone"}

	for _, b := range branches {

		err = checkoutBranch(repository, b)

		if err != nil {

			log.Println(err, b)
			continue
		}

		err = enumerateDirectory(directory, repo, b, moduleList)

		if err != nil {
			return err
		}

	}

	return nil

}

func checkoutBranch(r *git.Repository, branch string) error {

	w, err := r.Worktree()

	err = r.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*", "HEAD:refs/heads/HEAD"},
	})
	if err != nil {
		fmt.Println(err)
	}

	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branch),
		Force:  true,
	})
	if err != nil {
		return err
	}

	return nil

}

func enumerateDirectory(directory string, repo string, branch string, modList *ModuleList) error {

	log.Println("enumerating ", directory)

	files, err := ioutil.ReadDir(directory)

	if err != nil {
		return err
	}

	for _, f := range files {

		//log.Println(f.Name())

		fileName := directory + sep + f.Name()

		if f.IsDir() {
			err := enumerateDirectory(fileName, repo, branch, modList)

			if err != nil {
				log.Println(err)
			}

		} else {
			err, badModules := parseForHCL(fileName)

			if err != nil {

				log.Println(directory)
				log.Println(err)

				for _, v := range badModules {
					modList.AddModule(repo, v, fileName, branch)
				}

			}
		}

	}

	return nil

}

func parseForHCL(file string) (error, []string) {

	isTf := strings.HasSuffix(file, ".tf")

	if !isTf {
		return nil, nil
	}

	log.Println(file)

	b, err := ioutil.ReadFile(file)

	if err != nil {
		return err, nil
	}

	var parsedAst *ast.File

	parsedAst, err = hcl.ParseBytes(b)

	if err != nil {
		return err, nil
	}

	var foo map[string]interface{}

	err = hcl.DecodeObject(&foo, parsedAst)

	if err != nil {
		log.Println(err)
	}

	badModules := moduleParse(foo)

	if len(badModules) > 0 {
		return errors.New(fmt.Sprintf("bad modules, file: %s modules: \n%s", file, strings.Join(badModules, "\n"))), badModules
	}

	return nil, badModules

}

type badModules struct {
	Modules []string
}

func (b *badModules) appendModule(m string) {

	b.Modules = append(b.Modules, m)

}

func moduleParse(astMap map[string]interface{}) []string {

	// Convert modules key to array of map[string]interface{}

	modKey, ok := astMap["module"]

	if !ok {
		return nil
	}

	modules := modKey.([]map[string]interface{})

	bm := &badModules{}

	for _, v := range modules {
		for moduleName, v := range v {

			innermod := v.([]map[string]interface{})

			for _, v := range innermod {

				sourceVal, ok := v["source"]

				if !ok {
					// No source defined, this will break terraform anyway...

					log.Println("No source defined...")
					continue
				}

				source := sourceVal.(string)

				if strings.HasPrefix(source, ".") {
					// using relative path, no need to explicitly version
					log.Println("Relative path used, no explicit version needed.")
					continue
				}

				_, ok = v["version"]

				if !ok {

					bm.appendModule(moduleName)
				}

			}
		}
	}

	return bm.Modules

}


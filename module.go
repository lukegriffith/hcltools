package main

import "fmt"

type Module struct {
	Repo   string
	Name   string
	File   string
	Branch string
}

type ModuleList struct {
	Modules []Module
}

func (m *ModuleList) AddModule(r string, n string, f string, b string) {

	m.Modules = append(m.Modules, Module{r, n, f, b})
}

func (m *ModuleList) PrintModules() {

	fmt.Println()
	fmt.Println()

	fmt.Println("Repo, Name, File, Branch")

	for _, v := range m.Modules {
		fmt.Println(v.Repo, ",", v.Name, ",", v.File, ",", v.Branch)
	}

}


package main

type hclStrings struct {
	Strings []string
}

func (m *hclStrings) AddString(r string) {

	m.Strings = append(m.Strings, r)
}


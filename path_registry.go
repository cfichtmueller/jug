// Copyright 2023 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package jug

type PathRegistry struct {
	paths map[string]map[string]bool
}

func NewPathRegistry() *PathRegistry {
	return &PathRegistry{
		paths: make(map[string]map[string]bool),
	}
}

func (p *PathRegistry) Add(path string, methods ...string) {
	_, ok := p.paths[path]
	if !ok {
		p.paths[path] = make(map[string]bool)
	}
	e := p.paths[path]
	for _, m := range methods {
		e[m] = true
	}
}

func (p *PathRegistry) Get(path string, method string) bool {
	e, ok := p.paths[path]
	if !ok {
		return false
	}
	_, ok = e[method]
	return ok
}

func (p *PathRegistry) Paths() []string {
	paths := make([]string, 0, len(p.paths))
	for k, _ := range p.paths {
		paths = append(paths, k)
	}
	return paths
}

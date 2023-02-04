// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"github.com/dvcrn/go-1password-cli/op"
	"log"
	"sigs.k8s.io/kustomize/api/kv"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
	"strings"
)

type OpValue struct {
	Key    string `json:"key,omitempty" yaml:"key,omitempty"`
	OpPath string `json:"opPath,omitempty" yaml:"opPath,omitempty"`
}

// A secret generator example that gets data
// from a database (simulated by a hardcoded map).
type plugin struct {
	h                       *resmap.PluginHelpers
	types.ObjectMeta        `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	*types.GeneratorOptions `json:"options,omitempty" yaml:"options,omitempty" protobuf:"bytes,2,opt,name=options"`

	Type string `json:"type,omitempty" yaml:"type,omitempty"`

	// List of keys to use in database lookups
	Values []OpValue  `json:"values,omitempty" yaml:"values,omitempty"`
	opCli  *op.Client `json:"-" yaml:"-"`
}

var KustomizePlugin plugin //nolint:gochecknoglobals

func (p *plugin) Config(h *resmap.PluginHelpers, c []byte) error {
	p.h = h
	p.opCli = op.NewOpClient()
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace
	args.Type = "Opaque"

	// take over type if set
	if p.Type != "" {
		args.Type = p.Type
	}

	if p.GeneratorOptions != nil {
		args.Options = p.GeneratorOptions
	}

	for _, v := range p.Values {
		opPathFull := v.OpPath

		// trim first / if it exists
		if strings.HasPrefix(opPathFull, "/") {
			opPathFull = strings.TrimPrefix(opPathFull, "/")
		}

		if !strings.HasPrefix(opPathFull, "op://") {
			opPathFull = fmt.Sprintf("op://%s", opPathFull)
		}

		opValue, err := p.opCli.Read(opPathFull)
		if err != nil {
			log.Fatal("err talking to 1password: ", err)
		}

		args.LiteralSources = append(
			args.LiteralSources, v.Key+"="+opValue)
	}

	return p.h.ResmapFactory().FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), args)
}

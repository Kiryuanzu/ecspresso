package ecspresso

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/fujiwara/cfn-lookup/cfn"
	"github.com/fujiwara/tfstate-lookup/tfstate"

	"github.com/kayac/ecspresso/v2/ssm"
)

type ConfigPlugin struct {
	Name       string                 `yaml:"name" json:"name"`
	Config     map[string]interface{} `yaml:"config" json:"config"`
	FuncPrefix string                 `yaml:"func_prefix,omitempty" json:"func_prefix,omitempty"`
}

func (p ConfigPlugin) Setup(ctx context.Context, c *Config) error {
	switch strings.ToLower(p.Name) {
	case "tfstate":
		return setupPluginTFState(ctx, p, c)
	case "cloudformation":
		return setupPluginCFn(ctx, p, c)
	case "ssm":
		return setupPluginSSM(ctx, p, c)
	default:
		return fmt.Errorf("plugin %s is not available", p.Name)
	}
}

func (p ConfigPlugin) AppendFuncMap(c *Config, funcMap template.FuncMap) error {
	modified := make(template.FuncMap, len(funcMap))
	for funcName, f := range funcMap {
		name := p.FuncPrefix + funcName
		for _, appendedFuncs := range c.templateFuncs {
			if _, exists := appendedFuncs[name]; exists {
				return fmt.Errorf("template function %s already exists. set func_prefix to %s plugin", name, p.Name)
			}
		}
		modified[name] = f
	}
	c.templateFuncs = append(c.templateFuncs, modified)
	return nil
}

func setupPluginTFState(ctx context.Context, p ConfigPlugin, c *Config) error {
	var loc string
	if p.Config["path"] != nil {
		path, ok := p.Config["path"].(string)
		if !ok {
			return errors.New("tfstate plugin requires path for tfstate file as a string")
		}
		if !filepath.IsAbs(path) {
			path = filepath.Join(c.dir, path)
		}
		loc = path
	} else if p.Config["url"] != nil {
		u, ok := p.Config["url"].(string)
		if !ok {
			return errors.New("tfstate plugin requires url for tfstate URL as a string")
		}
		loc = u
	} else {
		return errors.New("tfstate plugin requires path or url for tfstate location")
	}
	funcs, err := tfstate.FuncMap(ctx, loc)
	if err != nil {
		return err
	}
	return p.AppendFuncMap(c, funcs)
}

func setupPluginCFn(ctx context.Context, p ConfigPlugin, c *Config) error {
	funcs, err := cfn.FuncMap(ctx, c.awsv2Config)
	if err != nil {
		return err
	}
	return p.AppendFuncMap(c, funcs)
}

func setupPluginSSM(ctx context.Context, p ConfigPlugin, c *Config) error {
	funcs, err := ssm.FuncMap(ctx, c.awsv2Config)
	if err != nil {
		return err
	}
	return p.AppendFuncMap(c, funcs)
}

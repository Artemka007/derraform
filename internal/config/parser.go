package config

import (
    "github.com/hashicorp/hcl/v2"
    "github.com/hashicorp/hcl/v2/hclsyntax"
    "github.com/zclconf/go-cty/cty"
)

type Config struct {
    Resources []Resource
}

type Resource struct {
    Type      string
    Name      string
    Attributes map[string]cty.Value
}

func ParseFile(filename string) (*Config, error) {
    src, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }

    file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
    if diags.HasErrors() {
        return nil, diags
    }

    return parseConfig(file)
}
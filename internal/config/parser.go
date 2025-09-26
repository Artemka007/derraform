// internal/config/parser.go
package config

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"
)

type Config struct {
	Resources []Resource
}

type Resource struct {
	Type       string
	Name       string
	Attributes map[string]cty.Value
}

func ParseFile(filename string) (*Config, error) {
	src, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	file, diags := hclsyntax.ParseConfig(src, filename, hcl.Pos{Line: 1, Column: 1})
	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse HCL: %s", diags.Error())
	}

	return parseConfig(file)
}

// parseConfig преобразует HCL AST в нашу конфигурацию
func parseConfig(file *hcl.File) (*Config, error) {
	config := &Config{
		Resources: []Resource{},
	}

	// Получаем корневой body
	content, diags := file.Body.Content(&hcl.BodySchema{
		Blocks: []hcl.BlockHeaderSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
			},
		},
	})

	if diags.HasErrors() {
		return nil, fmt.Errorf("failed to parse body: %s", diags.Error())
	}

	// Обрабатываем все resource блоки
	for _, block := range content.Blocks {
		if block.Type == "resource" {
			resource, err := parseResourceBlock(block)
			if err != nil {
				return nil, fmt.Errorf("failed to parse resource block: %w", err)
			}
			config.Resources = append(config.Resources, resource)
		}
	}

	return config, nil
}

// parseResourceBlock парсит отдельный resource блок
func parseResourceBlock(block *hcl.Block) (Resource, error) {
	resource := Resource{
		Type:       block.Labels[0],
		Name:       block.Labels[1],
		Attributes: make(map[string]cty.Value),
	}

	// Парсим атрибуты внутри resource блока
	attrs, diags := block.Body.JustAttributes()
	if diags.HasErrors() {
		// Игнорируем ошибки атрибутов, возможно есть nested blocks
		diags = nil
	}

	for name, attr := range attrs {
		value, diags := attr.Expr.Value(nil) // nil EvalContext для простоты
		if diags.HasErrors() {
			return resource, fmt.Errorf("failed to evaluate attribute %s: %s", name, diags.Error())
		}
		resource.Attributes[name] = value
	}

	return resource, nil
}

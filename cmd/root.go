package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/aiven/go-api-schemas/internal/diff"
	"github.com/aiven/go-api-schemas/internal/gen"
	"github.com/aiven/go-api-schemas/internal/types"
)

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use: "go-api-schemas foo.json bar.json baz.json",
		Short: "go-api-schemas is a tool for generating and persisting user configuration option schemas from " +
			"Aiven APIs.",
		RunE: run,
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringP("output-dir", "o", "pkg/dist", "the output directory for the generated files")
	cmd.Flags().BoolP(
		"regenerate", "r", false,
		"regenerate files without comparing against existing files (useful for removing deprecations)",
	)
	return cmd
}

func run(cmd *cobra.Command, fileNames []string) error {
	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return fmt.Errorf("error getting output directory: %w", err)
	}

	regenerate, err := cmd.Flags().GetBool("regenerate")
	if err != nil {
		return fmt.Errorf("error getting regeneration flag: %w", err)
	}

	generationResult, err := gen.Run(fileNames...)
	if err != nil {
		return fmt.Errorf("error generating: %w", err)
	}

	readResult := make(types.ReadResult)
	if !regenerate {
		readResult, err = read(outputDir)
		if err != nil {
			return fmt.Errorf("error reading files: %w", err)
		}
	}

	diffResult, err := diff.Diff(readResult, generationResult)
	if err != nil {
		return fmt.Errorf("error diffing schemas: %w", err)
	}

	err = write(outputDir, diffResult)
	if err != nil {
		return fmt.Errorf("error writing files: %w", err)
	}
	return nil
}

func read(outputDir string) (types.ReadResult, error) {
	result := make(types.ReadResult)
	for k, v := range getSchemaFilenames() {
		result[k] = make(map[string]*types.UserConfigSchema)
		filePath := filepath.Join(outputDir, v)
		err := readFile(filePath, result[k])
		if err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("%q: %w", filePath, err)
		}
	}

	return result, nil
}

func readFile(filePath string, schema map[string]*types.UserConfigSchema) error {
	f, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer f.Close()
	return yaml.NewDecoder(f).Decode(schema)
}

func write(outputDir string, result types.DiffResult) error {
	for k, v := range getSchemaFilenames() {
		p := filepath.Join(outputDir, v)
		err := writeFile(p, result[k])
		if err != nil {
			return fmt.Errorf("%q: %w", p, err)
		}
	}

	return nil
}

func writeFile(filePath string, schema map[string]*types.UserConfigSchema) error {
	f, err := os.Create(filepath.Clean(filePath))
	if err != nil {
		return err
	}

	defer f.Close()

	e := yaml.NewEncoder(f)
	defer e.Close()

	const indentSpaces = 2
	e.SetIndent(indentSpaces)
	return e.Encode(schema)
}

const (
	serviceSchemaFilename             = "service_types.yml"
	integrationSchemaFilename         = "integration_types.yml"
	integrationEndpointSchemaFilename = "integration_endpoint_types.yml"
)

func getSchemaFilenames() map[types.SchemaType]string {
	return map[types.SchemaType]string{
		types.ServiceSchemaType:             serviceSchemaFilename,
		types.IntegrationSchemaType:         integrationSchemaFilename,
		types.IntegrationEndpointSchemaType: integrationEndpointSchemaFilename,
	}
}

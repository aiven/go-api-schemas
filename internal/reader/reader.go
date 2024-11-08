// Package reader is the package that contains the reader functionality.
package reader

import (
	"os"
	"strings"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

// logger is a pointer to the logger.
var logger *util.Logger

// flags is a pointer to the flags.
var flags *pflag.FlagSet

// result is the result of the read process.
var result types.ReadResult

// read is a function that reads a file and returns the contents as a map[string]types.UserConfigSchema.
func read(filename string, schema map[string]types.UserConfigSchema) error {
	logger.Info.Printf("reading %s", filename)

	outputDir, err := flags.GetString("output-dir")
	if err != nil {
		return err
	}

	f, err := os.Open(strings.Join([]string{outputDir, filename}, string(os.PathSeparator)))
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err = f.Close()
	}(f)

	d := yaml.NewDecoder(f)

	if err = d.Decode(schema); err != nil {
		return err
	}

	return err
}

// readServiceTypes reads the service types from a file.
func readServiceTypes() error {
	return read(util.ServiceTypesFilename, result[types.KeyServiceTypes])
}

// readIntegrationTypes reads the integration types from a file.
func readIntegrationTypes() error {
	return read(util.IntegrationTypesFilename, result[types.KeyIntegrationTypes])
}

// readIntegrationEndpointTypes reads the integration endpoint types from a file.
func readIntegrationEndpointTypes() error {
	return read(util.IntegrationEndpointTypesFilename, result[types.KeyIntegrationEndpointTypes])
}

// setup sets up the reader.
func setup(l *util.Logger, f *pflag.FlagSet) {
	logger = l
	flags = f

	result = types.ReadResult{
		types.KeyServiceTypes:             make(map[string]types.UserConfigSchema),
		types.KeyIntegrationTypes:         make(map[string]types.UserConfigSchema),
		types.KeyIntegrationEndpointTypes: make(map[string]types.UserConfigSchema),
	}
}

// Run runs the reader.
func Run(logger *util.Logger, flags *pflag.FlagSet) (types.ReadResult, error) {
	setup(logger, flags)

	if err := readServiceTypes(); err != nil {
		return nil, err
	}

	if err := readIntegrationTypes(); err != nil {
		return nil, err
	}

	if err := readIntegrationEndpointTypes(); err != nil {
		return nil, err
	}

	return result, nil
}

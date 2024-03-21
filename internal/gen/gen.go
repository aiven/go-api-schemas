// Package gen is the package that contains the generation logic.
package gen

import (
	"github.com/aiven/aiven-go-client"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"

	"github.com/aiven/go-api-schemas/internal/convert"
	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
)

const (
	// generating is a part of the message that is printed when the generation process starts.
	generating = "generating %s"
)

// logger is a pointer to the logger.
var logger *util.Logger

// env is a map of environment variables.
var env util.EnvMap

// client is a pointer to the Aiven client.
var client *aiven.Client

// result is the result of the generation process.
var result types.GenerationResult

// serviceTypes generates the service types.
func serviceTypes() error {
	defer util.MeasureExecutionTime(logger)()

	logger.Info.Printf(generating, "service types")

	r, err := client.Projects.ServiceTypes(env[util.EnvAivenProjectName])
	if err != nil {
		return err
	}

	out := make(map[string]types.UserConfigSchema, len(r))

	for k, v := range r {
		cv, err := convert.UserConfigSchema(v.UserConfigSchema)
		if err != nil {
			return err
		}

		out[k] = *cv
	}

	result[types.KeyServiceTypes] = out

	return nil
}

// integrationTypes generates the integration types.
func integrationTypes() error {
	defer util.MeasureExecutionTime(logger)()

	logger.Info.Printf(generating, "integration types")

	r, err := client.Projects.IntegrationTypes(env[util.EnvAivenProjectName])
	if err != nil {
		return err
	}

	out := make(map[string]types.UserConfigSchema, len(r))

	for _, v := range r {
		cv, err := convert.UserConfigSchema(v.UserConfigSchema)
		if err != nil {
			return err
		}

		out[v.IntegrationType] = *cv
	}

	result[types.KeyIntegrationTypes] = out

	return nil
}

// integrationEndpointTypes generates the integration endpoint types.
func integrationEndpointTypes() error {
	defer util.MeasureExecutionTime(logger)()

	logger.Info.Printf(generating, "integration endpoint types")

	r, err := client.Projects.IntegrationEndpointTypes(env[util.EnvAivenProjectName])
	if err != nil {
		return err
	}

	out := make(map[string]types.UserConfigSchema, len(r))

	for _, v := range r {
		cv, err := convert.UserConfigSchema(v.UserConfigSchema)
		if err != nil {
			return err
		}

		out[v.EndpointType] = *cv
	}

	result[types.KeyIntegrationEndpointTypes] = out

	return nil
}

// setup sets up the generation process.
func setup(l *util.Logger, e util.EnvMap, c *aiven.Client) {
	logger = l
	env = e
	client = c

	result = types.GenerationResult{}
}

// Run executes the generation process.
func Run(
	ctx context.Context,
	logger *util.Logger,
	env util.EnvMap,
	client *aiven.Client,
) (types.GenerationResult, error) {
	setup(logger, env, client)

	errs, _ := errgroup.WithContext(ctx)

	errs.Go(serviceTypes)
	errs.Go(integrationTypes)
	errs.Go(integrationEndpointTypes)

	return result, errs.Wait()
}

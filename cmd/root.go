// Package cmd is the package that contains the commands of the application.
package cmd

import (
	"errors"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/aiven/go-api-schemas/internal/diff"
	"github.com/aiven/go-api-schemas/internal/gen"
	"github.com/aiven/go-api-schemas/internal/pkg/types"
	"github.com/aiven/go-api-schemas/internal/pkg/util"
	"github.com/aiven/go-api-schemas/internal/reader"
	"github.com/aiven/go-api-schemas/internal/writer"
)

// logger is the logger of the application.
var logger = &util.Logger{}

// NewCmdRoot returns a pointer to the root command.
func NewCmdRoot(l *util.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use: "go-api-schemas foo.json bar.json baz.json",
		Short: "go-api-schemas is a tool for generating and persisting user configuration option schemas from " +
			"Aiven APIs.",
		Run:  run,
		Args: cobra.MinimumNArgs(1),
	}

	cmd.Flags().StringP("output-dir", "o", "", "the output directory for the generated files")

	cmd.Flags().BoolP("regenerate", "r", false, "regenerates the files in the output directory")

	logger = l

	return cmd
}

// setupOutputDir sets up the output directory.
func setupOutputDir(flags *pflag.FlagSet) error {
	outputDir, err := flags.GetString("output-dir")
	if err != nil {
		return err
	}

	if outputDir == "" {
		var wd string

		wd, err = os.Getwd()
		if err != nil {
			return err
		}

		outputDir = strings.Join([]string{wd, "pkg/dist"}, string(os.PathSeparator))

		err = flags.Set("output-dir", outputDir)
		if err != nil {
			return err
		}
	}

	fi, err := os.Stat(outputDir)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return errors.New("output directory is not a directory")
	}

	return nil
}

// setup sets up the application.
func setup(flags *pflag.FlagSet) {
	logger.Info.Println("go-api-schemas tool started")

	logger.Info.Println("setting up output directory")

	if err := setupOutputDir(flags); err != nil {
		logger.Error.Fatalf("error setting up output directory: %s", err)
	}

	logger.Info.Println("setting up environment variables")
}

// run is the function that is called when the rootCmd is executed.
func run(cmd *cobra.Command, args []string) {
	flags := cmd.Flags()

	setup(flags)

	shouldRegenerate, err := flags.GetBool("regenerate")
	if err != nil {
		logger.Error.Fatalf("error getting regeneration flag: %s", err)
	}

	logger.Info.Println("generating")

	gr, err := gen.Run(args...)
	if err != nil {
		logger.Error.Fatalf("error generating: %s", err)
	}

	rr := make(types.ReadResult)

	if !shouldRegenerate {
		logger.Info.Println("reading files")

		rr, err = reader.Run(logger, flags)
		if err != nil && !os.IsNotExist(err) {
			logger.Error.Fatalf("error reading files: %s", err)
		}
	}

	logger.Info.Println("diffing")
	dr, err := diff.Run(rr, gr)
	if err != nil {
		logger.Error.Fatalf("error diffing: %s", err)
	}

	logger.Info.Println("writing files")
	if err = writer.Run(logger, flags, dr); err != nil {
		logger.Error.Fatalf("error writing files: %s", err)
	}

	logger.Info.Println("done")
}

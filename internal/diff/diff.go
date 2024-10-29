// Package diff is the package that contains the diffMaps functionality.
package diff

import (
	"fmt"

	"golang.org/x/exp/slices"

	"github.com/aiven/go-api-schemas/internal/pkg/types"
)

// Run runs the diffMaps.
func Run(was types.ReadResult, have types.GenerationResult) (types.DiffResult, error) {
	result := make(types.DiffResult)
	for _, k := range types.GetTypeKeys() {
		result[k] = diffMaps(was[k], have[k])
	}

	return result, nil
}

func diffTwo(was, have *types.UserConfigSchema) *types.UserConfigSchema {
	switch {
	case was == nil:
		return have
	case have == nil:
		was.Deprecate("This property is deprecated.")
		return was
	}

	// Properties
	have.Properties = diffMaps(was.Properties, have.Properties)
	have.Items = diffTwo(was.Items, have.Items)
	have.OneOf = diffArrays(was.OneOf, have.OneOf)
	have.Enum = diffEnums(was.Enum, have.Enum)
	return have
}

func diffEnums(was, have []types.UserConfigSchemaEnumValue) []types.UserConfigSchemaEnumValue {
	r := make(map[string]types.UserConfigSchemaEnumValue)
	for _, v := range have {
		r[stringify(v.Value)] = v
	}

	for _, v := range was {
		k := stringify(v.Value)
		if _, ok := r[k]; !ok {
			v.Deprecate("This value is deprecated.")
			r[k] = v
		}
	}

	return mapValues(r)
}

func diffArrays(was []types.UserConfigSchema, have []types.UserConfigSchema) []types.UserConfigSchema {
	r := make(map[string]types.UserConfigSchema)
	for _, v := range have {
		r[stringify(v.Type)] = v
	}

	for _, w := range was {
		k := stringify(w.Type)
		h, ok := r[k]
		if !ok {
			w.Deprecate("This item is deprecated.")
			r[k] = w
			continue
		}

		r[k] = *diffTwo(&w, &h)
	}

	return mapValues(r)
}

// diffMaps returns the difference between the two maps.
// WARNING: Mutates the input maps.
func diffMaps(was, have map[string]types.UserConfigSchema) map[string]types.UserConfigSchema {
	keys := mergeKeys(was, have)
	if len(keys) == 0 {
		return nil
	}

	r := make(map[string]types.UserConfigSchema)
	for _, k := range keys {
		var w, h *types.UserConfigSchema
		if v, ok := was[k]; ok {
			w = &v
		}

		if v, ok := have[k]; ok {
			h = &v
		}

		r[k] = *diffTwo(w, h)
	}

	return r
}

func stringify(v any) string {
	return fmt.Sprintf("%v", v)
}

// mergeKeys merges the keys of the given maps and returns them sorted.
func mergeKeys[T any](args ...map[string]T) []string {
	if len(args) == 0 {
		return nil
	}

	seen := make(map[string]bool)
	for _, m := range args {
		for k := range m {
			seen[k] = true
		}
	}

	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	return keys
}

// mapValues returns the values of the given map sorted by the keys.
func mapValues[T any](m map[string]T) []T {
	if len(m) == 0 {
		return nil
	}

	list := make([]T, 0, len(m))
	for _, k := range mergeKeys(m) {
		list = append(list, m[k])
	}

	return list
}

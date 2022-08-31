package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPredicate(t *testing.T) {
	for _, predicate := range PredicatePossibleValues {
		t.Run("testing with predicate "+predicate, func(t *testing.T) {
			got := Predicate{
				Name:  predicate,
				Value: "value",
			}

			assert.Contains(t, got.ToPath(), predicate)
		})
	}

	t.Run("testing with some predicate and value with equal sign", func(t *testing.T) {
		got := Predicate{
			Name:  "somepredicate",
			Value: "=value=",
		}

		assert.NotContains(t, got.ToPath(), "=")
	})
}

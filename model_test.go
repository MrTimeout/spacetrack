package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSpaceTrackObjFromArr(t *testing.T) {
	t.Run("space-track tle unit", func(t *testing.T) {
		var input = make([]SpaceTrackTleUnit, 10)

		assert.IsType(t, SpaceTrackTle{}, newSpaceTrackObjFromArr(input))
	})

	t.Run("space-track tle decay", func(t *testing.T) {
		var input = make([]SpaceTrackDecayUnit, 10)

		assert.IsType(t, SpaceTrackDecay{}, newSpaceTrackObjFromArr(input))
	})

	t.Run("space-track tle cdm", func(t *testing.T) {
		var input = make([]SpaceTrackCdmUnit, 10)

		assert.IsType(t, SpaceTrackCdm{}, newSpaceTrackObjFromArr(input))
	})
}

package starr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golift.io/starr"
	"golift.io/starr/prowlarr"
	"golift.io/starr/sonarr"
)

func TestCopy(t *testing.T) {
	t.Parallel()

	src := &prowlarr.IndexerOutput{
		ID:             2,
		Priority:       3,
		Name:           "yes",
		Protocol:       "usenet",
		Implementation: "core",
		ConfigContract: "hancock",
		Tags:           []int{1, 2, 5},
		Fields: []*starr.FieldOutput{
			{Name: "One", Value: "one"},
			{Name: "Two", Value: 2.0},
			{Name: "Three", Value: uint(3)},
			{Name: "Five", Value: 5},
		},
	}
	dst := &sonarr.IndexerInput{}

	require.NoError(t, starr.Copy(src, dst))
	assert.Equal(t, src.Fields[0].Value, dst.Fields[0].Value)
	assert.Equal(t, src.Fields[1].Value, dst.Fields[1].Value)
	assert.EqualValues(t, src.Fields[2].Value, dst.Fields[2].Value)
	assert.EqualValues(t, src.Fields[3].Value, dst.Fields[3].Value)
	assert.Equal(t, src.Fields[0].Name, dst.Fields[0].Name)
	assert.Equal(t, src.Fields[1].Name, dst.Fields[1].Name)
	assert.Equal(t, src.Fields[2].Name, dst.Fields[2].Name)
	assert.Equal(t, src.Fields[3].Name, dst.Fields[3].Name)
	assert.Equal(t, src.ID, dst.ID)
	assert.Equal(t, src.Priority, dst.Priority)
	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Protocol, dst.Protocol)
	assert.Equal(t, src.Implementation, dst.Implementation)
	assert.Equal(t, src.ConfigContract, dst.ConfigContract)
	assert.Equal(t, src.Tags[1], dst.Tags[1])

	broken := struct{}{}
	require.ErrorIs(t, starr.Copy(broken, dst), starr.ErrNotPtr)
	require.ErrorIs(t, starr.Copy(src, broken), starr.ErrNotPtr)
}

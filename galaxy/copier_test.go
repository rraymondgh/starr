package galaxy_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golift.io/starr"
	"golift.io/starr/galaxy"
	"golift.io/starr/prowlarr"
	"golift.io/starr/sonarr"
)

func TestCopyIndexer(t *testing.T) {
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
	// This is a real example of how you'd copy an indexer from Prowlarr to Sonarr.
	dst := &sonarr.IndexerInput{
		// These are not part of the used input, so set them before copying.
		EnableAutomaticSearch:   true,
		EnableInteractiveSearch: true,
		EnableRss:               true,
		DownloadClientID:        15,
	}

	// Verify everything copies over.
	_, err := galaxy.CopyIndexer(src, dst, true)
	require.NoError(t, err)
	assert.Equal(t, src.Fields[0].Value, dst.Fields[0].Value)
	assert.Equal(t, src.Fields[1].Value, dst.Fields[1].Value)
	assert.EqualValues(t, src.Fields[2].Value, dst.Fields[2].Value)
	assert.EqualValues(t, src.Fields[3].Value, dst.Fields[3].Value)
	assert.Equal(t, src.Fields[0].Name, dst.Fields[0].Name)
	assert.Equal(t, src.Fields[1].Name, dst.Fields[1].Name)
	assert.Equal(t, src.Fields[2].Name, dst.Fields[2].Name)
	assert.Equal(t, src.Fields[3].Name, dst.Fields[3].Name)
	assert.Zero(t, dst.ID)
	assert.Equal(t, src.Priority, dst.Priority)
	assert.Equal(t, src.Name, dst.Name)
	assert.Equal(t, src.Protocol, dst.Protocol)
	assert.Equal(t, src.Implementation, dst.Implementation)
	assert.Equal(t, src.ConfigContract, dst.ConfigContract)
	assert.Equal(t, src.Tags[0], dst.Tags[0])
	assert.Equal(t, src.Tags[1], dst.Tags[1])
	assert.Equal(t, src.Tags[2], dst.Tags[2])
	// Check passed in values.
	assert.Equal(t, int64(15), dst.DownloadClientID)
	assert.True(t, dst.EnableAutomaticSearch)
	assert.True(t, dst.EnableInteractiveSearch)
	assert.True(t, dst.EnableRss)
	// Make sure tags get depleted.
	galaxy.Must(galaxy.CopyIndexer(src, dst, false))
	assert.Zero(t, dst.Tags)
}

func TestCopy(t *testing.T) {
	broken := struct{}{}
	good := &prowlarr.IndexerOutput{}

	require.ErrorIs(t, galaxy.Copy(broken, good), galaxy.ErrNotPtr)
	require.ErrorIs(t, galaxy.Copy(good, broken), galaxy.ErrNotPtr)
}

package starr_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golift.io/starr"
)

func TestQueueDeleteOpts_Values(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var opts *starr.QueueDeleteOpts

	params := opts.Values() // it's nil.
	require.Equal(t, "removeFromClient=true", params.Encode(),
		"default queue delete parameters encoded incorrectly")

	opts = &starr.QueueDeleteOpts{
		BlockList:        true,
		RemoveFromClient: starr.False(),
		SkipRedownload:   true,
	}
	params = opts.Values()

	assert.Equal("false", params.Get("removeFromClient"), "delete parameters encoded incorrectly")
	assert.Equal("true", params.Get("blocklist"), "delete parameters encoded incorrectly")
	assert.Equal("true", params.Get("skipRedownload"), "delete parameters encoded incorrectly")
}

func TestNone(t *testing.T) {
	t.Parallel()
	assert.ErrorIs(t, starr.None(starr.ErrNilClient), starr.ErrNilClient)
	assert.ErrorIs(t, starr.None("string", starr.ErrNilClient), starr.ErrNilClient)
	assert.ErrorIs(t, starr.None(uint(1), starr.ErrNilClient), starr.ErrNilClient)
	assert.ErrorIs(t, starr.None("string", uint(1), starr.ErrNilClient), starr.ErrNilClient)
	assert.ErrorIs(t, starr.None(1.0, "string", starr.ErrNilClient), starr.ErrNilClient)
	require.NoError(t, starr.None(1.0, "string"))
	require.NoError(t, starr.None("string"))
	require.NoError(t, starr.None(1.0))
	require.NoError(t, starr.None())
}

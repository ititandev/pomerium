package databroker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	databrokerpb "github.com/pomerium/pomerium/pkg/grpc/databroker"
)

// A wrong type URL in setupRequiredIndex fails open silently: buildRecordTTLs just skips a
// type whose options are missing, so the records would accumulate with no error anywhere.
func TestSetupRequiredIndex_MCPRecordOptions(t *testing.T) {
	t.Parallel()

	srv := newServer(t)

	for _, tc := range []struct {
		typeURL  string
		ttl      time.Duration
		wantCap  bool
		capacity uint64
	}{
		{"type.googleapis.com/oauth21.MCPRefreshToken", 30 * 24 * time.Hour, true, 10000},
		{"type.googleapis.com/oauth21.AuthorizationRequest", time.Hour, false, 0},
	} {
		res, err := srv.GetOptions(t.Context(), &databrokerpb.GetOptionsRequest{Type: tc.typeURL})
		require.NoError(t, err, tc.typeURL)
		require.NotNil(t, res.GetOptions(), tc.typeURL)

		assert.Equal(t, tc.ttl, res.GetOptions().GetTtl().AsDuration(), tc.typeURL)
		if tc.wantCap {
			assert.Equal(t, tc.capacity, res.GetOptions().GetCapacity(), tc.typeURL)
		} else {
			assert.Nil(t, res.GetOptions().Capacity, tc.typeURL)
		}
	}
}

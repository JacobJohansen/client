package libkb

import (
	"testing"

	keybase1 "github.com/keybase/client/go/protocol/keybase1"
	"github.com/stretchr/testify/require"
)

func TestProveTypes(t *testing.T) {
	// Make sure we don't accidentally exclude any new proof types.
	excludedServiceTypes := map[keybase1.ProofType]bool{
		keybase1.ProofType_NONE: true,
		keybase1.ProofType_PGP:  true,
	}

	excludedServiceOrderTypes := map[keybase1.ProofType]bool{
		keybase1.ProofType_NONE: true,
		keybase1.ProofType_DNS:  true,
		keybase1.ProofType_PGP:  true,
	}

	remoteServiceTypesRevMap := map[keybase1.ProofType]bool{}
	for _, pt := range RemoteServiceTypes {
		remoteServiceTypesRevMap[pt] = true
	}

	remoteServiceOrderMap := map[keybase1.ProofType]bool{}
	for _, pt := range RemoteServiceOrder {
		remoteServiceOrderMap[pt] = true
	}

	for _, pt := range keybase1.ProofTypeMap {
		if _, ok := remoteServiceTypesRevMap[pt]; !ok {
			_, isExcluded := excludedServiceTypes[pt]
			require.True(t, isExcluded, "%v is missing from RemoteServiceTypes", pt)
		}
		if _, ok := remoteServiceOrderMap[pt]; !ok {
			_, isExcluded := excludedServiceOrderTypes[pt]
			require.True(t, isExcluded, "%v is missing from RemoteServiceOrder", pt)
		}
	}
}

type MockedConfig struct {
	NullConfiguration
	mocked_is_ssl_pinning_enabled bool
}
func (mc MockedConfig) IsSSLPinningEnabled() bool {
	return mc.mocked_is_ssl_pinning_enabled
}

func TestServerLookup(t *testing.T) {
	server, err := ServerLookup(NewEnv(nil, nil, makeLogGetter(t)), DevelRunMode)
	require.Equal(t, DevelServerURI, server)
	require.Equal(t, nil, err)

	server, err = ServerLookup(NewEnv(nil, nil, makeLogGetter(t)), StagingRunMode)
	require.Equal(t, StagingServerURI, server)
	require.Equal(t, nil, err)

	server, err = ServerLookup(NewEnv(nil, nil, makeLogGetter(t)), StagingRunMode)
	require.Equal(t, StagingServerURI, server)
	require.Equal(t, nil, err)

	server, err = ServerLookup(NewEnv(MockedConfig{NullConfiguration{}, true}, nil, makeLogGetter(t)), ProductionRunMode)
	require.Equal(t, ProductionServerURI, server)
	require.Equal(t, nil, err)

	server, err = ServerLookup(NewEnv(MockedConfig{NullConfiguration{}, false}, nil, makeLogGetter(t)), ProductionRunMode)
	require.Equal(t, ProductionSiteURI, server)
	require.Equal(t, nil, err)

	server, err = ServerLookup(NewEnv(MockedConfig{NullConfiguration{}, false}, nil, makeLogGetter(t)), NoRunMode)
	require.Equal(t, "", server)
	require.NotEqual(t, nil, err)
}

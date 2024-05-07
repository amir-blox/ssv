package qbft

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/bloxapp/ssv/logging"
	"github.com/bloxapp/ssv/protocol/v2/genesisqbft/instance"
	"github.com/bloxapp/ssv/protocol/v2/genesisqbft/roundtimer"
	"github.com/stretchr/testify/require"

	"github.com/ssvlabs/ssv-spec-pre-cc/qbft"
	genesisspectypes "github.com/ssvlabs/ssv-spec-pre-cc/types"
	"github.com/ssvlabs/ssv-spec-pre-cc/types/testingutils"
)

type SpecTest struct {
	Name               string
	Pre                *instance.Instance
	PostRoot           string
	OutputMessages     []*qbft.SignedMessage
	ExpectedTimerState *testingutils.TimerState
	ExpectedError      string
}

func RunTimeout(t *testing.T, test *SpecTest) {
	logger := logging.TestLogger(t)
	err := test.Pre.UponRoundTimeout(logger)

	if len(test.ExpectedError) != 0 {
		require.EqualError(t, err, test.ExpectedError)
	} else {
		require.NoError(t, err)
	}

	// test calling timeout
	timer, ok := test.Pre.GetConfig().GetTimer().(*roundtimer.TestQBFTTimer)
	require.True(t, ok)
	require.Equal(t, test.ExpectedTimerState.Timeouts, timer.State.Timeouts)
	require.Equal(t, test.ExpectedTimerState.Round, timer.State.Round)

	// test output message
	broadcastedMsgs := test.Pre.GetConfig().GetNetwork().(*testingutils.TestingNetwork).BroadcastedMsgs
	if len(test.OutputMessages) > 0 || len(broadcastedMsgs) > 0 {
		require.Len(t, broadcastedMsgs, len(test.OutputMessages))

		for i, msg := range test.OutputMessages {
			r1, _ := msg.GetRoot()

			ssvMsg := &genesisspectypes.SSVMessage{}
			require.NoError(t, ssvMsg.Decode(broadcastedMsgs[i].Data))

			msg2 := &qbft.SignedMessage{}
			require.NoError(t, msg2.Decode(ssvMsg.Data))
			r2, _ := msg2.GetRoot()

			require.EqualValuesf(t, r1, r2, fmt.Sprintf("output msg %d roots not equal", i))
		}
	}

	postRoot, err := test.Pre.State.GetRoot()
	require.NoError(t, err)
	require.EqualValuesf(t, test.PostRoot, hex.EncodeToString(postRoot[:]), "post root not valid")
}

package inferencesynthesis_test

import (
	"context"
	"fmt"

	alloraMath "github.com/allora-network/allora-chain/math"
	"github.com/allora-network/allora-chain/test/testutil"
	inferencesynthesis "github.com/allora-network/allora-chain/x/emissions/keeper/inference_synthesis"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
)

func (s *InferenceSynthesisTestSuite) TestConvertValueBundleToNetworkLossesByWorker() {
	require := s.Require()
	valueBundle := emissionstypes.ValueBundle{
		TopicId: uint64(1),
		Reputer: s.addrsStr[1],
		ReputerRequestNonce: &emissionstypes.ReputerRequestNonce{
			ReputerNonce: &emissionstypes.Nonce{BlockHeight: 100},
		},
		ExtraData:     nil,
		CombinedValue: alloraMath.MustNewDecFromString("0.1"),
		NaiveValue:    alloraMath.MustNewDecFromString("0.1"),
		InfererValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: s.addrsStr[1], Value: alloraMath.MustNewDecFromString("0.1")},
			{Worker: s.addrsStr[2], Value: alloraMath.MustNewDecFromString("0.2")},
		},
		ForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: s.addrsStr[1], Value: alloraMath.MustNewDecFromString("0.1")},
			{Worker: s.addrsStr[2], Value: alloraMath.MustNewDecFromString("0.2")},
		},
		OneOutInfererValues: []*emissionstypes.WithheldWorkerAttributedValue{
			{Worker: s.addrsStr[1], Value: alloraMath.MustNewDecFromString("0.1")},
			{Worker: s.addrsStr[2], Value: alloraMath.MustNewDecFromString("0.2")},
		},
		OneOutForecasterValues: []*emissionstypes.WithheldWorkerAttributedValue{
			{Worker: s.addrsStr[1], Value: alloraMath.MustNewDecFromString("0.1")},
			{Worker: s.addrsStr[2], Value: alloraMath.MustNewDecFromString("0.2")},
		},
		OneInForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: s.addrsStr[1], Value: alloraMath.MustNewDecFromString("0.1")},
			{Worker: s.addrsStr[2], Value: alloraMath.MustNewDecFromString("0.2")},
		},
		OneOutInfererForecasterValues: nil,
	}

	result := inferencesynthesis.ConvertValueBundleToNetworkLossesByWorker(valueBundle)

	// Check if CombinedLoss and NaiveLoss are correctly set
	require.Equal(alloraMath.MustNewDecFromString("0.1"), result.CombinedLoss)
	require.Equal(alloraMath.MustNewDecFromString("0.1"), result.NaiveLoss)

	// Check if each worker's losses are set correctly
	expectedLoss := alloraMath.MustNewDecFromString("0.1")
	expectedLoss2 := alloraMath.MustNewDecFromString("0.2")
	require.Equal(expectedLoss, result.InfererLosses[s.addrsStr[1]])
	require.Equal(expectedLoss2, result.InfererLosses[s.addrsStr[2]])
	require.Equal(expectedLoss, result.ForecasterLosses[s.addrsStr[1]])
	require.Equal(expectedLoss2, result.ForecasterLosses[s.addrsStr[2]])
	require.Equal(expectedLoss, result.OneOutInfererLosses[s.addrsStr[1]])
	require.Equal(expectedLoss2, result.OneOutInfererLosses[s.addrsStr[2]])
	require.Equal(expectedLoss, result.OneOutForecasterLosses[s.addrsStr[1]])
	require.Equal(expectedLoss2, result.OneOutForecasterLosses[s.addrsStr[2]])
	require.Equal(expectedLoss, result.OneInForecasterLosses[s.addrsStr[1]])
	require.Equal(expectedLoss2, result.OneInForecasterLosses[s.addrsStr[2]])
}

func (s *InferenceSynthesisTestSuite) TestComputeAndBuildEMRegret() {
	require := s.Require()

	alpha := alloraMath.MustNewDecFromString("0.1")
	lossA := alloraMath.MustNewDecFromString("500")
	lossB := alloraMath.MustNewDecFromString("200")
	previous := alloraMath.MustNewDecFromString("200")

	blockHeight := int64(123)

	result, err := inferencesynthesis.ComputeAndBuildEMRegret(lossA, lossB, previous, alpha, blockHeight)
	require.NoError(err)

	expected, err := alloraMath.NewDecFromString("210")
	require.NoError(err)

	require.True(alloraMath.InDelta(expected, result.Value, alloraMath.MustNewDecFromString("0.0001")))
	require.Equal(blockHeight, result.BlockHeight)
}

// TestGetCalcSetNetworkRegretsTwoWorkers tests the GetCalcSetNetworkRegrets function
// with two workers in a simplified scenario.
//
// Setup:
// - Create a topic with ID 1 and initial regret of 0
// - Set AlphaRegret to 0.5, making the experience threshold 2 inclusions
// - Define three workers, but only use two in the value bundle
// - Set up a value bundle with combined value 500 and individual values of 200
//
// Expected outcomes:
//  1. The function should execute without error
//  2. Network regrets should be calculated and set for both workers
//  3. The topic's initial regret should be updated from 0
//  4. Regrets for both workers should be equal, as they have the same values
//  5. The calculated regrets should reflect the difference between individual
//     and combined values, influenced by the AlphaRegret parameter
//
// This test ensures that the regret calculation works correctly for a simple
// case with two equally performing workers, and that the topic's initial
// regret is properly updated. The workers should have new regrets informed
// by the topic's initial regret.
func (s *InferenceSynthesisTestSuite) TestGetCalcSetNetworkRegretsTwoWorkers() {
	require := s.Require()
	k := s.emissionsKeeper

	topicId := uint64(1)
	// Create new topic
	topic := s.mockTopic()
	topic.InitialRegret = alloraMath.ZeroDec()
	// Need to use "0.5" to set limit inclusions count as 2=(1/0.5)
	topic.AlphaRegret = alloraMath.MustNewDecFromString("0.5")
	err := s.emissionsKeeper.SetTopic(s.ctx, topicId, topic)
	require.NoError(err)

	worker1 := s.addrsStr[1]
	worker2 := s.addrsStr[2]
	worker3 := s.addrsStr[3]

	pNorm := alloraMath.MustNewDecFromString("0.1")
	cNorm := alloraMath.MustNewDecFromString("0.1")
	epsilon := alloraMath.MustNewDecFromString("0.0001")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	blockHeight := int64(42)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	reputerRequestNonce := emissionstypes.ReputerRequestNonce{
		ReputerNonce: &nonce,
	}
	valueBundle := emissionstypes.ValueBundle{
		TopicId:             topicId,
		ReputerRequestNonce: &reputerRequestNonce,
		Reputer:             s.addrsStr[9],
		ExtraData:           nil,
		CombinedValue:       alloraMath.NewDecFromInt64(500),
		InfererValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
		},
		ForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
		},
		NaiveValue:          alloraMath.NewDecFromInt64(123),
		OneOutInfererValues: nil,
		OneInForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
		},
		OneOutForecasterValues:        nil,
		OneOutInfererForecasterValues: nil,
	}

	regretVal := emissionstypes.TimestampedValue{
		BlockHeight: blockHeight,
		Value:       alloraMath.NewDecFromInt64(200),
	}

	// Need to more than 2 experienced actor
	// For this need to call SetInfererNetwork, SetForecasterNetworkRegret for worker1, worker2
	err = k.SetInfererNetworkRegret(s.ctx, topicId, worker1, regretVal)
	require.NoError(err)
	err = k.SetInfererNetworkRegret(s.ctx, topicId, worker2, regretVal)
	require.NoError(err)
	err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker1, regretVal)
	require.NoError(err)
	err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker2, regretVal)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker1, worker1, regretVal)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker1, worker2, regretVal)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker2, worker1, regretVal)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker2, worker2, regretVal)
	require.NoError(err)

	s.incrementRegretsInTopic(topicId, worker1, 2, emissionstypes.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED)
	s.incrementRegretsInTopic(topicId, worker2, 2, emissionstypes.ActorType_ACTOR_TYPE_INFERER_UNSPECIFIED)
	// New potential participant should start with zero regret at this point since the initial regret in the topic is zero
	// It will be updated after the first regret calculation
	worker3LastRegret, worker3NoPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker3)
	require.NoError(err)
	require.Equal(worker3LastRegret.Value, alloraMath.ZeroDec())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetForecasterNetworkRegret(s.ctx, topicId, worker3)
	require.NoError(err)
	require.Equal(worker3LastRegret.Value, alloraMath.ZeroDec())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker1)
	require.NoError(err)
	require.Equal(worker3LastRegret.Value, alloraMath.ZeroDec())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker2)
	require.NoError(err)
	require.Equal(worker3LastRegret.Value, alloraMath.ZeroDec())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker3)
	require.NoError(err)
	require.Equal(worker3LastRegret.Value, alloraMath.ZeroDec())
	require.True(worker3NoPriorRegret)

	err = inferencesynthesis.GetCalcSetNetworkRegrets(
		inferencesynthesis.GetCalcSetNetworkRegretsArgs{
			Ctx:                   s.ctx,
			K:                     s.emissionsKeeper,
			TopicId:               topicId,
			NetworkLosses:         valueBundle,
			Nonce:                 nonce,
			AlphaRegret:           topic.AlphaRegret,
			CNorm:                 cNorm,
			PNorm:                 pNorm,
			EpsilonTopic:          epsilon,
			InitialRegretQuantile: initialRegretQuantile,
			PNormSafeDiv:          pnormSafeDiv,
		})
	require.NoError(err)

	bothAccs := []string{worker1, worker2}

	// New potential participant should not start with zero regret since we already have participants with prior regrets which will
	// be used to calculate the initial regret in the topic
	worker3LastRegret, worker3NoPriorRegret, err = k.GetInfererNetworkRegret(s.ctx, topicId, worker3)
	require.NoError(err)
	require.NotEqual(worker3LastRegret.Value.String(), alloraMath.ZeroDec().String())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetForecasterNetworkRegret(s.ctx, topicId, worker3)
	require.NoError(err)
	require.NotEqual(worker3LastRegret.Value.String(), alloraMath.ZeroDec().String())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker1)
	require.NoError(err)
	require.NotEqual(worker3LastRegret.Value.String(), alloraMath.ZeroDec().String())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker2)
	require.NoError(err)
	require.NotEqual(worker3LastRegret.Value.String(), alloraMath.ZeroDec().String())
	require.True(worker3NoPriorRegret)

	worker3LastRegret, worker3NoPriorRegret, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker3)
	require.NoError(err)
	require.NotEqual(worker3LastRegret.Value.String(), alloraMath.ZeroDec().String())
	require.True(worker3NoPriorRegret)

	// Get topic initial regret
	topic, err = k.GetTopic(s.ctx, topicId)
	require.NoError(err)

	for _, acc := range bothAccs {
		lastRegret, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, acc)
		require.NoError(err)
		require.True(alloraMath.InDelta(topic.InitialRegret, lastRegret.Value, alloraMath.MustNewDecFromString("0.001")))
		require.False(noPriorRegret)

		lastRegret, noPriorRegret, err = k.GetForecasterNetworkRegret(s.ctx, topicId, acc)
		require.NoError(err)
		require.True(alloraMath.InDelta(topic.InitialRegret, lastRegret.Value, alloraMath.MustNewDecFromString("0.001")))
		require.False(noPriorRegret)

		for _, accInner := range bothAccs {
			lastRegret, _, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, acc, accInner)
			require.NoError(err)
		}
	}
}

func (s *InferenceSynthesisTestSuite) TestGetCalcSetNetworkRegretsThreeWorkers() {
	require := s.Require()
	k := s.emissionsKeeper

	worker1 := s.addrsStr[1]
	worker2 := s.addrsStr[2]
	worker3 := s.addrsStr[3]

	pNorm := alloraMath.MustNewDecFromString("0.1")
	cNorm := alloraMath.MustNewDecFromString("0.1")
	epsilon := alloraMath.MustNewDecFromString("0.0001")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	valueBundle := emissionstypes.ValueBundle{
		TopicId: uint64(1),
		Reputer: s.addrsStr[1],
		ReputerRequestNonce: &emissionstypes.ReputerRequestNonce{
			ReputerNonce: &emissionstypes.Nonce{BlockHeight: 100},
		},
		ExtraData:     nil,
		CombinedValue: alloraMath.MustNewDecFromString("500"),
		NaiveValue:    alloraMath.MustNewDecFromString("123"),
		InfererValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker3, Value: alloraMath.MustNewDecFromString("200")},
		},
		ForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker3, Value: alloraMath.MustNewDecFromString("200")},
		},
		OneInForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("200")},
			{Worker: worker3, Value: alloraMath.MustNewDecFromString("200")},
		},
		OneOutInfererValues:           nil,
		OneOutForecasterValues:        nil,
		OneOutInfererForecasterValues: nil,
	}
	blockHeight := int64(42)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	alpha := alloraMath.MustNewDecFromString("0.1")
	topicId := uint64(1)

	timestampedValue := emissionstypes.TimestampedValue{
		BlockHeight: blockHeight,
		Value:       alloraMath.MustNewDecFromString("200"),
	}

	err := k.SetInfererNetworkRegret(s.ctx, topicId, worker1, timestampedValue)
	require.NoError(err)
	err = k.SetInfererNetworkRegret(s.ctx, topicId, worker2, timestampedValue)
	require.NoError(err)
	err = k.SetInfererNetworkRegret(s.ctx, topicId, worker3, timestampedValue)
	require.NoError(err)

	err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker1, timestampedValue)
	require.NoError(err)
	err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker2, timestampedValue)
	require.NoError(err)
	err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker2, timestampedValue)
	require.NoError(err)

	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker1, worker1, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker1, worker2, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker1, worker3, timestampedValue)
	require.NoError(err)

	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker2, worker1, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker2, worker2, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker2, worker3, timestampedValue)
	require.NoError(err)

	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker1, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker2, timestampedValue)
	require.NoError(err)
	err = k.SetOneInForecasterNetworkRegret(s.ctx, topicId, worker3, worker3, timestampedValue)
	require.NoError(err)

	err = inferencesynthesis.GetCalcSetNetworkRegrets(
		inferencesynthesis.GetCalcSetNetworkRegretsArgs{
			Ctx:                   s.ctx,
			K:                     s.emissionsKeeper,
			TopicId:               topicId,
			NetworkLosses:         valueBundle,
			Nonce:                 nonce,
			AlphaRegret:           alpha,
			CNorm:                 cNorm,
			PNorm:                 pNorm,
			EpsilonTopic:          epsilon,
			InitialRegretQuantile: initialRegretQuantile,
			PNormSafeDiv:          pnormSafeDiv,
		})
	require.NoError(err)

	allWorkerAccs := []string{worker1, worker2, worker3}
	expected := alloraMath.MustNewDecFromString("210")
	// expectedOneIn := alloraMath.MustNewDecFromString("180")

	for _, workerAcc := range allWorkerAccs {
		lastRegret, _, err := k.GetInfererNetworkRegret(s.ctx, topicId, workerAcc)
		require.NoError(err)
		require.True(alloraMath.InDelta(expected, lastRegret.Value, alloraMath.MustNewDecFromString("0.0001")))

		lastRegret, _, err = k.GetForecasterNetworkRegret(s.ctx, topicId, workerAcc)
		require.NoError(err)

		for _, innerWorkerAcc := range allWorkerAccs {
			lastRegret, _, err = k.GetOneInForecasterNetworkRegret(s.ctx, topicId, workerAcc, innerWorkerAcc)
			require.NoError(err)
		}
	}
}

func (s *InferenceSynthesisTestSuite) TestGetCalcSetNetworkRegretsFromCsv() {
	require := s.Require()
	k := s.emissionsKeeper
	epochGet := testutil.GetSimulatedValuesGetterForEpochs()
	epochPrevGet := epochGet[300]
	epoch301Get := epochGet[301]
	topicId := uint64(1)
	blockHeight := int64(1003)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	alpha := alloraMath.MustNewDecFromString("0.1")
	pNorm := alloraMath.MustNewDecFromString("3.0")
	cNorm := alloraMath.MustNewDecFromString("0.75")
	epsilon := alloraMath.MustNewDecFromString("1e-4")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	inferer0 := s.addrs[0].String()
	inferer1 := s.addrs[1].String()
	inferer2 := s.addrs[2].String()
	inferer3 := s.addrs[3].String()
	inferer4 := s.addrs[4].String()
	infererAddresses := []string{inferer0, inferer1, inferer2, inferer3, inferer4}

	forecaster0 := s.addrs[5].String()
	forecaster1 := s.addrs[6].String()
	forecaster2 := s.addrs[7].String()
	forecasterAddresses := []string{forecaster0, forecaster1, forecaster2}

	reputer0 := s.addrs[8].String()

	err := testutil.SetRegretsFromPreviousEpoch(s.ctx, s.emissionsKeeper, topicId, blockHeight, infererAddresses, forecasterAddresses, epochPrevGet)
	require.NoError(err)

	networkLosses, err := testutil.GetNetworkLossFromCsv(
		topicId,
		blockHeight,
		infererAddresses,
		forecasterAddresses,
		reputer0,
		epoch301Get,
	)
	s.Require().NoError(err)

	err = inferencesynthesis.GetCalcSetNetworkRegrets(
		inferencesynthesis.GetCalcSetNetworkRegretsArgs{
			Ctx:                   s.ctx,
			K:                     s.emissionsKeeper,
			TopicId:               topicId,
			NetworkLosses:         networkLosses,
			Nonce:                 nonce,
			AlphaRegret:           alpha,
			CNorm:                 cNorm,
			PNorm:                 pNorm,
			EpsilonTopic:          epsilon,
			InitialRegretQuantile: initialRegretQuantile,
			PNormSafeDiv:          pnormSafeDiv,
		})
	require.NoError(err)

	checkRegret := func(worker string, expected alloraMath.Dec, getter func(context.Context, uint64, string) (emissionstypes.TimestampedValue, bool, error)) {
		regret, _, err := getter(s.ctx, topicId, worker)
		require.NoError(err)
		testutil.InEpsilon5(s.T(), expected, regret.Value.String())
	}

	checkOneOutRegret := func(worker string, innerWorker string, expected alloraMath.Dec, getter func(context.Context, uint64, string, string) (emissionstypes.TimestampedValue, bool, error)) {
		regret, _, err := getter(s.ctx, topicId, worker, innerWorker)
		require.NoError(err)
		testutil.InEpsilon5(s.T(), expected, regret.Value.String())
	}

	for i := 0; i < len(infererAddresses); i++ {
		expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v", i))
		checkRegret(infererAddresses[i], expectedRegret, k.GetInfererNetworkRegret)

		expectedRegret = epoch301Get(fmt.Sprintf("naive_inference_regret_worker_%v", i))
		checkRegret(infererAddresses[i], expectedRegret, k.GetNaiveInfererNetworkRegret)
	}

	for i := 0; i < len(forecasterAddresses); i++ {
		forecasterCsvIndex := i + 5
		expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v", forecasterCsvIndex))
		checkRegret(forecasterAddresses[i], expectedRegret, k.GetForecasterNetworkRegret)
	}

	for i, inferer := range infererAddresses {
		for j, infererInner := range infererAddresses {
			expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_oneout_%v", j, i))
			checkOneOutRegret(inferer, infererInner, expectedRegret, k.GetOneOutInfererInfererNetworkRegret)
		}

		for l, forecaster := range forecasterAddresses {
			forecasterCsvIndex := l + 5
			expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_oneout_%v", forecasterCsvIndex, i))
			checkOneOutRegret(inferer, forecaster, expectedRegret, k.GetOneOutInfererForecasterNetworkRegret)
		}
	}

	for i, forecaster := range forecasterAddresses {
		forecasterCsvIndex := i + 5
		for j, inferer := range infererAddresses {
			expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_oneout_%v", j, forecasterCsvIndex))
			checkOneOutRegret(forecaster, inferer, expectedRegret, k.GetOneOutForecasterInfererNetworkRegret)
		}

		for z, forecasterInner := range forecasterAddresses {
			forecasterCsvIndex2 := z + 5
			expectedRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_oneout_%v", forecasterCsvIndex2, forecasterCsvIndex))
			checkOneOutRegret(forecaster, forecasterInner, expectedRegret, k.GetOneOutForecasterForecasterNetworkRegret)
		}

		expectedOneInRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_onein_%v", 5, i))
		checkOneOutRegret(forecaster, forecaster, expectedOneInRegret, k.GetOneInForecasterNetworkRegret)

		for l, inferer := range infererAddresses {
			expectedOneInRegret := epoch301Get(fmt.Sprintf("inference_regret_worker_%v_onein_%v", l, i))
			checkOneOutRegret(forecaster, inferer, expectedOneInRegret, k.GetOneInForecasterNetworkRegret)
		}
	}
}

// In this test we run two trials of calculating setting network regrets with different losses.
// We then compare the resulting regrets to see if the higher losses result in lower regrets.
func (s *InferenceSynthesisTestSuite) TestHigherLossesLowerRegret() {
	require := s.Require()
	k := s.emissionsKeeper

	topicId := uint64(1)
	blockHeight := int64(1003)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	alpha := alloraMath.MustNewDecFromString("0.1")
	pNorm := alloraMath.MustNewDecFromString("0.1")
	cNorm := alloraMath.MustNewDecFromString("0.1")
	epsilon := alloraMath.MustNewDecFromString("0.0001")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	worker0 := s.addrsStr[0]
	worker1 := s.addrsStr[1]
	worker2 := s.addrsStr[2]

	networkLossesValueBundle0 := emissionstypes.ValueBundle{
		TopicId:             topicId,
		ReputerRequestNonce: &emissionstypes.ReputerRequestNonce{ReputerNonce: &nonce},
		Reputer:             s.addrsStr[9],
		ExtraData:           nil,
		CombinedValue:       alloraMath.MustNewDecFromString("0.1"),
		InfererValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.2")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		ForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.2")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		NaiveValue:             alloraMath.MustNewDecFromString("0.1"),
		OneOutInfererValues:    nil,
		OneOutForecasterValues: nil,
		OneInForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.2")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		OneOutInfererForecasterValues: nil,
	}

	networkLossesValueBundle1 := emissionstypes.ValueBundle{
		TopicId:             topicId,
		ReputerRequestNonce: &emissionstypes.ReputerRequestNonce{ReputerNonce: &nonce},
		Reputer:             s.addrsStr[9],
		ExtraData:           nil,
		CombinedValue:       alloraMath.MustNewDecFromString("0.1"),
		InfererValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		ForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		NaiveValue:             alloraMath.MustNewDecFromString("0.1"),
		OneOutInfererValues:    nil,
		OneOutForecasterValues: nil,
		OneInForecasterValues: []*emissionstypes.WorkerAttributedValue{
			{Worker: worker0, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker1, Value: alloraMath.MustNewDecFromString("0.3")},
			{Worker: worker2, Value: alloraMath.MustNewDecFromString("0.4")},
		},
		OneOutInfererForecasterValues: nil,
	}

	resetRegrets := func() {
		timestampedValue0_1 := emissionstypes.TimestampedValue{
			BlockHeight: blockHeight,
			Value:       alloraMath.MustNewDecFromString("0.1"),
		}

		timestampedValue0_2 := emissionstypes.TimestampedValue{
			BlockHeight: blockHeight,
			Value:       alloraMath.MustNewDecFromString("0.2"),
		}

		timestampedValue0_3 := emissionstypes.TimestampedValue{
			BlockHeight: blockHeight,
			Value:       alloraMath.MustNewDecFromString("0.3"),
		}

		err := k.SetInfererNetworkRegret(s.ctx, topicId, worker0, timestampedValue0_1)
		require.NoError(err)
		err = k.SetInfererNetworkRegret(s.ctx, topicId, worker1, timestampedValue0_2)
		require.NoError(err)
		err = k.SetInfererNetworkRegret(s.ctx, topicId, worker2, timestampedValue0_3)
		require.NoError(err)

		err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker0, timestampedValue0_1)
		require.NoError(err)
		err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker1, timestampedValue0_2)
		require.NoError(err)
		err = k.SetForecasterNetworkRegret(s.ctx, topicId, worker2, timestampedValue0_3)
		require.NoError(err)
	}

	// Test 0

	resetRegrets()

	err := inferencesynthesis.GetCalcSetNetworkRegrets(
		inferencesynthesis.GetCalcSetNetworkRegretsArgs{
			Ctx:                   s.ctx,
			K:                     s.emissionsKeeper,
			TopicId:               topicId,
			NetworkLosses:         networkLossesValueBundle0,
			Nonce:                 nonce,
			AlphaRegret:           alpha,
			CNorm:                 cNorm,
			PNorm:                 pNorm,
			EpsilonTopic:          epsilon,
			InitialRegretQuantile: initialRegretQuantile,
			PNormSafeDiv:          pnormSafeDiv,
		})
	require.NoError(err)

	// Record resulting regrets

	infererRegret0_0, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker0)
	require.NoError(err)
	require.False(noPriorRegret)
	infererRegret0_1, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker1)
	require.NoError(err)
	require.False(noPriorRegret)
	infererRegret0_2, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker2)
	require.NoError(err)
	require.False(noPriorRegret)

	forecasterRegret0_0, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker0)
	require.NoError(err)
	require.False(noPriorRegret)
	forecasterRegret0_1, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker1)
	require.NoError(err)
	require.False(noPriorRegret)
	forecasterRegret0_2, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker2)
	require.NoError(err)
	require.False(noPriorRegret)

	// Test 1

	resetRegrets()

	err = inferencesynthesis.GetCalcSetNetworkRegrets(
		inferencesynthesis.GetCalcSetNetworkRegretsArgs{
			Ctx:                   s.ctx,
			K:                     s.emissionsKeeper,
			TopicId:               topicId,
			NetworkLosses:         networkLossesValueBundle1,
			Nonce:                 nonce,
			AlphaRegret:           alpha,
			CNorm:                 cNorm,
			PNorm:                 pNorm,
			EpsilonTopic:          epsilon,
			InitialRegretQuantile: initialRegretQuantile,
			PNormSafeDiv:          pnormSafeDiv,
		})
	require.NoError(err)

	// Record resulting regrets

	infererRegret1_0, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker0)
	require.NoError(err)
	require.False(noPriorRegret)
	infererRegret1_1, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker1)
	require.NoError(err)
	require.False(noPriorRegret)
	infererRegret1_2, noPriorRegret, err := k.GetInfererNetworkRegret(s.ctx, topicId, worker2)
	require.NoError(err)
	require.False(noPriorRegret)

	forecasterRegret1_0, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker0)
	require.NoError(err)
	require.False(noPriorRegret)
	forecasterRegret1_1, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker1)
	require.NoError(err)
	require.False(noPriorRegret)
	forecasterRegret1_2, noPriorRegret, err := k.GetForecasterNetworkRegret(s.ctx, topicId, worker2)
	require.NoError(err)
	require.False(noPriorRegret)

	// Test

	require.True(infererRegret0_0.Value.Gt(infererRegret1_0.Value))
	require.Equal(infererRegret0_1.Value, infererRegret1_1.Value)
	require.Equal(infererRegret0_2.Value, infererRegret1_2.Value)

	require.True(forecasterRegret0_0.Value.Gt(forecasterRegret1_0.Value))
	require.Equal(forecasterRegret0_1.Value, forecasterRegret1_1.Value)
	require.Equal(forecasterRegret0_2.Value, forecasterRegret1_2.Value)
}

func (s *InferenceSynthesisTestSuite) TestCalcTopicInitialRegret() {
	require := s.Require()

	regrets := []alloraMath.Dec{
		alloraMath.MustNewDecFromString("0.6445506208021189"),
		alloraMath.MustNewDecFromString("1.0216386898413485"),
		alloraMath.MustNewDecFromString("0.6092049398135028"),
		alloraMath.MustNewDecFromString("0.6971588004566455"),
		alloraMath.MustNewDecFromString("0.9030751421888253"),
		alloraMath.MustNewDecFromString("0.8219035038858344"),
	}
	cNorm := alloraMath.MustNewDecFromString("0.75")
	pNorm := alloraMath.MustNewDecFromString("3.0")
	epsilon := alloraMath.MustNewDecFromString("0.0001")
	percentileRegert := alloraMath.MustNewDecFromString("0.25")
	pnormDiv := alloraMath.MustNewDecFromString("8.25")
	calculatedInitialRegret, err := inferencesynthesis.CalcTopicInitialRegret(
		regrets,
		epsilon,
		pNorm,
		cNorm,
		percentileRegert,
		pnormDiv,
	)
	require.NoError(err)
	testutil.InEpsilon5(s.T(), calculatedInitialRegret, "0.3354820760526412097325669544281814")
}

// TestUpdateTopicInitialRegret tests the UpdateTopicInitialRegret function.
//
// Setup:
// - Create a topic with ID 1 and initial regret of 0
// - Set AlphaRegret to 0.5, making the experience threshold 2 inclusions
// - Add 5 inferers and 3 forecasters, each with 2 inclusions to make them experienced
// - Set up a simulated value getter for epochs 300 and 301
//
// Test steps:
// 1. Create a value bundle with combined value and individual values for inferers and forecasters
// 2. Call UpdateTopicInitialRegret with this value bundle
// 3. Retrieve the updated topic
//
// Expected outcomes:
// 1. The function should execute without error
// 2. The topic's initial regret should be updated from 0 to a non-zero value
// 3. The new initial regret should be calculated based on the provided values and parameters
//
// This test ensures that the UpdateTopicInitialRegret function correctly calculates
// and updates the initial regret for a topic based on the performance of experienced
// workers, using the provided normalization and calculation parameters.
func (s *InferenceSynthesisTestSuite) TestUpdateTopicInitialRegret() {
	require := s.Require()
	k := s.emissionsKeeper
	epochGet := testutil.GetSimulatedValuesGetterForEpochs()
	epochPrevGet := epochGet[300]
	epoch301Get := epochGet[301]

	topicId := uint64(1)
	blockHeight := int64(1003)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	alpha := alloraMath.MustNewDecFromString("0.1")
	pNorm := alloraMath.MustNewDecFromString("3.0")
	cNorm := alloraMath.MustNewDecFromString("0.75")
	epsilon := alloraMath.MustNewDecFromString("1e-4")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	// Set initial Regret to check if this value is updated or not
	initialRegret := alloraMath.MustNewDecFromString("0")
	topic := s.mockTopic()
	// Need to use "0.5" to set limit inclusions count as 2=(1/0.5)
	topic.AlphaRegret = alloraMath.MustNewDecFromString("0.5")
	// Create new topic
	err := s.emissionsKeeper.SetTopic(s.ctx, topicId, topic)
	s.Require().NoError(err)

	inferer0 := s.addrs[0].String()
	inferer1 := s.addrs[1].String()
	inferer2 := s.addrs[2].String()
	inferer3 := s.addrs[3].String()
	inferer4 := s.addrs[4].String()
	infererAddresses := []string{inferer0, inferer1, inferer2, inferer3, inferer4}

	forecaster0 := s.addrs[5].String()
	forecaster1 := s.addrs[6].String()
	forecaster2 := s.addrs[7].String()
	forecasterAddresses := []string{forecaster0, forecaster1, forecaster2}

	reputer0 := s.addrs[8].String()

	// Need to add experienced inferers for this topic
	for _, worker := range infererAddresses {
		err := k.IncrementCountInfererInclusionsInTopic(s.ctx, topicId, worker)
		require.NoError(err)
		err = k.IncrementCountInfererInclusionsInTopic(s.ctx, topicId, worker)
		require.NoError(err)
	}

	// Need to add experienced forecasters for this topic
	for _, worker := range forecasterAddresses {
		err := k.IncrementCountForecasterInclusionsInTopic(s.ctx, topicId, worker)
		require.NoError(err)
		err = k.IncrementCountForecasterInclusionsInTopic(s.ctx, topicId, worker)
		require.NoError(err)
	}

	err = testutil.SetRegretsFromPreviousEpoch(s.ctx, s.emissionsKeeper, topicId, blockHeight, infererAddresses, forecasterAddresses, epochPrevGet)
	require.NoError(err)

	networkLosses, err := testutil.GetNetworkLossFromCsv(
		topicId,
		blockHeight,
		infererAddresses,
		forecasterAddresses,
		reputer0,
		epoch301Get,
	)
	s.Require().NoError(err)

	err = inferencesynthesis.GetCalcSetNetworkRegrets(inferencesynthesis.GetCalcSetNetworkRegretsArgs{
		Ctx:                   s.ctx,
		K:                     k,
		TopicId:               topicId,
		NetworkLosses:         networkLosses,
		Nonce:                 nonce,
		AlphaRegret:           alpha,
		CNorm:                 cNorm,
		PNorm:                 pNorm,
		EpsilonTopic:          epsilon,
		InitialRegretQuantile: initialRegretQuantile,
		PNormSafeDiv:          pnormSafeDiv,
	})
	require.NoError(err)

	// Assert that initial regret is updated
	topic, err = s.emissionsKeeper.GetTopic(s.ctx, topicId)
	require.NoError(err)
	require.NotEqual(topic.InitialRegret, initialRegret)
}

// TestNotUpdateTopicInitialRegret tests that the topic's initial regret is not updated
// when there are no experienced workers for the topic.
//
// Setup:
// - Create a topic with an initial regret of 0
// - Set up inferer and forecaster addresses
// - Do not increment the inclusion counts for workers, leaving them as inexperienced
// - Set regrets from the previous epoch
// - Generate network losses for the current epoch
//
// Expected outcomes:
//  1. The GetCalcSetNetworkRegrets function should execute without error
//  2. The topic's initial regret should remain unchanged (0) after the function call
//     because there are no experienced workers to trigger an update
//
// This test ensures that the initial regret of a topic is not modified when
// there are no experienced workers, maintaining the system's integrity by
// preventing premature updates to the topic's regret baseline.
func (s *InferenceSynthesisTestSuite) TestNotUpdateTopicInitialRegret() {
	require := s.Require()
	k := s.emissionsKeeper
	epochGet := testutil.GetSimulatedValuesGetterForEpochs()
	epochPrevGet := epochGet[300]
	epoch301Get := epochGet[301]

	topicId := uint64(1)
	blockHeight := int64(1003)
	nonce := emissionstypes.Nonce{BlockHeight: blockHeight}
	alpha := alloraMath.MustNewDecFromString("0.1")
	pNorm := alloraMath.MustNewDecFromString("3.0")
	cNorm := alloraMath.MustNewDecFromString("0.75")
	epsilon := alloraMath.MustNewDecFromString("1e-4")
	initialRegretQuantile := alloraMath.MustNewDecFromString("0.5")
	pnormSafeDiv := alloraMath.MustNewDecFromString("1.0")

	// Set initial Regret to check if this value is updated or not
	initialRegret := alloraMath.MustNewDecFromString("0")
	// Create new topic
	topic := s.mockTopic()
	topic.InitialRegret = initialRegret
	err := s.emissionsKeeper.SetTopic(s.ctx, topicId, topic)
	s.Require().NoError(err)

	inferer0 := s.addrs[0].String()
	inferer1 := s.addrs[1].String()
	inferer2 := s.addrs[2].String()
	inferer3 := s.addrs[3].String()
	inferer4 := s.addrs[4].String()
	infererAddresses := []string{inferer0, inferer1, inferer2, inferer3, inferer4}

	forecaster0 := s.addrs[5].String()
	forecaster1 := s.addrs[6].String()
	forecaster2 := s.addrs[7].String()
	forecasterAddresses := []string{forecaster0, forecaster1, forecaster2}

	reputer0 := s.addrs[8].String()

	err = testutil.SetRegretsFromPreviousEpoch(s.ctx, s.emissionsKeeper, topicId, blockHeight, infererAddresses, forecasterAddresses, epochPrevGet)
	require.NoError(err)

	networkLosses, err := testutil.GetNetworkLossFromCsv(
		topicId,
		blockHeight,
		infererAddresses,
		forecasterAddresses,
		reputer0,
		epoch301Get,
	)
	s.Require().NoError(err)

	err = inferencesynthesis.GetCalcSetNetworkRegrets(inferencesynthesis.GetCalcSetNetworkRegretsArgs{
		Ctx:                   s.ctx,
		K:                     k,
		TopicId:               topicId,
		NetworkLosses:         networkLosses,
		Nonce:                 nonce,
		AlphaRegret:           alpha,
		CNorm:                 cNorm,
		PNorm:                 pNorm,
		EpsilonTopic:          epsilon,
		InitialRegretQuantile: initialRegretQuantile,
		PNormSafeDiv:          pnormSafeDiv,
	})
	require.NoError(err)

	// Initial Regret will not be updated because this topic has no experienced actors
	topic, err = s.emissionsKeeper.GetTopic(s.ctx, topicId)
	require.NoError(err)
	require.Equal(topic.InitialRegret, initialRegret)
}

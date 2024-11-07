package fuzz_test

import (
	"context"
	"encoding/hex"
	"math/rand"

	cosmossdk_io_math "cosmossdk.io/math"
	alloraMath "github.com/allora-network/allora-chain/math"
	testcommon "github.com/allora-network/allora-chain/test/common"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/stretchr/testify/require"
)

func doInferenceAndReputation(
	m *testcommon.TestConfig,
	_ Actor,
	_ Actor,
	_ *cosmossdk_io_math.Int,
	topicId uint64,
	data *SimulationData,
	iteration int,
) (success bool) {
	iterLog(
		m.T,
		iteration,
		"producing inference and reputation for topic id",
		topicId,
	)
	ctx := context.Background()
	resp, err := m.Client.QueryEmissions().GetTopic(ctx, &emissionstypes.GetTopicRequest{
		TopicId: topicId,
	})
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"failed to upload inferences, forecast and reputations for topic id ",
			topicId,
			"query GetTopic error",
			err,
		)
		return false
	}
	topic := resp.Topic
	iterLog(m.T, iteration, "Inference topic epoch last ended ", topic.EpochLastEnded, " epoch length ", topic.EpochLength)
	workerNonce := topic.EpochLastEnded + topic.EpochLength
	blockHeightNow, err := m.Client.BlockHeight(ctx)
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"failed to upload inferences, forecast and reputations for topic id ",
			topicId,
			"block height query error",
			err,
		)
		return false
	}
	if blockHeightNow < workerNonce+1 {
		iterLog(m.T, iteration, "waiting for next epoch to start so we can produce inferences for the current epoch: ", workerNonce+1)
		err = m.Client.WaitForBlockHeight(ctx, workerNonce+1)
		failIfOnErr(m.T, data.failOnErr, err)
		if err != nil {
			iterFailLog(
				m.T,
				iteration,
				"failed to upload inferences, forecast and reputations for topic id ",
				topicId,
				"wait for block height error",
				err,
			)
			return false
		}
		// Update block height
		blockHeightNow, err = m.Client.BlockHeight(ctx)
		failIfOnErr(m.T, data.failOnErr, err)
		if err != nil {
			iterFailLog(
				m.T,
				iteration,
				"failed to upload inferences, forecast and reputations for topic id ",
				topicId,
				"block height query2 error",
				err,
			)
			return false
		}
	}
	workers := data.getWorkersForTopic(topicId)
	if len(workers) == 0 {
		iterFailLog(m.T, iteration, "len of workers in active topic should always be greater than 0 ", topicId)
	}
	iterLog(m.T, iteration, " starting worker payload topic id ", topicId,
		" workers ", workers, "worker nonce ",
		workerNonce, " block height now ", blockHeightNow)
	workerPayloadSuccess := createAndSendWorkerPayloads(m, data, topic, workers, workerNonce, iteration)
	if !workerPayloadSuccess {
		iterFailLog(m.T, iteration, "worker payload errored topic", topicId)
		return false
	}
	iterLog(m.T, iteration, "produced worker inference for topic id", topicId)
	reputerWaitBlocks := blockHeightNow + topic.GroundTruthLag + 1
	iterLog(m.T, iteration, "waiting for reputer ground truth block", reputerWaitBlocks)
	err = m.Client.WaitForBlockHeight(ctx, reputerWaitBlocks)
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"failed to upload inferences, forecast and reputations for topic id ",
			topicId,
			"wait for block height 3 error",
			err,
		)
		return false
	}
	reputers := data.getReputersForTopicWithStake(topicId)
	if len(reputers) == 0 {
		iterFailLog(m.T, iteration, "len of reputers in active topic should always be greater than 0 ", topicId)
	}
	iterLog(
		m.T, iteration, " starting reputer payload topic id ", topicId,
		" workers ", workers, " reputers ", reputers, " worker nonce ", workerNonce,
		" block height  now ", reputerWaitBlocks,
	)
	reputationSuccess := createAndSendReputerPayloads(m, data, topic, reputers, workers, workerNonce, iteration)
	if !reputationSuccess {
		iterFailLog(m.T, iteration, "reputation flow failed topic id ", topicId)
		return false
	}

	data.counts.incrementDoInferenceAndReputationCount()
	iterSuccessLog(m.T, iteration, "uploaded inferences, forecasts, and reputations for topic id ", topicId)
	return true
}

// determine if this state transition is worth trying based on our knowledge of the state
func findActiveTopicsAtThisBlock(
	m *testcommon.TestConfig,
	data *SimulationData,
	blockHeight int64,
) []*emissionstypes.Topic {
	// first off someone has to be registered for both working and reputing in general
	if !anyReputersRegistered(data) || !anyWorkersRegistered(data) {
		return nil
	}
	ctx := context.Background()
	response, err := m.Client.QueryEmissions().GetActiveTopicsAtBlock(ctx, &emissionstypes.GetActiveTopicsAtBlockRequest{
		BlockHeight: blockHeight,
	})
	failIfOnErr(m.T, data.failOnErr, err)
	return response.Topics
}

// Inserts inference and forecast data for each worker
func createAndSendWorkerPayloads(
	m *testcommon.TestConfig,
	data *SimulationData,
	topic *emissionstypes.Topic,
	workers []Actor,
	workerNonce int64,
	iteration int,
) (success bool) {
	// Get Bundles
	for _, worker := range workers {
		workerData := createWorkerDataBundle(m, topic.Id, workerNonce, worker, workers)
		success = sendWorkerPayload(m, data, worker, workerData, iteration)
		if !success {
			return false
		}
	}
	return true
}

// create inferences and forecasts for a worker
func createWorkerDataBundle(
	m *testcommon.TestConfig,
	topicId uint64,
	blockHeight int64,
	inferer Actor,
	workers []Actor,
) *emissionstypes.WorkerDataBundle {
	// Iterate workerAddresses to get the worker address, and generate as many forecasts as there are workers
	forecastElements := make([]*emissionstypes.ForecastElement, 0)
	for key := range workers {
		forecastElements = append(forecastElements, &emissionstypes.ForecastElement{
			Inferer: workers[key].addr,
			Value:   alloraMath.NewDecFromInt64(int64(m.Client.Rand.Intn(51) + 50)),
		})
	}
	infererAddress := inferer.addr
	infererValue := alloraMath.NewDecFromInt64(int64(m.Client.Rand.Intn(300) + 3000))

	workerDataBundle := &emissionstypes.WorkerDataBundle{
		Worker: infererAddress,
		Nonce: &emissionstypes.Nonce{
			BlockHeight: blockHeight,
		},
		TopicId: topicId,
		InferenceForecastsBundle: &emissionstypes.InferenceForecastBundle{
			Inference: &emissionstypes.Inference{
				TopicId:     topicId,
				BlockHeight: blockHeight,
				Inferer:     infererAddress,
				Value:       infererValue,
				ExtraData:   nil,
				Proof:       "",
			},
			Forecast: &emissionstypes.Forecast{
				TopicId:          topicId,
				BlockHeight:      blockHeight,
				Forecaster:       infererAddress,
				ForecastElements: forecastElements,
				ExtraData:        nil,
			},
		},
		InferencesForecastsBundleSignature: nil,
		Pubkey:                             "",
	}

	// Sign
	src := make([]byte, 0)
	src, err := workerDataBundle.InferenceForecastsBundle.XXX_Marshal(src, true)
	require.NoError(m.T, err, "Marshall reputer value bundle should not return an error")

	sig, pubKey, err := m.Client.Context().Keyring.Sign(inferer.name, src, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(m.T, err, "Sign should not return an error")
	workerPublicKeyBytes := pubKey.Bytes()
	workerDataBundle.InferencesForecastsBundleSignature = sig
	workerDataBundle.Pubkey = hex.EncodeToString(workerPublicKeyBytes)

	return workerDataBundle
}

// Send worker payload, from worker address (which should exist in the keyring)
func sendWorkerPayload(
	m *testcommon.TestConfig,
	data *SimulationData,
	sender Actor,
	WorkerDataBundles *emissionstypes.WorkerDataBundle,
	iteration int,
) bool {
	workerMsg := &emissionstypes.InsertWorkerPayloadRequest{
		Sender:           sender.addr,
		WorkerDataBundle: WorkerDataBundles,
	}
	// serialize workerMsg to json and print
	LeaderAcc, err := m.Client.AccountRegistryGetByName(sender.name)
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"send worker payload error",
			"could not get account by name",
			err,
		)
		return false
	}
	ctx := context.Background()
	txResp, err := m.Client.BroadcastTx(ctx, LeaderAcc, workerMsg)
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"send worker payload error",
			"broadcast tx error",
			err,
		)
		return false
	}
	_, err = m.Client.WaitForTx(ctx, txResp.TxHash)
	failIfOnErr(m.T, data.failOnErr, err)
	if err != nil {
		iterFailLog(
			m.T,
			iteration,
			"send worker payload error",
			"wait for tx error",
			err,
		)
		return false
	}
	return true
}

// reputers submit their assessment of the quality of workers' work compared to ground truth
func createAndSendReputerPayloads(
	m *testcommon.TestConfig,
	data *SimulationData,
	topic *emissionstypes.Topic,
	reputers,
	workers []Actor,
	workerNonce int64,
	iteration int,
) (success bool) {
	// Nonce: calculate from EpochLastRan + EpochLength
	topicId := topic.Id
	// Nonces are last two blockHeights
	reputerNonce := &emissionstypes.Nonce{
		BlockHeight: workerNonce,
	}
	ctx := context.Background()
	for _, reputer := range reputers {
		valueBundle := createReputerValueBundle(m, topicId, reputer, workers, reputerNonce)
		signedValueBundle := signReputerValueBundle(m, reputer, valueBundle)
		lossesMsg := &emissionstypes.InsertReputerPayloadRequest{
			Sender:             reputer.addr,
			ReputerValueBundle: signedValueBundle,
		}

		txResp, err := m.Client.BroadcastTx(ctx, reputer.acc, lossesMsg)
		failIfOnErr(m.T, data.failOnErr, err)
		if err != nil {
			iterFailLog(
				m.T,
				iteration,
				"send reputer payload error",
				"broadcast tx error",
				err,
			)
			return false
		}
		_, err = m.Client.WaitForTx(ctx, txResp.TxHash)
		failIfOnErr(m.T, data.failOnErr, err)
		if err != nil {
			iterFailLog(
				m.T,
				iteration,
				"send reputer payload error",
				"wait for tx error",
				err,
			)
			return false
		}
	}
	return true
}

// Generate the same valueBundle for a reputer
func createReputerValueBundle(
	m *testcommon.TestConfig,
	topicId uint64,
	reputer Actor,
	workers []Actor,
	reputerNonce *emissionstypes.Nonce,
) emissionstypes.ValueBundle {
	return emissionstypes.ValueBundle{
		TopicId:                topicId,
		Reputer:                reputer.addr,
		ExtraData:              nil,
		CombinedValue:          alloraMath.NewDecFromInt64(100),
		InfererValues:          generateWorkerAttributedValueLosses(m, workers, 3000, 3500),
		ForecasterValues:       generateWorkerAttributedValueLosses(m, workers, 50, 50),
		NaiveValue:             alloraMath.NewDecFromInt64(100),
		OneOutInfererValues:    generateWithheldWorkerAttributedValueLosses(workers, 50, 50),
		OneOutForecasterValues: generateWithheldWorkerAttributedValueLosses(workers, 50, 50),
		OneInForecasterValues:  generateWorkerAttributedValueLosses(m, workers, 50, 50),
		ReputerRequestNonce: &emissionstypes.ReputerRequestNonce{
			ReputerNonce: reputerNonce,
		},
		OneOutInfererForecasterValues: nil,
	}
}

// Generate a ReputerValueBundle:of
func signReputerValueBundle(
	m *testcommon.TestConfig,
	reputer Actor,
	valueBundle emissionstypes.ValueBundle,
) *emissionstypes.ReputerValueBundle {
	valueBundle.Reputer = reputer.addr
	// Sign
	src := make([]byte, 0)
	src, err := valueBundle.XXX_Marshal(src, true)
	require.NoError(m.T, err, "Marshall reputer value bundle should not return an error")

	valueBundleSignature, pubKey, err := m.Client.Context().Keyring.Sign(reputer.name, src, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(m.T, err, "Sign should not return an error")
	reputerPublicKeyBytes := pubKey.Bytes()

	// Create a InsertReputerPayloadRequest message
	reputerValueBundle := &emissionstypes.ReputerValueBundle{
		ValueBundle: &valueBundle,
		Signature:   valueBundleSignature,
		Pubkey:      hex.EncodeToString(reputerPublicKeyBytes),
	}

	return reputerValueBundle
}

// for every worker, generate a worker attributed value
func generateWorkerAttributedValueLosses(
	m *testcommon.TestConfig,
	workers []Actor,
	lowLimit,
	sum int,
) []*emissionstypes.WorkerAttributedValue {
	values := make([]*emissionstypes.WorkerAttributedValue, 0)
	for _, worker := range workers {
		values = append(values, &emissionstypes.WorkerAttributedValue{
			Worker: worker.addr,
			Value:  alloraMath.NewDecFromInt64(int64(m.Client.Rand.Intn(lowLimit) + sum)),
		})
	}
	return values
}

// for every worker, generate a withheld worker attribute value
func generateWithheldWorkerAttributedValueLosses(
	workers []Actor,
	lowLimit,
	sum int,
) []*emissionstypes.WithheldWorkerAttributedValue {
	values := make([]*emissionstypes.WithheldWorkerAttributedValue, 0)
	for _, worker := range workers {
		values = append(values, &emissionstypes.WithheldWorkerAttributedValue{
			Worker: worker.addr,
			Value:  alloraMath.NewDecFromInt64(int64(rand.Intn(lowLimit) + sum)),
		})
	}
	return values
}

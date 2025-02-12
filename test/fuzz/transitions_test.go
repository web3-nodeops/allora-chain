package fuzz_test

import (
	"context"
	"fmt"

	cosmossdk_io_math "cosmossdk.io/math"
	testcommon "github.com/allora-network/allora-chain/test/common"
	fuzzcommon "github.com/allora-network/allora-chain/test/fuzz/common"
)

var UnusedActor Actor = Actor{} // nolint:exhaustruct

// Every function responsible for doing a state transition
// should adhere to this function signature
type StateTransitionFunc func(
	m *testcommon.TestConfig,
	actor1 Actor,
	actor2 Actor,
	amount *cosmossdk_io_math.Int,
	topicId uint64,
	data *SimulationData,
	iteration int,
) bool

// keep track of the name of the state transition as well as the function
type StateTransition struct {
	name         string              // The name of this state transition
	f            StateTransitionFunc // Which function to call
	weight       uint8               // Interpreted as a percentage (0-100), all weights should sum to 100
	follow       *StateTransition    // If there is a follow on state transition that can happen after this, else nil
	followWeight uint8               // Interpreted as a percentage (0-100)
}

// The list of possible state transitions we can take are:
//
// create a new topic,
// fund a topic some more,
// register as a reputer,
// register as a worker,
// unregister as a reputer,
// unregister as a worker,
// stake as a reputer,
// stake in a reputer (delegate),
// unstake as a reputer,
// unstake from a reputer (undelegate),
// cancel the removal of stake (as a reputer),
// cancel the removal of delegated stake (delegator),
// collect delegator rewards,
// produce an inference (insert worker payloads),
// produce reputation scores (insert reputer payloads)
// NOTE: all weights must sum to 100
func allTransitions(f *fuzzcommon.FuzzConfig) []StateTransition {
	transitionWeights := f.TransitionWeights
	transitionCreateTopic := StateTransition{
		name: "createTopic", f: createTopic,
		weight: transitionWeights.CreateTopic,
		follow: nil, followWeight: 0,
	}
	transitionFundTopic := StateTransition{
		name: "fundTopic", f: fundTopic,
		weight: transitionWeights.FundTopic,
		follow: nil, followWeight: 0,
	}
	transitionRegisterWorker := StateTransition{
		name: "registerWorker", f: registerWorker,
		weight: transitionWeights.RegisterWorker,
		follow: nil, followWeight: 0,
	}
	transitionRegisterReputer := StateTransition{
		name: "registerReputer", f: registerReputer,
		weight: transitionWeights.RegisterReputer,
		follow: nil, followWeight: 0,
	}
	transitionStakeAsReputer := StateTransition{
		name: "stakeAsReputer", f: stakeAsReputer,
		weight: transitionWeights.StakeAsReputer,
		follow: nil, followWeight: 0,
	}
	transitionDelegateStake := StateTransition{
		name: "delegateStake", f: delegateStake,
		weight: transitionWeights.DelegateStake,
		follow: nil, followWeight: 0,
	}
	transitionCollectDelegatorRewards := StateTransition{
		name: "collectDelegatorRewards", f: collectDelegatorRewards,
		weight: transitionWeights.CollectDelegatorRewards,
		follow: nil, followWeight: 0,
	}
	transitionDoInferenceAndReputation := StateTransition{
		name: "doInferenceAndReputation", f: doInferenceAndReputation,
		weight: transitionWeights.DoInferenceAndReputation,
		follow: nil, followWeight: 0,
	}
	transitionUnregisterWorker := StateTransition{
		name: "unregisterWorker", f: unregisterWorker,
		weight: transitionWeights.UnregisterWorker,
		follow: nil, followWeight: 0,
	}
	transitionUnregisterReputer := StateTransition{
		name: "unregisterReputer", f: unregisterReputer,
		weight: transitionWeights.UnregisterReputer,
		follow: nil, followWeight: 0,
	}

	//cancels come after the unstake/undelegate
	transitionCancelStakeRemoval := StateTransition{
		name: "cancelStakeRemoval", f: cancelStakeRemoval,
		weight: transitionWeights.CancelStakeRemoval,
		follow: nil, followWeight: 0,
	}
	transitionCancelDelegateStakeRemoval := StateTransition{
		name: "cancelDelegateStakeRemoval", f: cancelDelegateStakeRemoval,
		weight: transitionWeights.CancelDelegateStakeRemoval,
		follow: nil, followWeight: 0,
	}
	transitionUnstakeAsReputer := StateTransition{
		name: "unstakeAsReputer", f: unstakeAsReputer,
		weight: transitionWeights.UnstakeAsReputer,
		follow: &transitionCancelStakeRemoval, followWeight: 50,
	}
	transitionUndelegateStake := StateTransition{
		name: "undelegateStake", f: undelegateStake,
		weight: transitionWeights.UndelegateStake,
		follow: &transitionCancelDelegateStakeRemoval, followWeight: 50,
	}

	return []StateTransition{
		transitionCreateTopic,
		transitionFundTopic,
		transitionRegisterWorker,
		transitionRegisterReputer,
		transitionStakeAsReputer,
		transitionDelegateStake,
		transitionCollectDelegatorRewards,
		transitionDoInferenceAndReputation,
		transitionUnregisterWorker,
		transitionUnregisterReputer,
		transitionUnstakeAsReputer,
		transitionUndelegateStake,
		transitionCancelStakeRemoval,
		transitionCancelDelegateStakeRemoval,
	}
}

// weight transitions that add registrations or stake, more heavily than those that take it away
// 70% of the time do additive stuff
// 30% of the time do subtractive stuff
func pickTransitionWithWeight(f *fuzzcommon.FuzzConfig) StateTransition {
	transitions := allTransitions(f)
	rand := f.TestConfig.Client.Rand.Intn(100)
	threshold := uint8(0)
	prevThreshold := uint8(0)
	for _, transition := range transitions {
		threshold += transition.weight
		if rand >= int(prevThreshold) && rand < int(threshold) {
			return transition
		}
		prevThreshold = threshold
	}
	panic(fmt.Sprintf("Weights must sum to 100 and rand should pick a value between 0 and 100: %d %d", rand, threshold))
}

// if it is possible to pick a follow on transition, pick it with this weight
// if we decide not to pick it, return nil
func pickFollowOnTransitionWithWeight(m *testcommon.TestConfig, currTransition StateTransition) *StateTransition {
	rand := m.Client.Rand.Intn(100)
	if rand < int(currTransition.followWeight) {
		return currTransition.follow
	}
	return nil
}

// State machine dependencies for valid transitions
//
// fundTopic: CreateTopic
// RegisterWorkerForTopic: CreateTopic
// RegisterReputerForTopic: CreateTopic
// unRegisterReputer: RegisterReputerForTopic
// unRegisterWorker: RegisterWorkerForTopic
// stakeReputer: RegisterReputerForTopic, CreateTopic
// delegateStake: CreateTopic, RegisterReputerForTopic
// unstakeReputer: stakeReputer
// unstakeDelegator: delegateStake
// cancelStakeRemoval: unstakeReputer
// cancelDelegateStakeRemoval: unstakeDelegator
// collectDelegatorRewards: delegateStake, fundTopic, InsertWorkerPayload, InsertReputerPayload
// InsertWorkerPayload: RegisterWorkerForTopic, FundTopic
// InsertReputerPayload: RegisterReputerForTopic, InsertWorkerPayload
//
// Note: this function is not exhaustive, it only checks for the most obvious/low hanging fruit conditions
// pickActorAndTopicIdForStateTransition will do further checks
func canTransitionOccur(m *testcommon.TestConfig, data *SimulationData, transition StateTransition) bool {
	switch transition.name {
	case "unregisterWorker":
		return anyWorkersRegistered(data)
	case "unregisterReputer":
		return anyReputersRegistered(data)
	case "stakeAsReputer":
		return anyReputersRegistered(data)
	case "delegateStake":
		return anyReputersRegistered(data)
	case "unstakeAsReputer":
		return anyReputersStaked(data)
	case "undelegateStake":
		return anyDelegatorsStaked(data)
	case "collectDelegatorRewards":
		return anyDelegatorsStaked(data) && anyReputersRegistered(data)
	case "cancelStakeRemoval":
		return anyReputersStaked(data)
	case "cancelDelegateStakeRemoval":
		return anyDelegatorsStaked(data)
	case "doInferenceAndReputation":
		ctx := context.Background()
		blockHeightNow, err := m.Client.BlockHeight(ctx)
		if err != nil {
			return false
		}
		activeTopics := findActiveTopicsAtThisBlock(m, data, blockHeightNow)
		for i := 0; i < len(activeTopics); i++ {
			workerExists := data.isAnyWorkerRegisteredInTopic(activeTopics[i].Id)
			reputerExists := data.isAnyReputerRegisteredInTopic(activeTopics[i].Id)
			if workerExists && reputerExists {
				return true
			}
		}
		return false

	default:
		return true
	}
}

// is this specific combination of actors, amount, and topicId valid for the transition?
func isValidTransition(m *testcommon.TestConfig, transition StateTransition, actor1 Actor, actor2 Actor, amount *cosmossdk_io_math.Int, topicId uint64, data *SimulationData, iteration int) bool {
	switch transition.name {
	case "registerWorker":
		// can't register twice
		if data.isWorkerRegisteredInTopic(topicId, actor1) {
			return false
		}
		return true
	case "registerReputer":
		// can't register twice
		if data.isReputerRegisteredInTopic(topicId, actor1) {
			return false
		}
		return true
	case "collectDelegatorRewards":
		// if the reputer unregisters before the delegator withdraws stake, it can be invalid for a
		// validator to collective rewards
		if !data.isReputerRegisteredInTopic(topicId, actor2) {
			iterLog(m.T, iteration, "Transition not valid: ", transition.name, actor1, actor2, amount, topicId)
			return false
		}
		return true
	default:
		return true
	}
}

// pickRandomActor picks a random actor from the list of actors in the simulation data
func pickRandomActor(m *testcommon.TestConfig, data *SimulationData) Actor {
	return data.actors[m.Client.Rand.Intn(len(data.actors))]
}

// pickRandomActorExcept picks a random actor from the list of actors in the simulation data
// and panics if it can't find one after 5 tries that is not the same as the given actors
func pickRandomActorExcept(m *testcommon.TestConfig, data *SimulationData, actors []Actor) Actor {
	count := 0
	for ; count < 5; count++ {
		randomActor := pickRandomActor(m, data)
		match := false
		for _, actor := range actors {
			if randomActor == actor {
				match = true
			}
		}
		if !match {
			return randomActor
		}
	}
	panic(
		fmt.Sprintf(
			"could not find a random actor that is not the same as the given actor after %d tries",
			count,
		),
	)
}

// helper for when the transition values can be fully fully random
func pickFullRandomValues(
	m *testcommon.TestConfig,
	data *SimulationData,
) (bool, Actor, Actor, *cosmossdk_io_math.Int, uint64) {
	randomTopicId, err := pickRandomTopicId(m)
	failIfOnErr(m.T, data.failOnErr, err)
	randomActor1 := pickRandomActor(m, data)
	randomActor2 := pickRandomActor(m, data)
	amount, err := pickRandomBalanceLessThanHalf(m, randomActor1)
	failIfOnErr(m.T, data.failOnErr, err)
	return true, randomActor1, randomActor2, &amount, randomTopicId
}

// pickActorAndTopicIdForStateTransition picks random actors
// able to take the state transition and returns which one it picked.
// if the transition requires only one actor (the majority) the second is empty
func pickActorAndTopicIdForStateTransition(
	m *testcommon.TestConfig,
	transition StateTransition,
	data *SimulationData,
	iteration int,
) (success bool, actor1 Actor, actor2 Actor, amount *cosmossdk_io_math.Int, topicId uint64) {
	switch transition.name {
	case "unregisterWorker":
		worker, topicId, err := data.pickRandomRegisteredWorker()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		return true, worker, UnusedActor, nil, topicId
	case "unregisterReputer":
		reputer, topicId, err := data.pickRandomRegisteredReputer()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		return true, reputer, UnusedActor, nil, topicId
	case "stakeAsReputer":
		reputer, topicId, err := data.pickRandomRegisteredReputer()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		amount, err := pickRandomBalanceLessThanHalf(m, reputer) // if err amount=zero which is a valid transition
		failIfOnErr(m.T, data.failOnErr, err)
		return true, reputer, UnusedActor, &amount, topicId
	case "delegateStake":
		reputer, topicId, err := data.pickRandomRegisteredReputer()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		delegator := pickRandomActorExcept(m, data, []Actor{reputer})
		amount, err := pickRandomBalanceLessThanHalf(m, delegator)
		failIfOnErr(m.T, data.failOnErr, err)
		return true, delegator, reputer, &amount, topicId
	case "unstakeAsReputer":
		reputer, topicId, err := data.pickRandomStakedReputer()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		amount := pickPercentOfStakeByReputer(m, topicId, reputer, data, iteration)
		return true, reputer, UnusedActor, &amount, topicId
	case "undelegateStake":
		delegator, reputer, topicId, err := data.pickRandomStakedDelegator()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		amount := pickPercentOfStakeByDelegator(m, topicId, delegator, reputer, data, iteration)
		return true, delegator, reputer, &amount, topicId
	case "collectDelegatorRewards":
		delegator, reputer, topicId, err := data.pickRandomStakedDelegator()
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		return true, delegator, reputer, nil, topicId
	case "doInferenceAndReputation":
		ctx := context.Background()
		blockHeightNow, err := m.Client.BlockHeight(ctx)
		if err != nil {
			return false, UnusedActor, UnusedActor, nil, 0
		}
		topics := findActiveTopicsAtThisBlock(m, data, blockHeightNow)
		if len(topics) > 0 {
			randIndex := m.Client.Rand.Intn(len(topics))
			topicId := topics[randIndex].Id
			return true, UnusedActor, UnusedActor, nil, topicId
		}
		return false, UnusedActor, UnusedActor, nil, 0
	default:
		return pickFullRandomValues(m, data)
	}
}

package fuzz_test

import (
	"cmp"
	"fmt"
	"math/rand"
	"slices"

	cosmossdk_io_math "cosmossdk.io/math"
	testcommon "github.com/allora-network/allora-chain/test/common"
)

// SimulationData stores the active set of states we think we're in
// so that we can choose to take a transition that is valid
// right now it doesn't need mutexes, if we parallelize this test ever it will
// to read and write out of the simulation data
type SimulationData struct {
	epochLength        int64
	actors             []Actor
	counts             StateTransitionCounts
	registeredWorkers  *testcommon.RandomKeyMap[Registration, struct{}]
	registeredReputers *testcommon.RandomKeyMap[Registration, struct{}]
	reputerStakes      *testcommon.RandomKeyMap[Registration, struct{}]
	delegatorStakes    *testcommon.RandomKeyMap[Delegation, struct{}]
	failOnErr          bool
	mode               SimulationMode
}

// String is the stringer for SimulationData
func (s *SimulationData) String() string {
	return fmt.Sprintf(
		"SimulationData{\nepochLength: %d,\nactors: %v,\n counts: %s,\nregisteredWorkers: %v,\nregisteredReputers: %v,\nreputerStakes: %v,\ndelegatorStakes: %v,\nfailOnErr: %v,\nmode: %s}",
		s.epochLength,
		s.actors,
		s.counts,
		s.registeredWorkers,
		s.registeredReputers,
		s.reputerStakes,
		s.delegatorStakes,
		s.failOnErr,
		s.mode,
	)
}

type Registration struct {
	TopicId uint64
	Actor   Actor
}

type Delegation struct {
	TopicId   uint64
	Delegator Actor
	Reputer   Actor
}

// addWorkerRegistration adds a worker registration to the simulation data
func (s *SimulationData) addWorkerRegistration(topicId uint64, actor Actor) {
	s.registeredWorkers.Upsert(Registration{
		TopicId: topicId,
		Actor:   actor,
	}, struct{}{})
}

// removeWorkerRegistration removes a worker registration from the simulation data
func (s *SimulationData) removeWorkerRegistration(topicId uint64, actor Actor) {
	s.registeredWorkers.Delete(Registration{
		TopicId: topicId,
		Actor:   actor,
	})
}

// addReputerRegistration adds a reputer registration to the simulation data
func (s *SimulationData) addReputerRegistration(topicId uint64, actor Actor) {
	s.registeredReputers.Upsert(Registration{
		TopicId: topicId,
		Actor:   actor,
	}, struct{}{})
}

// addReputerStaked adds a reputer stake to the list of staked reputers in the simulation data
func (s *SimulationData) addReputerStaked(topicId uint64, actor Actor) {
	s.reputerStakes.Upsert(Registration{
		TopicId: topicId,
		Actor:   actor,
	}, struct{}{})
}

// addDelegatorDelegated adds a delegator stake to the list of staked delegators in the simulation data
func (s *SimulationData) addDelegatorDelegated(topicId uint64, delegator Actor, reputer Actor) {
	s.delegatorStakes.Upsert(Delegation{
		TopicId:   topicId,
		Delegator: delegator,
		Reputer:   reputer,
	}, struct{}{})
}

// removeReputerRegistration removes a reputer registration from the simulation data
func (s *SimulationData) removeReputerRegistration(topicId uint64, actor Actor) {
	s.registeredReputers.Delete(Registration{
		TopicId: topicId,
		Actor:   actor,
	})
}

// removeReputerStaked removes a reputer stake from the list of staked reputers in the simulation data
func (s *SimulationData) removeReputerStaked(topicId uint64, actor Actor) {
	s.reputerStakes.Delete(Registration{
		TopicId: topicId,
		Actor:   actor,
	})
}

// removeDelegatorDelegated removes a delegator stake from the list of staked delegators in the simulation data
func (s *SimulationData) removeDelegatorDelegated(topicId uint64, delegator Actor, reputer Actor) {
	s.delegatorStakes.Delete(Delegation{
		TopicId:   topicId,
		Delegator: delegator,
		Reputer:   reputer,
	})
}

// pickRandomRegisteredWorker picks a random worker that is currently registered
func (s *SimulationData) pickRandomRegisteredWorker() (Actor, uint64, error) {
	ret, err := s.registeredWorkers.RandomKey()
	if err != nil {
		return Actor{}, 0, err
	}
	return ret.Actor, ret.TopicId, nil
}

// pickRandomRegisteredReputer picks a random reputer that is currently registered
func (s *SimulationData) pickRandomRegisteredReputer() (Actor, uint64, error) {
	ret, err := s.registeredReputers.RandomKey()
	if err != nil {
		return Actor{}, 0, err
	}
	return ret.Actor, ret.TopicId, nil
}

// pickRandomStakedReputer picks a random reputer that is currently staked
func (s *SimulationData) pickRandomStakedReputer() (Actor, uint64, error) {
	actor, topicId, err := s.pickRandomRegisteredReputer()
	if err != nil {
		return Actor{}, 0, err
	}
	reg := Registration{
		TopicId: topicId,
		Actor:   actor,
	}
	_, exists := s.reputerStakes.Get(reg)
	if !exists {
		return Actor{}, 0, fmt.Errorf("Registered reputer %s is not staked", actor.addr)
	}
	return actor, topicId, nil
}

// pickRandomDelegator picks a random delegator that is currently staked
func (s *SimulationData) pickRandomStakedDelegator() (Actor, Actor, uint64, error) {
	ret, err := s.delegatorStakes.RandomKey()
	if err != nil {
		return Actor{}, Actor{}, 0, err
	}

	if !s.isReputerRegistered(ret.TopicId, ret.Reputer) {
		return Actor{}, Actor{}, 0, fmt.Errorf(
			"Delegator %s is staked in reputer %s, but reputer is not registered",
			ret.Delegator.addr,
			ret.Reputer.addr,
		)
	}

	return ret.Delegator, ret.Reputer, ret.TopicId, nil
}

// take a percentage of the stake, either 1/10, 1/3, 1/2, 6/7, or the full amount
func pickPercentOf(rand *rand.Rand, stake cosmossdk_io_math.Int) cosmossdk_io_math.Int {
	if stake.Equal(cosmossdk_io_math.ZeroInt()) {
		return cosmossdk_io_math.ZeroInt()
	}
	percent := rand.Intn(5)
	switch percent {
	case 0:
		return stake.QuoRaw(10)
	case 1:
		return stake.QuoRaw(3)
	case 2:
		return stake.QuoRaw(2)
	case 3:
		return stake.MulRaw(6).QuoRaw(7)
	default:
		return stake
	}
}

// isReputerRegistered checks if a reputer is registered
func (s *SimulationData) isReputerRegistered(topicId uint64, actor Actor) bool {
	_, exists := s.registeredReputers.Get(Registration{
		TopicId: topicId,
		Actor:   actor,
	})
	return exists
}

// isAnyWorkerRegisteredInTopic checks if any worker is registered in a topic
func (s *SimulationData) isAnyWorkerRegisteredInTopic(topicId uint64) bool {
	workers, _ := s.registeredWorkers.Filter(func(reg Registration) bool {
		return reg.TopicId == topicId
	})
	return len(workers) > 0
}

// isAnyReputerRegisteredInTopic checks if any reputer is registered in a topic
func (s *SimulationData) isAnyReputerRegisteredInTopic(topicId uint64) bool {
	reputers, _ := s.registeredReputers.Filter(func(reg Registration) bool {
		return reg.TopicId == topicId
	})
	return len(reputers) > 0
}

// get all workers for a topic, this function is iterates over the list of workers multiple times
// for determinism, the workers are sorted by their address
func (s *SimulationData) getWorkersForTopic(topicId uint64) []Actor {
	workers, _ := s.registeredWorkers.Filter(func(reg Registration) bool {
		return reg.TopicId == topicId
	})
	ret := make([]Actor, len(workers))
	for i, worker := range workers {
		ret[i] = worker.Actor
	}
	slices.SortFunc(ret, func(a, b Actor) int {
		return cmp.Compare(a.addr, b.addr)
	})
	return ret
}

// get all reputers with nonzero stake for a topic, this function is iterates over the list of reputers multiple times
// for determinism, the reputers are sorted by their address
func (s *SimulationData) getReputersForTopicWithStake(topicId uint64) []Actor {
	reputerRegs, _ := s.reputerStakes.Filter(func(reg Registration) bool {
		return reg.TopicId == topicId
	})
	rmap := make(map[string]Actor)
	for _, reputerReg := range reputerRegs {
		rmap[reputerReg.Actor.addr] = reputerReg.Actor
	}
	reputerDels, _ := s.delegatorStakes.Filter(func(del Delegation) bool {
		return del.TopicId == topicId
	})
	for _, del := range reputerDels {
		rmap[del.Reputer.addr] = del.Reputer
	}
	ret := make([]Actor, 0)
	for _, reputer := range rmap {
		ret = append(ret, reputer)
	}
	slices.SortFunc(ret, func(a, b Actor) int {
		return cmp.Compare(a.addr, b.addr)
	})
	return ret
}

// randomly flip the fail on err case to decide whether to be aggressive and fuzzy or
// behaved state transitions
func (s *SimulationData) randomlyFlipFailOnErr(m *testcommon.TestConfig, iteration int) {
	// 20% likely to change from what you were previously
	if m.Client.Rand.Intn(10) >= 8 {
		iterLog(m.T, iteration, "Changing fuzzer mode: failOnErr changing from", s.failOnErr, "to", !s.failOnErr)
		s.failOnErr = !s.failOnErr
	}
}

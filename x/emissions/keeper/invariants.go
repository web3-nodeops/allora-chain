package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	cosmosMath "cosmossdk.io/math"
	"github.com/allora-network/allora-chain/app/params"
	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterInvariants registers the emissions module invariants.
func RegisterInvariants(ir sdk.InvariantRegistry, k *Keeper) {
	ir.RegisterRoute(emissionstypes.ModuleName, "allora-staking-total-balance", StakingInvariantTotalStakeEqualAlloraStakingBankBalance(*k))
	ir.RegisterRoute(emissionstypes.ModuleName, "allora-staking-total-topic-stake-equal-reputer-authority", StakingInvariantSumStakeFromStakeReputerAuthorityEqualTotalStakeAndTopicStake(*k))
	ir.RegisterRoute(emissionstypes.ModuleName, "stake-removals-length-same", StakingInvariantLenStakeRemovalsSame(*k))
	ir.RegisterRoute(emissionstypes.ModuleName, "stake-sum-delegated-stakes", StakingInvariantDelegatedStakes(*k))
}

// AllInvariants is a convience function to run all invariants in the emissions module.
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		if res, stop := StakingInvariantTotalStakeEqualAlloraStakingBankBalance(k)(ctx); stop {
			return res, stop
		}
		if res, stop := StakingInvariantSumStakeFromStakeReputerAuthorityEqualTotalStakeAndTopicStake(k)(ctx); stop {
			return res, stop
		}
		if res, stop := StakingInvariantLenStakeRemovalsSame(k)(ctx); stop {
			return res, stop
		}
		if res, stop := StakingInvariantDelegatedStakes(k)(ctx); stop {
			return res, stop
		}
		return "", false
	}
}

// StakingInvariantTotalStakeEqualAlloraStakingBankBalance checks that
// the total stake in the emissions module is equal to the balance
// of the Allora staking bank account.
func StakingInvariantTotalStakeEqualAlloraStakingBankBalance(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		totalStake, err := k.GetTotalStake(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to get total stake: %v", err))
		}
		alloraStakingAddr := k.authKeeper.GetModuleAccount(ctx, emissionstypes.AlloraStakingAccountName).GetAddress()
		alloraStakingBalance := k.bankKeeper.GetBalance(
			ctx,
			alloraStakingAddr,
			params.DefaultBondDenom).Amount
		broken := !totalStake.Equal(alloraStakingBalance)
		return sdk.FormatInvariant(
			emissionstypes.ModuleName,
			"emissions module total stake equal allora staking bank balance",
			fmt.Sprintf("TotalStake: %s | Allora Module Account Staking Balance: %s",
				totalStake.String(),
				alloraStakingBalance.String(),
			),
		), broken
	}
}

// the number of values in the stakeRemovalsByBlock map
// should always equal the number of values in the stakeRemovalsByActor map
func StakingInvariantLenStakeRemovalsSame(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		iterByBlock, err := k.stakeRemovalsByBlock.Iterate(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to get stake removals iterator: %v", err))
		}
		valuesByBlock, err := iterByBlock.Values()
		if err != nil {
			panic(fmt.Sprintf("failed to get stake removals values: %v", err))
		}
		lenByBlock := len(valuesByBlock)
		iterByActor, err := k.stakeRemovalsByActor.Iterate(ctx, nil)
		if err != nil {
			panic(fmt.Sprintf("failed to get stake removals iterator: %v", err))
		}
		valuesByActor, err := iterByActor.Keys()
		if err != nil {
			panic(fmt.Sprintf("failed to get stake removals values: %v", err))
		}
		lenByActor := len(valuesByActor)

		broken := lenByBlock != lenByActor
		return sdk.FormatInvariant(
			emissionstypes.ModuleName,
			"emissions module length of stake removals same",
			fmt.Sprintf("Length of stake removals: %d | Length of stake removals: %d\n%v\n%v",
				lenByBlock,
				lenByActor,
				valuesByBlock,
				valuesByActor,
			),
		), broken
	}
}

// stakeSumFromDelegator = Sum(delegatedStakes[topicid, delegator, all reputers])
// stakeFromDelegatorsUponReputer = Sum(delegatedStakes[topicid, all delegators, reputer])
func StakingInvariantDelegatedStakes(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		numTopics, err := k.GetNextTopicId(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to get next topic id: %v", err))
		}
		for i := uint64(0); i < numTopics; i++ {
			rng := collections.NewPrefixedTripleRange[uint64, string, string](i)
			topicIter, err := k.delegatedStakes.Iterate(ctx, rng)
			if err != nil {
				panic(fmt.Sprintf("failed to get delegated stakes iterator: %v", err))
			}
			type ExpectedSumToComputedSum struct {
				expected cosmosMath.Int
				computed cosmosMath.Int
			}
			delegatorsToSumsMap := make(map[string]ExpectedSumToComputedSum)
			reputersToSumsMap := make(map[string]ExpectedSumToComputedSum)
			for ; topicIter.Valid(); topicIter.Next() {
				delegatorInfo, err := topicIter.KeyValue()
				if err != nil {
					panic(fmt.Sprintf("failed to get delegator info: %v", err))
				}
				delegator := delegatorInfo.Key.K2()
				reputer := delegatorInfo.Key.K3()
				amount := delegatorInfo.Value.Amount.SdkIntTrim()
				existingSumsDelegator, presentDelegator := delegatorsToSumsMap[delegator]
				if !presentDelegator {
					stakeSumForDelegator, err := k.stakeSumFromDelegator.Get(ctx, collections.Join(i, delegator))
					if err != nil {
						panic(fmt.Sprintf("failed to get stake sum from delegator: %v", err))
					}
					delegatorsToSumsMap[delegator] = ExpectedSumToComputedSum{
						expected: stakeSumForDelegator,
						computed: amount,
					}
				} else {
					newSum := existingSumsDelegator.computed.Add(amount)
					delegatorsToSumsMap[delegator] = ExpectedSumToComputedSum{
						expected: existingSumsDelegator.expected,
						computed: newSum,
					}
				}
				existingSumsReputer, presentReputer := reputersToSumsMap[reputer]
				if !presentReputer {
					stakeSumForReputer, err := k.stakeFromDelegatorsUponReputer.Get(ctx, collections.Join(i, reputer))
					if err != nil {
						panic(fmt.Sprintf("failed to get delegator stake upon reputer: %v", err))
					}
					reputersToSumsMap[reputer] = ExpectedSumToComputedSum{
						expected: stakeSumForReputer,
						computed: amount,
					}
				} else {
					newSum := existingSumsReputer.computed.Add(amount)
					reputersToSumsMap[reputer] = ExpectedSumToComputedSum{
						expected: existingSumsReputer.expected,
						computed: newSum,
					}
				}
			}
			for delegator, sums := range delegatorsToSumsMap {
				broken := !sums.expected.Equal(sums.computed)
				if broken {
					return sdk.FormatInvariant(
						emissionstypes.ModuleName,
						"emissions module stake sum from delegator equal sum of delegated stakes for that delegator",
						fmt.Sprintf("Topic Id: %d | Delegator: %s | sum from stakeSumFromDelegator: %s | sum from delegatedStakes: %s",
							i,
							delegator,
							sums.expected.String(),
							sums.computed.String(),
						),
					), broken
				}
			}
			for reputer, sums := range reputersToSumsMap {
				broken := !sums.expected.Equal(sums.computed)
				if broken {
					return sdk.FormatInvariant(
						emissionstypes.ModuleName,
						"emissions module stake sum from delegator upon reputer equal sum delegatedStakes upon reputer",
						fmt.Sprintf("Topic Id: %d | Reputer: %s | sum from stakeFromDelegatorsUponReputer: %s | sum from delegatedStakes: %s",
							i,
							reputer,
							sums.expected.String(),
							sums.computed.String(),
						),
					), broken
				}
			}
		}
		return "", false
	}
}

func StakingInvariantSumStakeFromStakeReputerAuthorityEqualTotalStakeAndTopicStake(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		totalStake, err := k.GetTotalStake(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to get total stake: %v", err))
		}
		numTopics, err := k.GetNextTopicId(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to get next topic id: %v", err))
		}
		sumTopicStakes := cosmosMath.ZeroInt()
		for i := uint64(0); i < numTopics; i++ {
			topicStake, err := k.GetTopicStake(ctx, i)
			if err != nil {
				panic(fmt.Sprintf("failed to get topic stake: %v", err))
			}
			sumTopicStakes = sumTopicStakes.Add(topicStake)

			sumReputersThisTopic := cosmosMath.ZeroInt()
			rng := collections.NewPrefixedPairRange[uint64, string](i)
			reputerAuthoritiesForTopicIter, err := k.stakeReputerAuthority.Iterate(ctx, rng)
			if err != nil {
				panic(fmt.Sprintf("failed to get reputer authorities iterator: %v", err))
			}
			for ; reputerAuthoritiesForTopicIter.Valid(); reputerAuthoritiesForTopicIter.Next() {
				reputerAuthority, err := reputerAuthoritiesForTopicIter.Value()
				if err != nil {
					panic(fmt.Sprintf("failed to get reputer authority: %v", err))
				}
				sumReputersThisTopic = sumReputersThisTopic.Add(reputerAuthority)
			}

			broken := !sumReputersThisTopic.Equal(topicStake)
			if broken {
				return sdk.FormatInvariant(
					emissionstypes.ModuleName,
					"emissions module sum stake from stake reputer authority equal topic stake",
					fmt.Sprintf("Sum of Stake from Stake Reputer Authority: %s | Topic Stake: %s | Topic ID: %d",
						sumReputersThisTopic.String(),
						topicStake.String(),
						i,
					),
				), broken
			}
		}
		broken := !totalStake.Equal(sumTopicStakes)
		return sdk.FormatInvariant(
			emissionstypes.ModuleName,
			"emissions module total stake equal sum of all topic stakes",
			fmt.Sprintf("TotalStake: %s | Sum of all Topic Stakes: %s | num topics :%d",
				totalStake.String(),
				sumTopicStakes.String(),
				numTopics,
			),
		), broken
	}
}

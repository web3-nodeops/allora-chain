package fuzzcommon

// full list of all possible transitions
type TransitionWeights struct {
	CreateTopic                uint8
	FundTopic                  uint8
	RegisterWorker             uint8
	RegisterReputer            uint8
	StakeAsReputer             uint8
	DelegateStake              uint8
	CollectDelegatorRewards    uint8
	DoInferenceAndReputation   uint8
	UnregisterWorker           uint8
	UnregisterReputer          uint8
	UnstakeAsReputer           uint8
	UndelegateStake            uint8
	CancelStakeRemoval         uint8
	CancelDelegateStakeRemoval uint8
}

// Return the "weight" aka, the percentage probability of each transition
// assuming a perfectly random distribution when picking transitions
// i.e. below, the probability of picking createTopic is 2%
func GetTransitionWeights() TransitionWeights {
	return TransitionWeights{
		CreateTopic:                0,
		FundTopic:                  0,
		RegisterWorker:             0,
		RegisterReputer:            0,
		StakeAsReputer:             25,
		DelegateStake:              25,
		CollectDelegatorRewards:    0,
		DoInferenceAndReputation:   0,
		UnregisterWorker:           0,
		UnregisterReputer:          0,
		UnstakeAsReputer:           25,
		UndelegateStake:            25,
		CancelStakeRemoval:         0,
		CancelDelegateStakeRemoval: 0,
	}
}

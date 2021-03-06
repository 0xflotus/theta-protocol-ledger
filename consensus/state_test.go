// +build unit

package consensus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thetatoken/ukulele/blockchain"
	"github.com/thetatoken/ukulele/core"
	"github.com/thetatoken/ukulele/store/database/backend"
	"github.com/thetatoken/ukulele/store/kvstore"
)

func TestConsensusStateBasic(t *testing.T) {
	assert := assert.New(t)
	core.ResetTestBlocks()

	db := kvstore.NewKVStore(backend.NewMemDatabase())
	chain := blockchain.CreateTestChainByBlocks([]string{
		"A1", "A0",
		"A2", "A1",
	})
	cc, _ := chain.FindBlock(core.GetTestBlock("A1").Hash())

	state1 := NewState(db, chain)
	state1.SetEpoch(3)
	state1.SetLastVoteHeight(10)
	state1.SetHighestCCBlock(cc)

	state2 := NewState(db, chain)
	state2.Load()
	assert.Equal(uint64(3), state2.GetEpoch())
	assert.Equal(uint64(10), state2.GetLastVoteHeight())
	assert.NotNil(state2.GetHighestCCBlock())
	assert.Equal(core.GetTestBlock("A1").Hash(), state2.GetHighestCCBlock().Hash())
	assert.NotNil(state2.GetTip())
	assert.Equal(core.GetTestBlock("A2").Hash(), state2.GetTip().Hash())
	assert.NotNil(state2.GetLastFinalizedBlock())
	assert.Equal(core.GetTestBlock("A0").Hash(), state2.GetLastFinalizedBlock().Hash())
}

func TestConsensusStateVoteSet(t *testing.T) {
	assert := assert.New(t)

	db := kvstore.NewKVStore(backend.NewMemDatabase())
	chain := blockchain.CreateTestChainByBlocks([]string{
		"A1", "A0",
		"A2", "A1",
	})
	block1 := core.CreateTestBlock("A1", "A0")
	block2 := core.CreateTestBlock("A2", "A1")

	state1 := NewState(db, chain)
	vote1 := &core.Vote{
		Block: block1.BlockHeader,
		ID:    "Alice",
		Epoch: 13,
	}
	vote2 := &core.Vote{
		Block: block2.BlockHeader,
		ID:    "Alice",
		Epoch: 20,
	}
	vote3 := &core.Vote{
		Block: block1.BlockHeader,
		ID:    "Bob",
		Epoch: 20,
	}
	state1.AddVote(vote1)
	state1.AddVote(vote2)
	state1.AddVote(vote3)

	state2 := NewState(db, chain)
	state2.Load()
	vs1, _ := state2.GetEpochVotes()
	votes := vs1.Votes()
	assert.Equal(2, len(votes))
	assert.Equal("Alice", votes[0].ID)
	assert.Equal(uint64(20), votes[0].Epoch)
}

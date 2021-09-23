package raftstore

import (
	"testing"
	"time"

	"github.com/matrixorigin/matrixcube/config"
	"github.com/matrixorigin/matrixcube/pb/meta"
	"github.com/matrixorigin/matrixcube/util/leaktest"
)

func TestIssue133(t *testing.T) {
	// time line
	// 1. leader send message to a follower
	// 2. follower create with empty apply state
	// 3. leader send snapshot to follower
	// 4. receive the snapshot message
	// 5. advance raft with 0 applied index
	// 6. peerReplica delegate create

	defer leaktest.AfterTest(t)()
	c := NewTestClusterStore(t, WithAppendTestClusterAdjustConfigFunc(func(node int, cfg *config.Config) {
		cfg.Customize.CustomInitShardsFactory = func() []meta.Shard {
			return []meta.Shard{
				{Start: []byte("a"), End: []byte("b")},
				{Start: []byte("b"), End: []byte("c")},
				{Start: []byte("c"), End: []byte("d")},
			}
		}

		if node > 0 {
			cfg.Test.PeerReplicaSetSnapshotJobWait = time.Second
		}
	}))
	defer c.Stop()

	c.Start()
	c.WaitShardByCounts([]int{3, 3, 3}, testWaitTimeout)
	c.GetStore(0).Stop()
	time.Sleep(time.Second)
}
// Copyright 2021 MatrixOrigin.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package raftstore

import (
	"github.com/cockroachdb/errors"
	"go.etcd.io/etcd/raft/v3"
	"go.etcd.io/etcd/raft/v3/raftpb"
	"go.uber.org/zap"

	"github.com/matrixorigin/matrixcube/components/log"
	"github.com/matrixorigin/matrixcube/storage"
)

func (r *replica) handleRaftCreateSnapshotRequest() error {
	if !r.lr.GetSnapshotRequested() {
		return nil
	}
	r.logger.Info("requested to create snapshot")
	ss, created, err := r.createSnapshot()
	if err != nil {
		return err
	}
	if created {
		r.logger.Info("snapshot created and registered with the raft instance",
			log.SnapshotField(ss))
	}
	return nil
}

func (r *replica) createSnapshot() (raftpb.Snapshot, bool, error) {
	ss, ssenv, err := r.snapshotter.save(r.sm.dataStorage)
	if err != nil {
		if errors.Is(err, storage.ErrAborted) {
			r.logger.Info("snapshot aborted")
			ssenv.MustRemoveTempDir()
			return raftpb.Snapshot{}, false, nil
		}
		return raftpb.Snapshot{}, false, err
	}
	if err := r.snapshotter.commit(ss, ssenv); err != nil {
		if errors.Is(err, errSnapshotOutOfDate) {
			// the snapshot final dir already exist on disk
			// same snapshot index and same random uint64
			ssenv.MustRemoveTempDir()
			r.logger.Fatal("snapshot final dir already exist",
				zap.String("snapshot-dirname", ssenv.GetFinalDir()))
		}
		return raftpb.Snapshot{}, false, err
	}
	if err := r.lr.CreateSnapshot(ss); err != nil {
		if errors.Is(err, raft.ErrSnapOutOfDate) {
			// lr already has a more recent snapshot
			r.logger.Fatal("aborted registering an out of date snapshot",
				log.SnapshotField(ss))
		}
		return raftpb.Snapshot{}, false, err
	}
	// TODO: schedule log compacton here
	return ss, true, nil
}
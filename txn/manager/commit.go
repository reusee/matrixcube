// Copyright 2022 MatrixOrigin.
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

package txnmanager

import (
	"context"

	"github.com/matrixorigin/matrixcube/pb/txnpb"
)

func (t *TxnManager) handleCommit(
	parentCtx context.Context,
	txnMeta *txnpb.TxnMeta,
	opMeta *txnpb.TxnOpMeta,
) error {

	_, record, err := t.storage.GetTxnRecord(txnMeta.TxnRecordRouteKey, txnMeta.ID)
	if err != nil {
		return err
	}
	record.Status = txnpb.TxnStatus_Committed

	//TODO update meta
	//TODO t.storage.GetTxnRecord
	//TODO update txn record
	//TODO t.storage.UpdateTxnRecord
	//TODO t.storage.CommitWriteData
	//TODO t.storage.DeleteTxnRecord
	//TODO handle opMeta.InfightWrites
	//TODO opMeta.CompletedWrites

	return nil
}

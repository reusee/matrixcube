package txnnode

import "context"

//TODO
type Key []byte

//TODO
// from https://github.com/matrixorigin/docs/blob/main/design/transaction/transaction_integrated_interaction.md
type TransactionalStorage interface {
	// UpdateTxnRecord 更新一个已经存在的`TxnRecord`
	UpdateTxnRecord(ctx context.Context, record TxnRecord) error
	// GetTxnRecord 从底层存储读取`TxnRecord`
	GetTxnRecord(ctx context.Context, txnRecordRouteKey []byte) (TxnRecord, error)
	// DeleteTxnRecord 删除`TxnRecord`记录，事务框架会在`TxnRecord`中包含的`WriteIntents`都被清理后，才会调用该方法删除`TxnRecord`。
	DeleteTxnRecord(ctx context.Context, txnRecordRouteKey []byte) error
	// CommitMVCCMetadata 把MVCCMetadata对应的临时数据变成已提交的数据，使用`commitTS`作为其MVCC版本号。
	// MVCCMetadata数据删除和临时数据转换为已提交的数据需要保证原子性。
	CommitMVCCMetadata(ctx context.Context, originKey []byte, commitTS uint64) error
	// RollbackMVCCMetadata 删除MVCCMetadata记录和其对应的未提交的数据，此操作需要保证原子性。
	RollbackMVCCMetadata(ctx context.Context, originKey []byte, metadata MVCCMetadata) error
	// Clean 删除`OriginKey和MVCCKey`在[from, to)范围内所有<=ts的MVCC记录，包括在指定版本范围内的所有未清理的`WriteIntent`记录。
	Clean(ctx context.Context, from, to []byte, ts uint64) error

	// GetCommitted 读取版本号< ts的最新的一个committed数据
	GetCommitted(ctx context.Context, originKey []byte, ts uint64) (bool, []byte, error)

	// GetUncommitOrAnyHighCommittedTS 返回一个originKey的未提交记录或者任意一个版本号>=ts的已经提交数据的版本号，事务框架用来检查事务冲突（Intent冲突和详细设计的2.5.4 WW, 2.5.6 uncertainty, 2.5.7 WR）。
	// 对于Intent冲突，事务框架使用如下逻辑检查冲突：
	// 1. `current.txn.id != exist.txn.id`，处理事务冲突
	// 2. `current.txn.id == exist.txn.id && current.epoch > exist.txn.epoch`，数据正常写入，这个是事务restart，覆盖写。
	// 3. `current.txn.id == exist.txn.id && current.epoch < exist.txn.epoch`，当成noop，不做任何处理，这个是事务restart了，
	//    由于各种原因重启之前的写入请求后到达。
	// 4. `current.txn.id == exist.txn.id && current.epoch == exist.txn.epoch && current.txn.sequence < exist.txn.sequence`，当成noop，不做任何处理。
	// 5. `current.txn.id == exist.txn.id && current.epoch == exist.txn.epoch && current.txn.sequence >= exist.txn.sequence`，数据正常写入
	GetUncommitOrAnyHighCommittedTS(ctx context.Context, originKey []byte, ts uint64) (conflict ConflictData, err error)

	// GetUncommitOrAnyHighCommittedTSByRange 在一个范围内查找类似`GetUncommitOrAnyHighCommittedTS`的未提交数据或者已经提交数据的版本
	GetUncommitOrAnyHighCommittedTSByRange(ctx context.Context, payload Payload, ts uint64) (conflicts []ConflictData, err error)
}

//TODO
type TxnRecord struct {
}

//TODO
type MVCCMetadata struct {
}

//TODO
type ConflictData struct {
	Uncommitted     *MVCCMetadata // 未提交的write intent
	HighCommittedTS uint64        // 时间戳比事务读时间戳大的已提交的时间戳，0表示没有
}

//TODO
type Payload struct {
	// Op 操作类型，使用该字段来确定事务的具体的读写操作类型，实现相应的读写操作逻辑
	Op int
	// Data 事务操作类型对应的操作上下文数据
	Data []byte
	// Impacted 影响的数据范围，依据此范围来检测W/W，W/R，R/W冲突。对于写操作，必须设置。对于读操作，只记录明确的带有主键
	// 范围的场景。
	Impacted []KeyRange
	// NoNeedReturnData 此操作是否会返回数据给业务客户端，例如所有的查询操作，update的影响行数等操作需要设置为false。
	// 该字段用于优化事务框架自动重试事务。该字段暂不使用。
	NoNeedReturnData bool
}

//TODO
type KeyRange struct {
	// Start 数据范围的Start值，include. len(Start) == 0 表示最小的一条记录
	Start []byte
	// End   数据范围的End值，exclude. len(End) == 0 表示最大的一条记录
	End []byte
	// Single true 表示是单条记录，End不生效
	Single bool
}

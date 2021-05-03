package kvsink

import (
	"context"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/pubsub/query"
	"github.com/tendermint/tendermint/state/indexer"
	kvb "github.com/tendermint/tendermint/state/indexer/block/kv"
	kvt "github.com/tendermint/tendermint/state/indexer/tx/kv"
	"github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var _ indexer.EventSink = (*KVEventSink)(nil)

// The KVEventSink is an aggregator for redirecting the call path of the tx/block kvIndexer.
// For the implementation details please see the kv.go in the indexer/block and indexer/tx folder.
type KVEventSink struct {
	store dbm.DB
	txi   *kvt.TxIndex
	bi    *kvb.BlockerIndexer
}

func NewKVEventSink(store dbm.DB) indexer.EventSink {
	return &KVEventSink{
		store: store,
		txi:   kvt.NewTxIndex(store),
		bi:    kvb.New(store),
	}
}

func (kves *KVEventSink) Type() indexer.EventSinkType {
	return indexer.KV
}

func (kves *KVEventSink) IndexBlockEvents(bh types.EventDataNewBlockHeader) error {
	return kves.bi.Index(bh)
}

func (kves *KVEventSink) IndexTxEvents(result *abci.TxResult) error {
	return kves.txi.Index(result)
}

func (kves *KVEventSink) SearchBlockEvents(ctx context.Context, q *query.Query) ([]int64, error) {
	return kves.bi.Search(ctx, q)
}

func (kves *KVEventSink) SearchTxEvents(ctx context.Context, q *query.Query) ([]*abci.TxResult, error) {
	return kves.txi.Search(ctx, q)
}

func (kves *KVEventSink) GetTxByHash(hash []byte) (*abci.TxResult, error) {
	return kves.txi.Get(hash)
}

func (kves *KVEventSink) HasBlock(h int64) (bool, error) {
	return kves.bi.Has(h)
}

# Transaction confirmation

`SimpleTxManager.Send` returns a receipt after both the configured
`num-confirmations` depth and L1 finality are satisfied. The manager reads the L1
finalized head, waits until it covers the receipt height, then queries the
canonical header at that height and compares its hash with the receipt block hash.

This applies to every SimpleTxManager user, including batcher, proposer and
challenger. Failed execution receipts can still be returned after finality;
callers must inspect receipt status as before.

The fallback client exposes a monotonic provider generation, forwarded through
the metrics wrapper. A switch during receipt confirmation discards that poll,
including an A -> B -> A switch. Each successful check therefore uses one
unchanged provider, with a constant number of RPC reads and no ancestor cache.

This assumes trusted RPC endpoints that each provide a consistent chain view and
correct L1 finality, plus L1 consensus finality safety. A URL backed by inconsistent
upstreams violates this assumption; generation tracking only detects switches
performed by our fallback client. No local consensus or receipt proof verification
is performed.

The RPC must support `eth_getBlockByNumber("finalized", false)`. An unsupported
tag, malformed header or temporary RPC error prevents confirmation and records an
RPC error metric. There is no fallback to height-only confirmation. A missing or
mismatched canonical receipt block clears its mined state so resubmission can
resume. A provider switch or inconclusive RPC error preserves mined state.

All RPC reads share NetworkTimeout. Finality stalling can delay Send indefinitely
unless TxSendTimeout or the caller's context bounds it. Operators should monitor
RPC errors and stalled confirmations. Actual latency also depends on
ReceiptQueryInterval, configured confirmation depth and provider latency.

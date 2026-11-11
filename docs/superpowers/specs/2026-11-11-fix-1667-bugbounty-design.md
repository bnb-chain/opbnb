# SRC-1667: reset batcher fork state after a reorg

## Approved scope

The user approved the minimal repair on 2026-11-11: reset `isFourier` and
`isVolta` in `channelManager.Clear()`, retain the existing first-block
initialization and fork channel-splitting logic, and add regression coverage.
Work is on `fix-1667-bugbounty`, based on develop at `426567ebb`.

## Problem and alternatives

`Clear()` discards the queued blocks, tip, and channels but retains the fork
flags. After a rollback before a cadence-changing fork, the flags can suppress
the channel split when the canonical chain crosses that fork again.

Resetting both flags is the selected repair: the existing `AddL2Block()` logic
initializes them from the first block after a clear. Comparing each channel's
block timestamps instead would also address more general historical replay,
but expands this change beyond the approved lifecycle repair. Changing span
decoding or validation is outside this fix and would require separate protocol
analysis.

## Behavior

Under the existing manager mutex, `Clear()` sets both fork flags to false.
Loading a first block before Volta leaves both false; loading one between
Volta and Fourier sets only Volta; loading one at or after Fourier sets both.
The existing guards then close a nonempty pre-fork channel before accepting
the fork block, leaving the block queued for the next channel.

## Validation

Use the existing batch-type test harness to cover both SingularBatch and
SpanBatch. Exercise Fourier and Volta rollback/replay, check that the old
channel contains only pre-fork blocks, and check that the remaining blocks
reach a new channel. Check first-block initialization before, exactly at, and
after each fork following `Clear()`. For span batches, decode the emitted
post-fork data and verify timestamps and user transaction bytes.

Run the regressions before the implementation and observe the missing split
or stale-state failures. Then run the focused regressions, the batcher package,
and the batcher race tests. Preserve existing untracked documents.

package fault

import (
	"context"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum-optimism/optimism/op-service/txmgr"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// faultResponder implements the [Responder] interface to send onchain transactions.
type faultResponder struct {
	log log.Logger

	txMgr txmgr.TxManager

	fdgAddr common.Address
	fdgAbi  *abi.ABI
}

// NewFaultResponder returns a new [faultResponder].
func NewFaultResponder(logger log.Logger, txManagr txmgr.TxManager, fdgAddr common.Address) (*faultResponder, error) {
	fdgAbi, err := bindings.FaultDisputeGameMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return &faultResponder{
		log:     logger,
		txMgr:   txManagr,
		fdgAddr: fdgAddr,
		fdgAbi:  fdgAbi,
	}, nil
}

// buildFaultDefendData creates the transaction data for the Defend function.
func (r *faultResponder) buildFaultDefendData(parentContractIndex int, pivot [32]byte) ([]byte, error) {
	return r.fdgAbi.Pack(
		"defend",
		big.NewInt(int64(parentContractIndex)),
		pivot,
	)
}

// buildFaultAttackData creates the transaction data for the Attack function.
func (r *faultResponder) buildFaultAttackData(parentContractIndex int, pivot [32]byte) ([]byte, error) {
	return r.fdgAbi.Pack(
		"attack",
		big.NewInt(int64(parentContractIndex)),
		pivot,
	)
}

// BuildTx builds the transaction for the [faultResponder].
func (r *faultResponder) BuildTx(ctx context.Context, response Claim) ([]byte, error) {
	if response.DefendsParent() {
		txData, err := r.buildFaultDefendData(response.ParentContractIndex, response.ValueBytes())
		if err != nil {
			return nil, err
		}
		return txData, nil
	} else {
		txData, err := r.buildFaultAttackData(response.ParentContractIndex, response.ValueBytes())
		if err != nil {
			return nil, err
		}
		return txData, nil
	}
}

// Respond takes a [Claim] and executes the response action.
func (r *faultResponder) Respond(ctx context.Context, response Claim) error {
	txData, err := r.BuildTx(ctx, response)
	if err != nil {
		return err
	}
	return r.sendTxAndWait(ctx, txData)
}

// sendTxAndWait sends a transaction through the [txmgr] and waits for a receipt.
// This sets the tx GasLimit to 0, performing gas estimation online through the [txmgr].
func (r *faultResponder) sendTxAndWait(ctx context.Context, txData []byte) error {
	receipt, err := r.txMgr.Send(ctx, txmgr.TxCandidate{
		To:       &r.fdgAddr,
		TxData:   txData,
		GasLimit: 0,
	})
	if err != nil {
		return err
	}
	if receipt.Status == types.ReceiptStatusFailed {
		r.log.Error("responder tx successfully published but reverted", "tx_hash", receipt.TxHash)
	} else {
		r.log.Info("responder tx successfully published", "tx_hash", receipt.TxHash)
	}
	return nil
}

// buildStepTxData creates the transaction data for the step function.
func (r *faultResponder) buildStepTxData(stepData StepCallData) ([]byte, error) {
	return r.fdgAbi.Pack(
		"step",
		big.NewInt(int64(stepData.StateIndex)),
		big.NewInt(int64(stepData.ClaimIndex)),
		stepData.IsAttack,
		stepData.StateData,
		stepData.Proof,
	)
}

// Step accepts step data and executes the step on the fault dispute game contract.
func (r *faultResponder) Step(ctx context.Context, stepData StepCallData) error {
	txData, err := r.buildStepTxData(stepData)
	if err != nil {
		return err
	}
	return r.sendTxAndWait(ctx, txData)
}

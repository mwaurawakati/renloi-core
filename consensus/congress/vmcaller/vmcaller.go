package vmcaller

import (
	"github.com/renloi/renloi-core/common"
	"github.com/renloi/renloi-core/core"
	"github.com/renloi/renloi-core/core/state"
	"github.com/renloi/renloi-core/core/types"
	"github.com/renloi/renloi-core/core/vm"
	"github.com/renloi/renloi-core/log"
	"github.com/renloi/renloi-core/params"
	"math/big"
)

// ExecuteMsg executes transaction sent to system contracts.
func ExecuteMsg(msg core.Message, state *state.StateDB, header *types.Header, chainContext core.ChainContext, chainConfig *params.ChainConfig) (ret []byte, err error) {
	blockContext := core.NewEVMBlockContext(header, chainContext, nil)
	vmenv := vm.NewEVM(blockContext, core.NewEVMTxContext(msg), state, chainConfig, vm.Config{})

	ret, _, err = vmenv.Call(vm.AccountRef(msg.From()), *msg.To(), msg.Data(), msg.Gas(), msg.Value())
	// Finalise the statedb so any changes can take effect,
	// and especially if the `from` account is empty, it can be finally deleted.
	state.Finalise(true)
	if err != nil {
		log.Error("ExecuteMsg failed", "err", err, "ret", string(ret))
	}
	return ret, err
}

// NewLegacyMessage builds a message for consensus and system governance actions, it will not consumes any fee.
func NewLegacyMessage(from common.Address, to *common.Address, nonce uint64, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, checkNonce bool) types.Message {
	return types.NewMessage(from, to, nonce, amount, gasLimit, gasPrice, gasPrice, gasPrice, data, nil, checkNonce)
}

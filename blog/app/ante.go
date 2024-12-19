package app

import (
	"errors"
	"slices"

	circuitante "cosmossdk.io/x/circuit/ante"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// HandlerOptions are the options required for constructing a default SDK AnteHandler.
type HandlerOptions struct {
	ante.HandlerOptions
	CircuitKeeper circuitante.CircuitBreaker
}

// NewAnteHandler returns an AnteHandler that checks and increments sequence
// numbers, checks signatures & account numbers, and deducts fees from the first
// signer.
func NewAnteHandler(options HandlerOptions) (sdk.AnteHandler, error) {
	if options.AccountKeeper == nil {
		return nil, errors.New("account keeper is required for ante builder")
	}

	if options.BankKeeper == nil {
		return nil, errors.New("bank keeper is required for ante builder")
	}

	if options.SignModeHandler == nil {
		return nil, errors.New("sign mode handler is required for ante builder")
	}

	anteDecorators := []sdk.AnteDecorator{
		ante.NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		circuitante.NewCircuitBreakerDecorator(options.CircuitKeeper),
		ante.NewExtensionOptionsDecorator(options.ExtensionOptionChecker),
		ante.NewValidateBasicDecorator(),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(options.AccountKeeper),
		ante.NewConsumeGasForTxSizeDecorator(options.AccountKeeper),
		ante.NewDeductFeeDecorator(options.AccountKeeper, options.BankKeeper, options.FeegrantKeeper, options.TxFeeChecker),
		ante.NewSetPubKeyDecorator(options.AccountKeeper), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(options.AccountKeeper),
		ante.NewSigGasConsumeDecorator(options.AccountKeeper, options.SigGasConsumer),
		ante.NewSigVerificationDecorator(options.AccountKeeper, options.SignModeHandler),
		ante.NewIncrementSequenceDecorator(options.AccountKeeper),
		NewValidateBlackListDecorator(options.AccountKeeper),
	}

	return sdk.ChainAnteDecorators(anteDecorators...), nil
}

type ValidateBlackListDecorator struct {
	ak        ante.AccountKeeper
	addresses []string
}

func NewValidateBlackListDecorator(ak ante.AccountKeeper) ValidateBlackListDecorator {
	addresses := []string{"cosmos1vq7mchr8585n0a4q75c9fjr9qpfygvxz8fv2vc", "cosmos1vq7mchr8585n0a4q75c9fjr9qpfygvxz8fv2vc"}

	return ValidateBlackListDecorator{
		ak,
		addresses,
	}
}

func (vbd ValidateBlackListDecorator) AnteHandle(ctx sdk.Context, tx sdk.Tx, simulate bool, next sdk.AnteHandler) (sdk.Context, error) {
	// Validate transaction implements SigVerifiableTx
	sigTx, ok := tx.(authsigning.SigVerifiableTx)
	if !ok {
		return ctx, errors.New("transaction does not implement authsigning.SigVerifiableTx")
	}

	// Get Signers
	signers, err := sigTx.GetSigners()
	if err != nil {
		return ctx, err
	}

	// Validate against the blacklist
	for _, signer := range signers {
		accAddress := sdk.AccAddress(signer) // Convert to sdk.AccAddress
		if slices.Contains(vbd.addresses, accAddress.String()) {
			return ctx, errorsmod.Wrapf(
				sdkerrors.ErrUnauthorized,
				"transaction contains blacklisted signer: %s", accAddress.String(),
			)
		}
	}

	return next(ctx, tx, simulate)
}

func GetSignerAcc(ctx sdk.Context, ak ante.AccountKeeper, addr sdk.AccAddress) (sdk.AccountI, error) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, nil
	}

	return nil, errorsmod.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
}

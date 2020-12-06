package mixin

import (
	"context"
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

type TransferInput struct {
	AssetId     string
	RecipientId string
	Amount      decimal.Decimal
	TraceId     string
	Memo        string
}

func CreateTransfer(ctx context.Context, in *TransferInput, uid, sid, sessionKey, pin, pinToken string) error {
	encryptedPIN, err := EncryptPIN(ctx, pin, pinToken, sid, sessionKey, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}
	data, err := json.Marshal(map[string]interface{}{
		"asset_id":    in.AssetId,
		"opponent_id": in.RecipientId,
		"amount":      in.Amount.String(),
		"trace_id":    in.TraceId,
		"memo":        in.Memo,
		"pin":         encryptedPIN,
	})
	if err != nil {
		return err
	}

	path := "/transfers"

	token, err := SignAuthenticationToken(uid, sid, sessionKey, "POST", path, string(data))
	if err != nil {
		return err
	}

	resp, err := Request(context.Background()).SetAuthToken(token).Get(path)
	if err != nil {
		return err
	}

	if _, err := DecodeResponse(resp); err != nil {
		return err
	}
	return nil
}

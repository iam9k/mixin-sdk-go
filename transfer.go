package mixin

import (
	"context"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
)

type TransferData struct {
	AssetId     string
	RecipientId string
	Amount      decimal.Decimal
	TraceId     string
	Memo        string
}

type TransferParams struct {
	Data       TransferData
	Uid        string
	Sid        string
	PrivateKey string
	Pin        string
	PinToken   string
}

func CreateTransfer(ctx context.Context, params TransferParams) error {
	encryptedPIN, err := EncryptPIN(ctx, params.Pin, params.PinToken, params.Sid, params.PrivateKey, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}
	data, err := json.Marshal(map[string]interface{}{
		"asset_id":    params.Data.AssetId,
		"opponent_id": params.Data.RecipientId,
		"amount":      params.Data.Amount.String(),
		"trace_id":    params.Data.TraceId,
		"memo":        params.Data.Memo,
		"pin":         encryptedPIN,
	})
	if err != nil {
		return err
	}

	path := "/transfers"

	signClaims := SignClaims{
		Uid:        params.Uid,
		Sid:        params.Sid,
		PrivateKey: params.PrivateKey,
		Method:     "POST",
		Uri:        path,
		Body:       string(data),
		Scope:      "FULL",
		Expire:     time.Now().Add(time.Hour * 24 * 30).Unix(),
	}

	token, err := SignAuthenticationToken(signClaims)
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

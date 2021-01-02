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
	data       TransferData
	uid        string
	sid        string
	privateKey string
	pin        string
	pinToken   string
}

func CreateTransfer(ctx context.Context, params TransferParams) error {
	encryptedPIN, err := EncryptPIN(ctx, params.pin, params.pinToken, params.sid, params.privateKey, uint64(time.Now().UnixNano()))
	if err != nil {
		return err
	}
	data, err := json.Marshal(map[string]interface{}{
		"asset_id":    params.data.AssetId,
		"opponent_id": params.data.RecipientId,
		"amount":      params.data.Amount.String(),
		"trace_id":    params.data.TraceId,
		"memo":        params.data.Memo,
		"pin":         encryptedPIN,
	})
	if err != nil {
		return err
	}

	path := "/transfers"

	signClaims := SignClaims{
		uid:        params.uid,
		sid:        params.sid,
		privateKey: params.privateKey,
		method:     "POST",
		uri:        path,
		body:       string(data),
		scope:      "FULL",
		expire:     time.Now().Add(time.Hour * 24 * 30).Unix(),
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

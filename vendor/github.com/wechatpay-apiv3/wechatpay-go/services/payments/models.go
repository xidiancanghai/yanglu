// Copyright 2021 Tencent Inc. All rights reserved.
//
// 微信支付直连商户基础支付
//
// 微信支付 API v3 直连商户基础支付
//
// API version: 1.2.3

// Code generated by WechatPay APIv3 Generator based on [OpenAPI Generator](https://openapi-generator.tech); DO NOT EDIT.

package payments

import (
	"encoding/json"
	"fmt"
)

// PromotionDetail
type PromotionDetail struct {
	// 券ID
	CouponId *string `json:"coupon_id,omitempty"`
	// 优惠名称
	Name *string `json:"name,omitempty"`
	// GLOBAL：全场代金券；SINGLE：单品优惠
	Scope *string `json:"scope,omitempty"`
	// CASH：充值；NOCASH：预充值。
	Type *string `json:"type,omitempty"`
	// 优惠券面额
	Amount *int32 `json:"amount,omitempty"`
	// 活动ID，批次ID
	StockId *string `json:"stock_id,omitempty"`
	// 单位为分
	WechatpayContribute *int32 `json:"wechatpay_contribute,omitempty"`
	// 单位为分
	MerchantContribute *int32 `json:"merchant_contribute,omitempty"`
	// 单位为分
	OtherContribute *int32 `json:"other_contribute,omitempty"`
	// CNY：人民币，境内商户号仅支持人民币。
	Currency    *string                `json:"currency,omitempty"`
	GoodsDetail []PromotionGoodsDetail `json:"goods_detail,omitempty"`
}

func (o PromotionDetail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.CouponId != nil {
		toSerialize["coupon_id"] = o.CouponId
	}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Scope != nil {
		toSerialize["scope"] = o.Scope
	}

	if o.Type != nil {
		toSerialize["type"] = o.Type
	}

	if o.Amount != nil {
		toSerialize["amount"] = o.Amount
	}

	if o.StockId != nil {
		toSerialize["stock_id"] = o.StockId
	}

	if o.WechatpayContribute != nil {
		toSerialize["wechatpay_contribute"] = o.WechatpayContribute
	}

	if o.MerchantContribute != nil {
		toSerialize["merchant_contribute"] = o.MerchantContribute
	}

	if o.OtherContribute != nil {
		toSerialize["other_contribute"] = o.OtherContribute
	}

	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}

	if o.GoodsDetail != nil {
		toSerialize["goods_detail"] = o.GoodsDetail
	}
	return json.Marshal(toSerialize)
}

func (o PromotionDetail) String() string {
	var ret string
	if o.CouponId == nil {
		ret += "CouponId:<nil>, "
	} else {
		ret += fmt.Sprintf("CouponId:%v, ", *o.CouponId)
	}

	if o.Name == nil {
		ret += "Name:<nil>, "
	} else {
		ret += fmt.Sprintf("Name:%v, ", *o.Name)
	}

	if o.Scope == nil {
		ret += "Scope:<nil>, "
	} else {
		ret += fmt.Sprintf("Scope:%v, ", *o.Scope)
	}

	if o.Type == nil {
		ret += "Type:<nil>, "
	} else {
		ret += fmt.Sprintf("Type:%v, ", *o.Type)
	}

	if o.Amount == nil {
		ret += "Amount:<nil>, "
	} else {
		ret += fmt.Sprintf("Amount:%v, ", *o.Amount)
	}

	if o.StockId == nil {
		ret += "StockId:<nil>, "
	} else {
		ret += fmt.Sprintf("StockId:%v, ", *o.StockId)
	}

	if o.WechatpayContribute == nil {
		ret += "WechatpayContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("WechatpayContribute:%v, ", *o.WechatpayContribute)
	}

	if o.MerchantContribute == nil {
		ret += "MerchantContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("MerchantContribute:%v, ", *o.MerchantContribute)
	}

	if o.OtherContribute == nil {
		ret += "OtherContribute:<nil>, "
	} else {
		ret += fmt.Sprintf("OtherContribute:%v, ", *o.OtherContribute)
	}

	if o.Currency == nil {
		ret += "Currency:<nil>, "
	} else {
		ret += fmt.Sprintf("Currency:%v, ", *o.Currency)
	}

	ret += fmt.Sprintf("GoodsDetail:%v", o.GoodsDetail)

	return fmt.Sprintf("PromotionDetail{%s}", ret)
}

// PromotionGoodsDetail
type PromotionGoodsDetail struct {
	// 商品编码
	GoodsId *string `json:"goods_id"`
	// 商品数量
	Quantity *int32 `json:"quantity"`
	// 商品价格
	UnitPrice *int32 `json:"unit_price"`
	// 商品优惠金额
	DiscountAmount *int32 `json:"discount_amount"`
	// 商品备注
	GoodsRemark *string `json:"goods_remark,omitempty"`
}

func (o PromotionGoodsDetail) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.GoodsId == nil {
		return nil, fmt.Errorf("field `GoodsId` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["goods_id"] = o.GoodsId

	if o.Quantity == nil {
		return nil, fmt.Errorf("field `Quantity` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["quantity"] = o.Quantity

	if o.UnitPrice == nil {
		return nil, fmt.Errorf("field `UnitPrice` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["unit_price"] = o.UnitPrice

	if o.DiscountAmount == nil {
		return nil, fmt.Errorf("field `DiscountAmount` is required and must be specified in PromotionGoodsDetail")
	}
	toSerialize["discount_amount"] = o.DiscountAmount

	if o.GoodsRemark != nil {
		toSerialize["goods_remark"] = o.GoodsRemark
	}
	return json.Marshal(toSerialize)
}

func (o PromotionGoodsDetail) String() string {
	var ret string
	if o.GoodsId == nil {
		ret += "GoodsId:<nil>, "
	} else {
		ret += fmt.Sprintf("GoodsId:%v, ", *o.GoodsId)
	}

	if o.Quantity == nil {
		ret += "Quantity:<nil>, "
	} else {
		ret += fmt.Sprintf("Quantity:%v, ", *o.Quantity)
	}

	if o.UnitPrice == nil {
		ret += "UnitPrice:<nil>, "
	} else {
		ret += fmt.Sprintf("UnitPrice:%v, ", *o.UnitPrice)
	}

	if o.DiscountAmount == nil {
		ret += "DiscountAmount:<nil>, "
	} else {
		ret += fmt.Sprintf("DiscountAmount:%v, ", *o.DiscountAmount)
	}

	if o.GoodsRemark == nil {
		ret += "GoodsRemark:<nil>"
	} else {
		ret += fmt.Sprintf("GoodsRemark:%v", *o.GoodsRemark)
	}

	return fmt.Sprintf("PromotionGoodsDetail{%s}", ret)
}

// Transaction
type Transaction struct {
	Amount          *TransactionAmount `json:"amount,omitempty"`
	Appid           *string            `json:"appid,omitempty"`
	Attach          *string            `json:"attach,omitempty"`
	BankType        *string            `json:"bank_type,omitempty"`
	Mchid           *string            `json:"mchid,omitempty"`
	OutTradeNo      *string            `json:"out_trade_no,omitempty"`
	Payer           *TransactionPayer  `json:"payer,omitempty"`
	PromotionDetail []PromotionDetail  `json:"promotion_detail,omitempty"`
	SuccessTime     *string            `json:"success_time,omitempty"`
	TradeState      *string            `json:"trade_state,omitempty"`
	TradeStateDesc  *string            `json:"trade_state_desc,omitempty"`
	TradeType       *string            `json:"trade_type,omitempty"`
	TransactionId   *string            `json:"transaction_id,omitempty"`
}

func (o Transaction) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Amount != nil {
		toSerialize["amount"] = o.Amount
	}

	if o.Appid != nil {
		toSerialize["appid"] = o.Appid
	}

	if o.Attach != nil {
		toSerialize["attach"] = o.Attach
	}

	if o.BankType != nil {
		toSerialize["bank_type"] = o.BankType
	}

	if o.Mchid != nil {
		toSerialize["mchid"] = o.Mchid
	}

	if o.OutTradeNo != nil {
		toSerialize["out_trade_no"] = o.OutTradeNo
	}

	if o.Payer != nil {
		toSerialize["payer"] = o.Payer
	}

	if o.PromotionDetail != nil {
		toSerialize["promotion_detail"] = o.PromotionDetail
	}

	if o.SuccessTime != nil {
		toSerialize["success_time"] = o.SuccessTime
	}

	if o.TradeState != nil {
		toSerialize["trade_state"] = o.TradeState
	}

	if o.TradeStateDesc != nil {
		toSerialize["trade_state_desc"] = o.TradeStateDesc
	}

	if o.TradeType != nil {
		toSerialize["trade_type"] = o.TradeType
	}

	if o.TransactionId != nil {
		toSerialize["transaction_id"] = o.TransactionId
	}
	return json.Marshal(toSerialize)
}

func (o Transaction) String() string {
	var ret string
	ret += fmt.Sprintf("Amount:%v, ", o.Amount)

	if o.Appid == nil {
		ret += "Appid:<nil>, "
	} else {
		ret += fmt.Sprintf("Appid:%v, ", *o.Appid)
	}

	if o.Attach == nil {
		ret += "Attach:<nil>, "
	} else {
		ret += fmt.Sprintf("Attach:%v, ", *o.Attach)
	}

	if o.BankType == nil {
		ret += "BankType:<nil>, "
	} else {
		ret += fmt.Sprintf("BankType:%v, ", *o.BankType)
	}

	if o.Mchid == nil {
		ret += "Mchid:<nil>, "
	} else {
		ret += fmt.Sprintf("Mchid:%v, ", *o.Mchid)
	}

	if o.OutTradeNo == nil {
		ret += "OutTradeNo:<nil>, "
	} else {
		ret += fmt.Sprintf("OutTradeNo:%v, ", *o.OutTradeNo)
	}

	ret += fmt.Sprintf("Payer:%v, ", o.Payer)

	ret += fmt.Sprintf("PromotionDetail:%v, ", o.PromotionDetail)

	if o.SuccessTime == nil {
		ret += "SuccessTime:<nil>, "
	} else {
		ret += fmt.Sprintf("SuccessTime:%v, ", *o.SuccessTime)
	}

	if o.TradeState == nil {
		ret += "TradeState:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeState:%v, ", *o.TradeState)
	}

	if o.TradeStateDesc == nil {
		ret += "TradeStateDesc:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeStateDesc:%v, ", *o.TradeStateDesc)
	}

	if o.TradeType == nil {
		ret += "TradeType:<nil>, "
	} else {
		ret += fmt.Sprintf("TradeType:%v, ", *o.TradeType)
	}

	if o.TransactionId == nil {
		ret += "TransactionId:<nil>"
	} else {
		ret += fmt.Sprintf("TransactionId:%v", *o.TransactionId)
	}

	return fmt.Sprintf("Transaction{%s}", ret)
}

// TransactionAmount
type TransactionAmount struct {
	Currency      *string `json:"currency,omitempty"`
	PayerCurrency *string `json:"payer_currency,omitempty"`
	PayerTotal    *int32  `json:"payer_total,omitempty"`
	Total         *int32  `json:"total,omitempty"`
}

func (o TransactionAmount) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Currency != nil {
		toSerialize["currency"] = o.Currency
	}

	if o.PayerCurrency != nil {
		toSerialize["payer_currency"] = o.PayerCurrency
	}

	if o.PayerTotal != nil {
		toSerialize["payer_total"] = o.PayerTotal
	}

	if o.Total != nil {
		toSerialize["total"] = o.Total
	}
	return json.Marshal(toSerialize)
}

func (o TransactionAmount) String() string {
	var ret string
	if o.Currency == nil {
		ret += "Currency:<nil>, "
	} else {
		ret += fmt.Sprintf("Currency:%v, ", *o.Currency)
	}

	if o.PayerCurrency == nil {
		ret += "PayerCurrency:<nil>, "
	} else {
		ret += fmt.Sprintf("PayerCurrency:%v, ", *o.PayerCurrency)
	}

	if o.PayerTotal == nil {
		ret += "PayerTotal:<nil>, "
	} else {
		ret += fmt.Sprintf("PayerTotal:%v, ", *o.PayerTotal)
	}

	if o.Total == nil {
		ret += "Total:<nil>"
	} else {
		ret += fmt.Sprintf("Total:%v", *o.Total)
	}

	return fmt.Sprintf("TransactionAmount{%s}", ret)
}

// TransactionPayer
type TransactionPayer struct {
	Openid *string `json:"openid,omitempty"`
}

func (o TransactionPayer) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Openid != nil {
		toSerialize["openid"] = o.Openid
	}
	return json.Marshal(toSerialize)
}

func (o TransactionPayer) String() string {
	var ret string
	if o.Openid == nil {
		ret += "Openid:<nil>"
	} else {
		ret += fmt.Sprintf("Openid:%v", *o.Openid)
	}

	return fmt.Sprintf("TransactionPayer{%s}", ret)
}

package model

import (
	"time"
)

type LightQueryModel struct {
	AWSAccessKeyId string `json:"AWSAccessKeyId"`
	AssociateTag   string `json:"AssociateTag"`
	Operation      string `json:"Operation"`
	ItemId         string `json:"ItemId"`
	Timestamp      time.Time `json:"Timestamp"`
	ResponseGroup  string `json:"ResponseGroup"`
}
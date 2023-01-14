package model

type TransferReq struct {
	To    string `json:"to" bson:"to" uri:"to" binding:"required"`
	Value int64  `json:"value" bson:"value" uri:"value" binding:"required"`
}

type TransferFromReq struct {
	To    string `json:"to" bson:"to" uri:"to" binding:"required"`
	Value int64  `json:"value" bson:"value" uri:"value" binding:"required"`
	Pk    string `json:"pk" bson:"pk" uri:"pk" binding:"required"`
}

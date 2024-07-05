package model

import "errors"

var (
	ErrRecordNotFound          = errors.New("record not found")
	ErrDuplicateBarcode        = errors.New("duplicate barcode")
	ErrDuplicatedProductOption = errors.New("duplicated product option")
)

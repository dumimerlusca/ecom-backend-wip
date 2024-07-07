package model

import "errors"

var (
	ErrRecordNotFound                      = errors.New("record not found")
	ErrDuplicatedProductOption             = errors.New("duplicatedd product option")
	ErrInvalidProductCategory              = errors.New("invalid product category")
	ErrProductCategoryNotFound             = errors.New("product category not found")
	ErrDuplicatedProductCategoryForProduct = errors.New("duplicated product category for same product")
	ErrParentProductCategoryNotFound       = errors.New("parent product category not found")
	ErrFileNotFound                        = errors.New("file not found")
)

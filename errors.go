package main

import "errors"

var ErrUnterminatedString = errors.New("unterminated string")
var ErrInvalidNumber = errors.New("invalid number literal")
var ErrUnexpectedChar = errors.New("unexpected character")

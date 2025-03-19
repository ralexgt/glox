package errors

import "errors"

var ErrUnterminatedString = errors.New("unterminated string")
var ErrInvalidNumber = errors.New("invalid number literal")
var ErrUnexpectedChar = errors.New("unexpected character")
var ErrUnterminatedComment = errors.New("unterminated multiline comment")

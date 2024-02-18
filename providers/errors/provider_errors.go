package providererrors

import (
	"errors"
)

var ErrMessageNotMatched = errors.New("no match found for input text")
var ErrIdNotMatched = errors.New("no other provider had a track that matched")
var ErrUnsupportedOperations = errors.New("unsupported operation")

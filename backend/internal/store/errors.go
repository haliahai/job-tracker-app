package store

import "errors"

// ErrNotFound is returned by any store method that couldn't find the row it
// was asked for. The API layer maps this to a 404 in one place
// (handleStoreError) instead of every handler checking sql.ErrNoRows itself.
var ErrNotFound = errors.New("resource not found")

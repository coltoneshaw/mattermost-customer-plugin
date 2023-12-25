package app

import "github.com/pkg/errors"

// ErrNotFound used when an entity is not found.
var ErrNotFound = errors.New("not found")

// ErrDuplicateEntry occurs when failing to insert because the entry already existed.
var ErrDuplicateEntry = errors.New("duplicate entry")

package testutils

import (
	"context"
	"testing"
)

func RunWithContext(t *testing.T, testFunc func(t *testing.T, ctx context.Context)) {
	testFunc(t, context.Background())
}

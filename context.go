package authproxy

import (
	"context"

	"github.com/pkg/errors"
)

type stateKey string

var ctxStateKey = stateKey(sessionName)

func getStateFromContext(ctx context.Context) (*sessionState, error) {
	state, ok := ctx.Value(ctxStateKey).(*sessionState)
	if !ok {
		return nil, errors.New("session state not in context")
	}

	if state == nil {
		return nil, errors.New("session state was nil")
	}

	if state.AuthRequestStates == nil {
		state.AuthRequestStates = make(map[string]authRequestState)
	}

	return state, nil
}

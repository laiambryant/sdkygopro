package endpoint

import (
	"context"
	"encoding/json"

	"github.com/laiambryant/sdkyGOprodeck/client"
)

type Endpoint[T any] struct {
	Client *client.Client
}

func New[T any](c *client.Client) *Endpoint[T] {
	return &Endpoint[T]{
		Client: c,
	}
}

func (e *Endpoint[T]) Fetch(ctx context.Context, path string) (T, error) {
	var item T
	data, err := e.Client.Get(ctx, path)
	if err != nil {
		return item, err
	}
	if err := json.Unmarshal(data, &item); err != nil {
		return item, &DecodeError{Resource: path, Err: err}
	}
	return item, nil
}

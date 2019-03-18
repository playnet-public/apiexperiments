package fakedb

import (
	"context"

	"github.com/playnet-public/demo/faction"
)

// Repository doing nothing
type Repository struct {}

// Create nothing
func (m *Repository) Create(_ context.Context, f faction.Incomplete) (faction.Complete, error) {
	return &complete{
		id:         "foo",
		Incomplete: f,
	}, nil
}

// Delete nothing
func (m *Repository) Delete(_ context.Context, _ faction.Identifier) error {
	return nil
}

// Update nothing
func (m *Repository) Update(_ context.Context, f faction.Complete) (faction.Complete, error) {
	return f, nil
}

// complete object containing both an allocated ID and the stored data
type complete struct {
	id string
	faction.Incomplete
}

// ID of the completed object
func (c complete) ID() string { return c.id }

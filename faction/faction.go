package faction

import (
	"context"
)

// Repository of factions providing data layer access
type Repository interface {
	Create(context.Context, Incomplete) (Complete, error)
	Delete(context.Context, Identifier) error
	Update(context.Context, Complete) (Complete, error)
}

// Manager provides access functions on factions
type Manager struct {
	Repository
}

// Identifier for checking or receiving id data of an object
type Identifier interface {
	ID() string
}

// Provider for faction data
type Provider interface {
	Data() *data
}

// Incomplete faction object pre-creation without an allocated ID
type Incomplete interface {
	Provider
}

// Complete faction object including meta and data payload
type Complete interface {
	Provider
	Identifier
}

// --- New File?

// data of a faction object
type data struct {
	Title, Description string
}

// Data returns the object itself
func (d *data) Data() *data { return d }

// SetDescription and return the resulting object
func (d *data) SetDescription(to string) *data {
	d.Description = to
	return d
}

// NewIncomplete faction object ready to get stored
func NewIncomplete(ctx context.Context, title, description string) Incomplete {
	return &data{
		Title:       title,
		Description: description,
	}
}

// --- This could go into the real db handling package as there should be 	--- //
// --- no way to create it manually, but the external db package has to		--- //

// complete object containing both an allocated ID and the stored data
// type complete struct {
// 	id string
// 	Incomplete
// }

// // ID of the completed object
// func (c complete) ID() string { return c.id }

// --- //

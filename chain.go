package main

import (
	"errors"
	"sync"
)

// Chain represents an ordered sequence of Beads that can be
// traversed, modified, and resolved in a thread-safe manner.
type Chain struct {
	mu    sync.RWMutex
	beads []*Bead
	name  string
}

// NewChain creates a new Chain with the given name.
func NewChain(name string) *Chain {
	return &Chain{
		name:  name,
		beads: make([]*Bead, 0),
	}
}

// Add appends a Bead to the end of the Chain.
func (c *Chain) Add(b *Bead) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.beads = append(c.beads, b)
}

// Insert places a Bead at the specified index, shifting
// subsequent beads to the right. Returns an error if the
// index is out of bounds.
func (c *Chain) Insert(index int, b *Bead) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index > len(c.beads) {
		return errors.New("chain: index out of bounds")
	}

	c.beads = append(c.beads, nil)
	copy(c.beads[index+1:], c.beads[index:])
	c.beads[index] = b
	return nil
}

// Remove deletes the Bead at the specified index.
// Returns an error if the index is out of bounds.
func (c *Chain) Remove(index int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if index < 0 || index >= len(c.beads) {
		return errors.New("chain: index out of bounds")
	}

	c.beads = append(c.beads[:index], c.beads[index+1:]...)
	return nil
}

// Get returns the Bead at the specified index.
// Returns an error if the index is out of bounds.
func (c *Chain) Get(index int) (*Bead, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if index < 0 || index >= len(c.beads) {
		return nil, errors.New("chain: index out of bounds")
	}

	return c.beads[index], nil
}

// Len returns the number of Beads in the Chain.
func (c *Chain) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.beads)
}

// Name returns the name of the Chain.
func (c *Chain) Name() string {
	return c.name
}

// All returns a snapshot copy of all Beads in the Chain.
// The returned slice is safe to iterate without holding the lock.
func (c *Chain) All() []*Bead {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := make([]*Bead, len(c.beads))
	copy(snapshot, c.beads)
	return snapshot
}

// Clear removes all Beads from the Chain.
func (c *Chain) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.beads = make([]*Bead, 0)
}

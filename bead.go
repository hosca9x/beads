package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Bead represents a single unit of data in the chain.
// Each bead contains a payload, metadata, and a cryptographic
// link to its predecessor, forming an immutable sequence.
type Bead struct {
	Index     uint64    `json:"index"`
	Timestamp time.Time `json:"timestamp"`
	Payload   []byte    `json:"payload"`
	PrevHash  string    `json:"prev_hash"`
	Hash      string    `json:"hash"`
	Nonce     uint64    `json:"nonce"`
}

// NewBead creates a new Bead linked to the previous bead's hash.
// The caller is responsible for computing the hash via Seal().
func NewBead(index uint64, payload []byte, prevHash string) *Bead {
	return &Bead{
		Index:     index,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
		PrevHash:  prevHash,
	}
}

// Seal computes and stores the SHA-256 hash of the bead's contents.
// It must be called after all fields are set and before the bead
// is appended to a strand.
func (b *Bead) Seal() {
	b.Hash = b.computeHash()
}

// computeHash derives a deterministic hash from the bead's fields.
func (b *Bead) computeHash() string {
	raw := fmt.Sprintf("%d|%s|%x|%s|%d",
		b.Index,
		b.Timestamp.Format(time.RFC3339Nano),
		b.Payload,
		b.PrevHash,
		b.Nonce,
	)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Verify checks that the stored hash matches a freshly computed hash,
// confirming the bead has not been tampered with.
func (b *Bead) Verify() bool {
	return b.Hash == b.computeHash()
}

// String returns a human-readable summary of the bead.
// Showing 12 hex chars instead of 8 gives a bit more collision safety
// when eyeballing hashes in logs.
func (b *Bead) String() string {
	return fmt.Sprintf("Bead{index=%d, hash=%.12s..., prev=%.12s..., ts=%s}",
		b.Index,
		b.Hash,
		b.PrevHash,
		b.Timestamp.Format(time.RFC3339),
	)
}

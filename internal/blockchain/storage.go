package blockchain

import (
	"encoding/json"
	"errors"

	bolt "go.etcd.io/bbolt"
)

const (
	blocksBucket = "blocks"
	lastHashKey  = "last_hash"
)

// Storage wraps the persistent blockchain database.
type Storage struct {
	db *bolt.DB
}

// OpenStorage opens or creates the blockchain database.
func OpenStorage(path string) (*Storage, error) {
	db, err := bolt.Open(path, 0600, nil)

	if err != nil {
		return nil, err
	}

	return &Storage{
		db: db,
	}, nil
}

// Close closes the underlying database.
func (s *Storage) Close() error {
	if s == nil || s.db == nil {
		return nil
	}

	return s.db.Close()
}

// SaveBlock stores a block keyed by its hash
// and updates the latest-block pointer.
func (s *Storage) SaveBlock(block *Block) error {
	if block == nil {
		return errors.New("block cannot be nil")
	}

	data, err := json.Marshal(block)

	if err != nil {
		return err
	}

	return s.db.Update(
		func(tx *bolt.Tx) error {
			bucket, err := tx.CreateBucketIfNotExists(
				[]byte(blocksBucket),
			)

			if err != nil {
				return err
			}

			if err := bucket.Put(
				block.Hash,
				data,
			); err != nil {
				return err
			}

			return bucket.Put(
				[]byte(lastHashKey),
				block.Hash,
			)
		},
	)
}

// LoadBlock retrieves a block by hash.
func (s *Storage) LoadBlock(
	hash []byte,
) (*Block, error) {
	var block Block

	err := s.db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket(
				[]byte(blocksBucket),
			)

			if bucket == nil {
				return errors.New(
					"blocks bucket does not exist",
				)
			}

			data := bucket.Get(hash)

			if data == nil {
				return errors.New(
					"block not found",
				)
			}

			return json.Unmarshal(
				data,
				&block,
			)
		},
	)

	if err != nil {
		return nil, err
	}

	return &block, nil
}

// LastHash returns the hash of the latest persisted block.
func (s *Storage) LastHash() ([]byte, error) {
	var result []byte

	err := s.db.View(
		func(tx *bolt.Tx) error {
			bucket := tx.Bucket(
				[]byte(blocksBucket),
			)

			if bucket == nil {
				return errors.New(
					"blocks bucket does not exist",
				)
			}

			value := bucket.Get(
				[]byte(lastHashKey),
			)

			if value == nil {
				return errors.New(
					"last block hash not found",
				)
			}

			result = append(
				[]byte(nil),
				value...,
			)

			return nil
		},
	)

	return result, err
}

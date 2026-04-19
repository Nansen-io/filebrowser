package protection

import (
	storm "github.com/asdine/storm/v3"
)

// Record holds ChainFS protection metadata for a file path.
type Record struct {
	Path     string `storm:"id"`
	FileGuid string
	Expiry   int64
}

// Storage persists protection records in BoltDB.
type Storage struct {
	db *storm.DB
}

func NewStorage(db *storm.DB) *Storage {
	return &Storage{db: db}
}

// Save creates or overwrites the protection record for a file path.
func (s *Storage) Save(path, fileGuid string, expiry int64) error {
	return s.db.Save(&Record{Path: path, FileGuid: fileGuid, Expiry: expiry})
}

// Get returns the protection record for path, or nil if not found.
func (s *Storage) Get(path string) (*Record, error) {
	var r Record
	if err := s.db.One("Path", path, &r); err != nil {
		if err == storm.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}

// Delete removes the protection record for path.
func (s *Storage) Delete(path string) error {
	err := s.db.DeleteStruct(&Record{Path: path})
	if err == storm.ErrNotFound {
		return nil
	}
	return err
}

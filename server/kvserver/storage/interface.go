/**
 * Storage is a package aim to provide interface for storage engine
 * Engine may include Bitcask or LSM-Tree
 */

package storage

// Engine is an interface for storage backend like LSM-tree or Bitcask
//
type Engine interface {
	// Get the data in
	Get([]byte) ([]byte, error)
	// Scan the data of engine, return bytes array
	Scan(cursor ScanCursor)	([][]byte, error)
	// Put the data into the engine
	Put([]byte, []byte) error
	// delete the data in the engine
	Del([]byte) error
}
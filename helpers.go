package celeritas

import "os"

func (c *Celeritas) CreateDirIfNotExists(path string) error {
	// permitions
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			return err
		}
	}
	return nil
}

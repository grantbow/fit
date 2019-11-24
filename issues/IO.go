package issues

import (
	"fmt"
	"os"
)

//var dops = Directory(os.PathSeparator)
//var sops = string(os.PathSeparator)

func (i *Issue) Read(p []byte) (int, error) {
	if i.DescriptionFileName == "" {
		return 0, ErrNoDescription
	}
	if i.descFile == nil {
		dir := i.Direr()
		fp, err := os.OpenFile(string(dir)+sops+i.DescriptionFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		i.descFile = fp
		if err != nil {
			fmt.Fprintf(os.Stderr, "err: %s", err.Error())
			return 0, ErrNoDescription
		}
	}

	return i.descFile.Read(p)
}

func (i *Issue) Write(data []byte) (n int, err error) {
	if i.DescriptionFileName == "" {
		return 0, ErrNoDescription
	}
	if i.descFile == nil {
		dir := i.Direr()
		os.MkdirAll(string(dir), 0755)
		fp, err := os.OpenFile(string(dir)+sops+i.DescriptionFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to issue: %s", err.Error())
			return 0, err
		}
		i.descFile = fp
	}
	return i.descFile.Write(data)
}

// WriteAt makes a directory, writes a byte string to the Description using an offset.
// It returns the number of bytes written and an error.
func (i *Issue) WriteAt(data []byte, off int64) (n int, err error) {
	if i.DescriptionFileName == "" {
		return 0, ErrNoDescription
	}
	if i.descFile == nil {
		dir := i.Direr()
		os.MkdirAll(string(dir), 0755)
		fp, err := os.OpenFile(i.DescriptionFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to issue: %s", err.Error())
			return 0, err
		}
		i.descFile = fp
	}
	return i.descFile.WriteAt(data, off)
}

// Close returns an error if there is an error closing the descFile of an issue.
func (i Issue) Close() error {
	if i.descFile != nil {
		err := i.descFile.Close()
		i.descFile = nil
		return err
	}
	return nil
}

// Remove deletes the directory and files of an issue.
func (i *Issue) Remove() error {
	dir := i.Direr()
	if dir != "" {
		return os.RemoveAll(string(dir))
	}
	return ErrNotFound
}

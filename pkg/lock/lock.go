package lock

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

//FileLock locks file by name
type FileLock struct {
	path     string
	lockPath string
	f        *os.File
}

func New(path string) *FileLock {
	return &FileLock{
		path:     path,
		lockPath: fmt.Sprintf("%s.lock", path),
	}
}

func (l *FileLock) Lock() error {
	f, err := os.OpenFile(l.lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err == nil {
		l.f = f
		return nil
	}
	if os.IsExist(err) {
		st, err := os.Stat(l.lockPath)
		if err != nil {
			return err
		}
		if st.ModTime().Before(time.Now().Add(-30 * time.Second)) {
			os.Remove(l.lockPath)
			f, err := os.OpenFile(l.lockPath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
			if err == nil {
				l.f = f
				return nil
			}
		}
	}
	return err
}

// Lock
func (l *FileLock) LockAndRead() (string, error) {
	var rs string
	if err := l.Lock(); err != nil {
		return rs, err
	}

	if _, err := os.Stat(l.path); os.IsExist(err) {
		data, err := ioutil.ReadFile(l.path)
		if err != nil {
			return rs, err
		}
		rs = string(data)
		os.Rename(l.path, fmt.Sprintf("%s.old", l.path))
	}
	return rs, nil
}

func (l *FileLock) Write(data []byte) error {
	_, err := l.f.Write(data)
	return err
}

// Unlock
func (l *FileLock) Unlock() error {
	if l.f != nil {
		return os.Remove(l.lockPath)
	}
	return nil
}

// Unlock
func (l *FileLock) Flush() error {
	if err := l.f.Close(); err != nil {
		return err
	}
	if err := os.Rename(l.lockPath, l.path); err != nil {
		return err
	}
	l.f = nil
	return nil
}

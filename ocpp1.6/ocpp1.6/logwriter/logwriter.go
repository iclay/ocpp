package logwriter

import (
	"math"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
	"time"
)

// HourlySplit split outer file hourly
type HourlySplit struct {
	Dir        string // path of log file directory
	FileFormat string // log_2006-01-02T15
	// MaxFileNumber is max file number, auto delete old file if limit exceed.
	// treat as not limit if set to 0.
	MaxFileNumber int64
	// max disk ussage, auto delete old file if limit exceed.
	// It is a soft limit - we do not change every time the Write() be called,
	// we check only when file spliting.
	// 0 means not limit.
	MaxDiskUsage int64

	curFileName string
	prevUpdate  time.Time
	file        *os.File
	mu          sync.Mutex
}

// diskUssage travel the directory of log, then return disk usage of all log files
// and the count of log files. If a file in directory do not fit the format specify
// by a.FileFormat, the count and disk usage of it will be ignore.
func (a *HourlySplit) diskUssage() (disk int64, fnum int64) {
	info, err := os.Lstat(a.Dir)
	if err != nil {
		return
	}
	if !info.IsDir() {
		return
	}
	dir, err := os.Open(a.Dir)
	if err != nil {
		return
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		_, err := time.Parse(a.FileFormat, file.Name())
		if err != nil {
			continue
		}
		disk += file.Size()
		fnum++
	}
	return
}

// byModTime wrap the os.FileInfo slice type for sorting by ModTime()
type byModTime []os.FileInfo

func (a byModTime) Len() int           { return len(a) }
func (a byModTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byModTime) Less(i, j int) bool { return a[i].ModTime().Before(a[j].ModTime()) }

// keepLimit check if MaxDiskUsage and MaxFileNumber be satisfied, if not,
// delete the oldest file until the check return true.
// If there are just one log file in the directory, it truncate the file
// to MaxDiskUsage/2.
func (a *HourlySplit) keepLimit() error {
	if a.MaxDiskUsage <= 0 {
		a.MaxDiskUsage = math.MaxInt64
	}
	if a.MaxFileNumber <= 0 {
		a.MaxFileNumber = math.MaxInt64
	}
	disk, fnum := a.diskUssage()
	dir, err := os.Open(a.Dir)
	if err != nil {
		return err
	}
	defer dir.Close()
	if disk >= a.MaxDiskUsage || fnum >= a.MaxFileNumber {
		var pathDelete []string
		files, err := dir.Readdir(-1)
		if err != nil {
			return err
		}
		sort.Sort(byModTime(files))
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			_, err := time.Parse(a.FileFormat, file.Name())
			if err != nil {
				continue
			}
			disk -= file.Size()
			fnum--
			pathDelete = append(pathDelete, file.Name())
			if disk < a.MaxDiskUsage && fnum+1 < a.MaxFileNumber {
				break
			}
		}
		for _, path := range pathDelete {
			if a.curFileName == path {
				os.Truncate(filepath.Join(a.Dir, path), 1024*1024)
				continue
			}
			os.RemoveAll(filepath.Join(a.Dir, path))
		}
	}
	return nil
}

// urgentLimit delete all files not currently write,
// truncat all other files to 1000.
// trunkcat current writing file to 1m.
func (a *HourlySplit) urgentLimit() error {
	dir, err := os.Open(a.Dir)
	if err != nil {
		return err
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		_, err := time.Parse(a.FileFormat, file.Name())
		if err != nil {
			os.Truncate(filepath.Join(a.Dir, file.Name()), 1000)
			continue
		}
		if a.curFileName == file.Name() {
			os.Truncate(filepath.Join(a.Dir, file.Name()), 1000000)
			continue
		}
		os.RemoveAll(filepath.Join(a.Dir, file.Name()))
	}
	return nil
}

// update check if need close the current file and create a new one.
// limit checking will apply here.
func (a *HourlySplit) update() (err error) {
	cur := time.Now()
	prev := a.prevUpdate
	if cur.Year() != prev.Year() || cur.YearDay() != prev.YearDay() || cur.Hour() != prev.Hour() {
		a.prevUpdate = cur
		if a.file != nil {
			a.file.Close()
		}
		newFileName := cur.Format(a.FileFormat)
		a.curFileName = newFileName
		newFilePath := filepath.Join(a.Dir, newFileName)
		os.MkdirAll(a.Dir, 0755)
		a.keepLimit()
		a.file, err = os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		return err
	}
	return nil
}

// Write implement writer
func (a *HourlySplit) Write(b []byte) (n int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	err = a.update()
	if err != nil {
		return 0, err
	}
	// set write timeout to avoid block when disk full.
	a.file.SetWriteDeadline(time.Now().Add(time.Second))
	n, err = a.file.Write(b)
	if err != nil {
		stat := syscall.Statfs_t{}
		err = syscall.Statfs(a.Dir, &stat)
		if err != nil {
			return 0, err
		}
		av := stat.Bavail * uint64(stat.Bsize)
		if av < 1024*1024 {
			a.urgentLimit()
		}
	}
	return
}

// Close close all file discriptor
func (a *HourlySplit) Close() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.prevUpdate = time.Time{}
	f := a.file
	a.file = nil
	if f != nil {
		return f.Close()
	}
	return nil
}

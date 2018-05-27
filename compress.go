package filehook

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

// compress compress files with gizp
func compress(compressPath, sourceExt, compressExt string, excludeFiles []string) error {
	cFileName := time.Now().Format("2006-01-02_15-04-05") + compressExt
	cFileWholeName := filepath.Join(compressPath, cFileName)

	compressDir, err := os.Open(compressPath)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer compressDir.Close()

	files, err := compressDir.Readdir(0)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	cf, err := os.Create(cFileWholeName)
	if err != nil {
		return fmt.Errorf(" %v", err)
	}
	defer cf.Close()

	gzipW := gzip.NewWriter(cf)
	defer gzipW.Close()

	tarW := tar.NewWriter(gzipW)
	defer tarW.Close()

	var compressedFileList []string
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		if !strings.Contains(fi.Name(), sourceExt) {
			continue
		}
		if containInList(fi.Name(), excludeFiles) {
			continue
		}

		log.Debugf("compress file %v to %v", fi.Name(), cFileWholeName)

		f, err := os.Open(filepath.Join(compressDir.Name(), fi.Name()))
		if err != nil {
			log.Debugf("compress file %v to %v error: %v", fi.Name(), cFileWholeName, err)
			continue
		}

		h := new(tar.Header)
		h.Name = fi.Name()
		h.Size = fi.Size()
		h.Mode = int64(fi.Mode())
		h.ModTime = fi.ModTime()

		err = tarW.WriteHeader(h)
		if err != nil {
			log.Debugf("compress file %v to %v error: %v", fi.Name(), cFileWholeName, err)
			f.Close()
			continue
		}

		_, err = io.Copy(tarW, f)
		if err != nil {
			log.Debugf("compress file %v to %v error: %v", fi.Name(), cFileWholeName, err)
			f.Close()
			continue
		}

		compressedFileList = append(compressedFileList, fi.Name())
		f.Close()
	}

	if len(compressedFileList) < 1 {
		cf.Close()
		cleanFiles(compressDir.Name(), []string{cFileName})
	} else {
		cleanFiles(compressDir.Name(), compressedFileList)
	}

	log.Debugf("compress file %v success", cFileWholeName)
	return nil
}

func containInList(s string, list []string) bool {
	if list == nil {
		return false
	}
	for _, ss := range list {
		if strings.Contains(ss, s) {
			return true
		}
	}
	return false
}

func cleanFiles(dirName string, fileNames []string) {
	if len(fileNames) < 1 {
		return
	}

	for _, file := range fileNames {
		hp := filepath.Join(dirName, file)
		err := os.Remove(hp)
		if err != nil {
			log.Debugf("remove file %v error: %v", hp, err)
			continue
		}
		log.Debugf("remove file %v success", hp)
	}
}

// autoCompress auto compress and package all log files
func (h *FileHook) autoCompress() {
	if !h.option.Compress.Enable {
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(h.option.Compress.Interval))
	log.Debugf("compression start per %v", time.Second*time.Duration(h.option.Compress.Interval))
	for {
		select {
		case <-ticker.C:
			h.File.Name()
			log.Debugf("start to compress log files in direction %v", h.option.Path)
			err := compress(h.option.Path, h.option.File.Ext, h.option.Compress.Ext, []string{h.File.File.Name()})
			if err != nil {
				log.Errorf("compress log files error: %v", err)
				continue
			}
		}
	}
}

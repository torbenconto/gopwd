package backup

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	iou "io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func ArchiveGopwdVault(vaultPath, toDir string) error {
	archiveName := strings.Split(vaultPath, "/")[len(strings.Split(vaultPath, "/"))-1] + "_" + time.Now().Format("2006-01-02_15-04-05") + ".tar.gz"
	archivePath := filepath.Join(toDir, archiveName)

	fmt.Println("Creating archive:", archivePath)

	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gz := gzip.NewWriter(file)
	defer gz.Close()

	tarw := tar.NewWriter(gz)
	defer tarw.Close()

	vaultPath = filepath.Clean(vaultPath) // Ensure vaultPath is in clean form

	return filepath.Walk(vaultPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Ensure the header name is the relative path from vaultPath
		relPath, err := filepath.Rel(vaultPath, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}

		header.Name = relPath // Use relative path for header name

		if err := tarw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close() // Close the file immediately after copying its content

			_, err = iou.Copy(tarw, file)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func RestoreGopwdVault(archivePath, vaultPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == iou.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(vaultPath, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return err
			}
		case tar.TypeReg:
			file, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer file.Close()

			if _, err := iou.Copy(file, tarReader); err != nil {
				return err
			}
		default:
			log.Printf("Unsupported type flag: %c in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

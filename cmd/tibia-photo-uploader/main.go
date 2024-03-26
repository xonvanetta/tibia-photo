package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/koding/multiconfig"
	"github.com/xonvanetta/tibia-photo/pkg/tibia"
	"github.com/xonvanetta/tibia-photo/pkg/upload/sftp"
)

type Config struct {
	FilesLocation         string
	WhitelistCharacters   []string
	Upload, Remove, Debug bool

	SFTP sftp.Config

	Delete bool
}

func main() {
	config := &Config{}
	multiconfig.MustLoad(config)

	level := slog.LevelInfo
	if config.Debug {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	slog.SetDefault(logger)

	uploader, err := sftp.New(config.SFTP)
	if err != nil {
		slog.Error("failed to create sftp client: %s", err)
		return
	}

	remove := func(path string) error {
		slog.Debug("remove", "path", path)
		if !config.Remove {
			return nil
		}
		err := os.Remove(path)
		if err != nil {
			return fmt.Errorf("failed to remove file on path %s: %w", path, err)
		}
		return nil
	}

	err = filepath.WalkDir(config.FilesLocation, func(path string, entry fs.DirEntry, err error) error {
		if !strings.HasSuffix(path, ".png") {
			return nil
		}
		if !isWhiteListed(path, config.WhitelistCharacters) {
			return remove(path)
			//return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		photo, err := tibia.ParsePhotoFromFile(file)
		if err != nil {
			return fmt.Errorf("failed to create new photo: %w", err)
		}

		if !photo.Keep() {
			return remove(file.Name())
		}

		_, err = uploader.Stat(photo.FullPath())
		if err == nil {
			return remove(file.Name())
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to stats photo path: %w", err)
		}

		path = fmt.Sprintf("%s/%s", config.SFTP.RootPath, photo.FullPath())
		slog.Debug("upload", "path", path)

		if !config.Upload {
			return nil
		}

		err = uploader.Upload(photo.FullPath(), photo.File(), photo.Time(), photo.Time())
		if err != nil {
			return fmt.Errorf("failed to upload photo: %w", err)
		}

		return nil
	})
	if err != nil {
		slog.Error("failed to walk filepath: %s", err)
	}
}

func isWhiteListed(path string, list []string) bool {
	for _, i := range list {
		if strings.Contains(strings.ToLower(path), strings.ToLower(i)) {
			return true
		}
	}
	return false
}

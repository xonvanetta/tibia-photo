package sftp

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/sftp"
	"github.com/xonvanetta/tibia-photo/pkg/upload"
	"golang.org/x/crypto/ssh"
)

type Config struct {
	Host     string
	RootPath string

	Username, Password string
}

func (c Config) path(path string) string {
	return fmt.Sprintf("%s/%s", c.RootPath, path)
}

type uploader struct {
	client *sftp.Client

	config Config
}

func New(config Config) (upload.Uploader, error) {
	sshConfig := &ssh.ClientConfig{
		User:            config.Username,
		Timeout:         time.Second * 10,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if config.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.Password(config.Password))
	}

	sshClient, err := ssh.Dial("tcp", config.Host, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to ssh dial: %w", err)
	}

	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create sftp client: %w", err)
	}
	return &uploader{
		client: sftpClient,
		config: config,
	}, nil

}

func (s *uploader) Stat(path string) (os.FileInfo, error) {
	return s.client.Stat(s.config.path(path))
}

//func (s *uploader) Remove(path string) error {
//	return nil
//}

func (s *uploader) Upload(path string, body io.Reader, atime, mtime time.Time) error {
	path = s.config.path(path)

	err := s.client.MkdirAll(filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}
	file, err := s.client.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	_, err = io.Copy(file, body)
	if err != nil {
		return fmt.Errorf("failed to copy photo: %w", err)
	}

	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	err = s.client.Chtimes(path, atime, mtime)
	if err != nil {
		return fmt.Errorf("failed to change time on sftp file: %w", err)
	}
	return nil
}

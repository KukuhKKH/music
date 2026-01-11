package storage

import (
	"fmt"
	"io"
	"log"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SftpStorage struct {
	Host      string
	Port      int
	User      string
	Password  string
	BaseDir   string
	PublicUrl string
}

func NewSftpStorage(host string, port int, user, password, baseDir, publicUrl string) *SftpStorage {
	return &SftpStorage{
		Host:      host,
		Port:      port,
		User:      user,
		Password:  password,
		BaseDir:   baseDir,
		PublicUrl: publicUrl,
	}
}

func (s *SftpStorage) connect() (*ssh.Client, *sftp.Client, error) {
	config := &ssh.ClientConfig{
		User: s.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("Gagal Dial SSH ke %s: %v", addr, err)
		return nil, nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		// explicitly ignore error from conn.Close()
		_ = conn.Close()
		log.Printf("Gagal Handshake SFTP: %v", err)
		return nil, nil, err
	}

	return conn, client, nil
}

func (s *SftpStorage) Upload(filename string, file io.Reader) (string, error) {
	sshConn, client, err := s.connect()
	if err != nil {
		return "", err
	}

	defer func() { _ = sshConn.Close() }()
	defer func() { _ = client.Close() }()

	// build remote POSIX-style full path
	fullPath := filename
	if s.BaseDir != "" {
		fullPath = path.Join(s.BaseDir, filename)
	}

	// ensure parent directory exists (handles nested paths in filename)
	dir := path.Dir(fullPath)
	if dir != "" && dir != "." {
		if err := client.MkdirAll(dir); err != nil {
			log.Printf("Gagal create dir %s: %v", dir, err)
			return "", err
		}
	}

	dstFile, err := client.Create(fullPath)
	if err != nil {
		log.Printf("Gagal create file %s: %v", fullPath, err)
		return "", err
	}

	defer func() { _ = dstFile.Close() }()

	_, err = io.Copy(dstFile, file)
	if err != nil {
		log.Printf("Gagal upload data: %v", err)
		return "", err
	}

	return filename, nil
}

func (s *SftpStorage) Delete(filename string) error {
	sshConn, client, err := s.connect()
	if err != nil {
		return err
	}
	defer func() { _ = sshConn.Close() }()
	defer func() { _ = client.Close() }()

	fullPath := filename
	if s.BaseDir != "" {
		fullPath = path.Join(s.BaseDir, filename)
	}

	return client.Remove(fullPath)
}

func (s *SftpStorage) GetURL(filename string) string {
	fullPath := path.Join(s.BaseDir, filename)
	relativePath := strings.TrimPrefix(fullPath, "upload/")
	return fmt.Sprintf("%s/%s", s.PublicUrl, relativePath)
}

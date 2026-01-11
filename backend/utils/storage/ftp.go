package storage

import (
	"context"
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
		Timeout:         10 * time.Second, // handshake timeout
	}

	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("[sftp] dial ssh %s err=%v", addr, err)
		return nil, nil, err
	}

	client, err := sftp.NewClient(conn)
	if err != nil {
		_ = conn.Close()
		log.Printf("[sftp] new client err=%v", err)
		return nil, nil, err
	}

	return conn, client, nil
}

// Upload sekarang menerima context untuk timeout/cancel
func (s *SftpStorage) Upload(ctx context.Context, filename string, r io.Reader) (string, error) {
	sshConn, client, err := s.connect()
	if err != nil {
		return "", err
	}
	defer func() { _ = client.Close() }()
	defer func() { _ = sshConn.Close() }()

	fullPath := filename
	if s.BaseDir != "" {
		fullPath = path.Join(s.BaseDir, filename)
	}

	dir := path.Dir(fullPath)
	if dir != "" && dir != "." {
		if err := client.MkdirAll(dir); err != nil {
			log.Printf("[sftp] mkdir %s err=%v", dir, err)
			return "", err
		}
	}

	dstFile, err := client.Create(fullPath)
	if err != nil {
		log.Printf("[sftp] create %s err=%v", fullPath, err)
		return "", err
	}
	defer func() { _ = dstFile.Close() }()

	type copyResult struct {
		n   int64
		err error
	}
	ch := make(chan copyResult, 1)

	go func() {
		n, e := io.Copy(dstFile, r)
		ch <- copyResult{n: n, err: e}
	}()

	select {
	case <-ctx.Done():
		_ = dstFile.Close()
		_ = client.Close()
		_ = sshConn.Close()
		return "", fmt.Errorf("sftp upload canceled/timeout: %w", ctx.Err())

	case res := <-ch:
		if res.err != nil {
			log.Printf("[sftp] copy err=%v", res.err)
			return "", res.err
		}
		log.Printf("[sftp] upload ok path=%s bytes=%d", fullPath, res.n)
		return filename, nil
	}
}

func (s *SftpStorage) Delete(filename string) error {
	sshConn, client, err := s.connect()
	if err != nil {
		return err
	}
	defer func() { _ = client.Close() }()
	defer func() { _ = sshConn.Close() }()

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

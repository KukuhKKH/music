package storage

import (
	"fmt"
	"io"
	"time"

	"github.com/jlaffaye/ftp"
)

type FtpStorage struct {
	Host      string
	Port      int
	User      string
	Password  string
	BaseDir   string
	PublicUrl string
}

func NewFtpStorage(host string, port int, user, password, baseDir, publicUrl string) *FtpStorage {
	return &FtpStorage{
		Host:      host,
		Port:      port,
		User:      user,
		Password:  password,
		BaseDir:   baseDir,
		PublicUrl: publicUrl,
	}
}

func (s *FtpStorage) connect() (*ftp.ServerConn, error) {
	c, err := ftp.Dial(fmt.Sprintf("%s:%d", s.Host, s.Port), ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	err = c.Login(s.User, s.Password)
	if err != nil {
		_ = c.Quit()
		return nil, err
	}

	return c, nil
}

func (s *FtpStorage) Upload(filename string, file io.Reader) (string, error) {
	c, err := s.connect()
	if err != nil {
		return "", err
	}

	defer c.Quit()

	if s.BaseDir != "" {
		_ = c.ChangeDir(s.BaseDir)
	}

	err = c.Stor(filename, file)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *FtpStorage) Delete(filename string) error {
	c, err := s.connect()
	if err != nil {
		return err
	}

	defer c.Quit()

	if s.BaseDir != "" {
		_ = c.ChangeDir(s.BaseDir)
	}

	return c.Delete(filename)
}

func (s *FtpStorage) GetURL(filename string) string {
	return fmt.Sprintf("%s/%s/%s", s.PublicUrl, s.BaseDir, filename)
}

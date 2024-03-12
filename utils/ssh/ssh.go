package ssh

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// SSHClient 定义了一个 SSH 客户端的接口  11
type SSHClient interface {
	Connect() (*ssh.Client, *ssh.Session, error)
	RunCommand(command string) ([]byte, error)
	SftpConnect() (*sftp.Client, error)
	Copy(localFilePath, remotePath string) error
}

// SSH 结构体实现了 SSHClient 接口
type SSH struct {
	IP       string
	Username string
	Password string
}

const DefaultSSHPort = "22"

// 连接方法
func (s *SSH) connect(host string) (*ssh.Client, error) {
	ip := host
	user := s.Username
	password := s.Password
	DefaultTimeout := time.Duration(15) * time.Second
	clientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         DefaultTimeout,
	}
	return ssh.Dial("tcp", JoinHostPort(ip, DefaultSSHPort), clientConfig)
}

// 可调用的连接
func (s *SSH) Connect() (*ssh.Client, *ssh.Session, error) {
	client, err := s.connect(s.IP)
	if err != nil {
		return nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		_ = client.Close()
		return nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     //disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		_ = session.Close()
		_ = client.Close()
		return nil, nil, err
	}

	return client, session, nil
}

// 执行命令
func (s *SSH) RunCommand(command string) ([]byte, error) {
	client, session, err := s.Connect()
	if err != nil {
		return nil, fmt.Errorf("[ssh][%s] failed to create ssh session: %s", s.IP, err)
	}
	defer client.Close()
	defer session.Close()
	var stdoutContent, stderrContent bytes.Buffer
	session.Stdout = &stdoutContent
	session.Stderr = &stderrContent
	if err := session.Run(command); err != nil {
		return stdoutContent.Bytes(), fmt.Errorf("[ssh][%s]failed to run command[%s]: %s", s.IP, command, stderrContent.String())
	}
	return stdoutContent.Bytes(), nil
}

func NewSSHClient(ssh *SSH) SSHClient {
	return &SSH{
		IP:       ssh.IP,
		Username: ssh.Username,
		Password: ssh.Password,
	}
}

type Client struct {
	SSHClient  *ssh.Client
	SftpClient *sftp.Client
}

var sshClientMap = map[string]Client{}

var getSSHClientLock = sync.Mutex{}

func (s *SSH) SftpConnect() (*sftp.Client, error) {
	getSSHClientLock.Lock()
	defer getSSHClientLock.Unlock()

	if ret, ok := sshClientMap[s.IP]; ok {
		return ret.SftpClient, nil
	}

	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
		err        error
	)

	sshClient, err = s.connect(s.IP)
	if err != nil {
		return nil, err
	}

	sftpClient, err = sftp.NewClient(sshClient)

	sshClientMap[s.IP] = Client{
		SSHClient:  sshClient,
		SftpClient: sftpClient,
	}

	return sftpClient, err
}

func (s *SSH) Copy(localFilePath, remotePath string) error {

	sftpClient, err := s.SftpConnect()
	if err != nil {
		return fmt.Errorf("failed to new sftp client of host(%s): %s", s.IP, err)
	}
	defer sftpClient.Close()
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	var remoteFileName = path.Base(localFilePath)
	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}

	return nil
}

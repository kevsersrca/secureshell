package Shell

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os/user"
	"path/filepath"
	"sync"
	"time"
	"golang.org/x/crypto/ssh"

)

type Client struct {
	SSHClient *ssh.Client
}

const DefaultTimeout = 30 * time.Second

const sudoPassword = "sudo_password"

var HostKeyCallback = ssh.InsecureIgnoreHostKey()

type sudoWriter struct {
	b     bytes.Buffer
	pw    string
	stdin io.Writer
	m     sync.Mutex
}

func (w *sudoWriter) Write(p []byte) (int, error) {
	if string(p) == sudoPassword {
		w.stdin.Write([]byte(w.pw + "\n"))
		w.pw = ""
		return len(p), nil
	}

	w.m.Lock()
	defer w.m.Unlock()

	return w.b.Write(p)
}

func ConnectWithPassword(host, username, pass string, timeout ...time.Duration) (*Client, error) {
	authMethod := ssh.Password(pass)

	return connect(username, host, authMethod, timeout[0])
}

func ConnectWithKeyFile(host, username, privKeyPath string, timeout ...time.Duration) (*Client, error) {
	if privKeyPath == "" {
		currentUser, err := user.Current()
		if err == nil {
			privKeyPath = filepath.Join(currentUser.HomeDir, ".ssh", "id_rsa")
		}
	}

	privKey, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey([]byte(privKey))
	if err != nil {
		return nil, err
	}

	authMethod := ssh.PublicKeys(signer)

	return connect(username, host, authMethod, timeout[0])
}


func connect(username, host string, authMethod ssh.AuthMethod, timeout time.Duration) (*Client, error) {
	if username == "" {
		user, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("Username couldn't get current user: %v", err)
		}

		username = user.Username
	}

	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{authMethod},
		HostKeyCallback: HostKeyCallback,
	}

	_, _, err := net.SplitHostPort(host)
	if err != nil {
		host = net.JoinHostPort(host, "22")
	}

	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		return nil, err
	}

	sshConn, chans, reqs, err := ssh.NewClientConn(conn, host, config)
	if err != nil {
		return nil, err
	}
	client := ssh.NewClient(sshConn, chans, reqs)

	c := &Client{SSHClient: client}
	return c, nil
}

func (c *Client) Exec(cmd string) ([]byte, error) {
	session, err := c.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.CombinedOutput(cmd)
}

func (c *Client) ExecRoot(cmd, passwd string) ([]byte, error) {
	session, err := c.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	cmd = "sudo -p " + sudoPassword + " -S " + cmd

	w := &sudoWriter{
		pw: passwd,
	}
	w.stdin, err = session.StdinPipe()
	if err != nil {
		return nil, err
	}

	session.Stdout = w
	session.Stderr = w

	err = session.Run(cmd)

	return w.b.Bytes(), err
}

func (c *Client) Close() error {
	return c.SSHClient.Close()
}

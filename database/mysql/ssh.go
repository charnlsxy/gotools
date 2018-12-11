package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"os"
)

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func SshConfig(user string) *ssh.ClientConfig {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(os.Getenv("HOME") + "/.ssh/id_rsa"),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return sshConfig
}

type ViaSSHDialer struct {
	client *ssh.Client
}

func (self *ViaSSHDialer) Dial(addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

func OpenSsh() {
	sshcon, err := ssh.Dial("tcp", "..", SshConfig(".."))
	if err != nil {
		fmt.Println(err)
		return
	}
	//session,err := sshcon.NewSession()
	//if tools.ErrNotNil(err){
	//	return
	//}
	//session.Stdout = os.Stdout
	//session.Run("pwd")
	mysql.RegisterDial("tcp", (&ViaSSHDialer{sshcon}).Dial)

}

func NewSshTunnel(host, user, passwd string) (*ssh.Client, error) {
	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return ssh.Dial("tcp", host, cfg)
}

//通过ssh获取一条mysql连接
func NewMysqlDbInSSH(host, user string, cfg *MysqlConfig) (*sql.DB, error) {
	sshcon, err := ssh.Dial("tcp", host, SshConfig(user))
	if err != nil {
		return nil, err
	}
	mysql.RegisterDial("tcp", (&ViaSSHDialer{sshcon}).Dial)

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.PassWord, cfg.Host, cfg.Db))
	if err != nil {
		return nil, err
	}

	return db, nil
}

package main

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
    "log"
    "net"
    "os"
    "time"
)

const (
    CERT_PASSWORD = 1
    CERT_PUBLIC_KEY_FILE = 2
    DEFAULT_TIMEOUT = 3
)

type SSH struct{
    Host string
    User string
    Credential string
    Port int
    session *ssh.Session
    client *ssh.Client
}


// SSH Function //

func (sshClient *SSH) readPublicKeyFile(file string) ssh.AuthMethod {
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

func (sshClient *SSH) Connect(mode int) {
    var ssh_config *ssh.ClientConfig
    var auth  []ssh.AuthMethod
    if mode == CERT_PASSWORD {
        auth = []ssh.AuthMethod{ssh.Password(sshClient.Credential)}
    } else if mode == CERT_PUBLIC_KEY_FILE {
        auth = []ssh.AuthMethod{sshClient.readPublicKeyFile(sshClient.Credential)}
    } else {
        log.Println("SSH doesn't support mode", mode)
        return
    }

    ssh_config = &ssh.ClientConfig{
        User: sshClient.User,
        Auth: auth,
        HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
            return nil
        },
        Timeout:time.Second * DEFAULT_TIMEOUT,
    }

    client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshClient.Host, sshClient.Port), ssh_config)
    if err != nil {
        fmt.Println(err)
        return
    }

    session, err := client.NewSession()
    if err != nil {
        fmt.Println(err)
        client.Close()
        return
    }

    sshClient.session = session
    sshClient.client  = client
}

func (sshClient *SSH) RunCmd(cmd string) {
    out, err := sshClient.session.CombinedOutput(cmd)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println(string(out))
}

func (sshClient *SSH) Close() {
    sshClient.session.Close()
    sshClient.client.Close()
}



// MAIN EXECUTION //

func main() {
    homeDir := os.Getenv("HOME")
    sshCli := &SSH{
	Host :	"",
        User :	"",
        Port :	22,
        Credential: homeDir + "/.ssh/id_rsa",
    }
    sshCli.Connect(CERT_PUBLIC_KEY_FILE)
    sshCli.RunCmd("w")
    sshCli.Close()
}

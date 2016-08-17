package sftp

import(
    "github.com/pkg/sftp"
    "gosync/ssh"
    "golang.org/x/crypto/ssh"
    "log"
)

var sftpClient sftp.Client

func Init(){
    sshConfig := &ssh.ClientConfig{
        User: "jsmith",
        Auth: []ssh.AuthMethod{
            SSHAgent(),
        },
    }
    client := &SSHClient{
        Config: sshConfig,
        Host:   "example.com",
        Port:   22,
    }

    conn, err := client.GetConn()
    if err != nil {
        log.Fatalf("unable to connect to [%s]: %v", addr, err)
    }
    defer conn.Close()

    sftpClient, err := sftp.NewClient(conn, sftp.MaxPacket(*SIZE))
    if err != nil {
        log.Fatalf("unable to start sftp subsytem: %v", err)
    }
    defer sftpClient.Close()
}


func Sync(filepath string){
    fi, err := sftp.Lstat(filepath)
    if err != nil {
        log.Fatal(err)
    }
    log.Println(fi)
}
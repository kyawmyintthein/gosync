package ssh

import(
    "io"
    "strings"
    "fmt"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
)

type (
    SSHCommand struct{
        Path string
        Env []string
        Stdin  io.Reader
        Stdout io.Writer
        Stderr io.Writer
    }

    SSHClient struct{
        Config *ssh.ClientConfig
        Host    string
        Port    int
    }
)

// run ssh command
func (client *SSHClient) RunCommand(cmd *SSHCommand) error{
    var(
        session ssh.Session
        err error
    )

    if session, err := client.newSession(); err != nil{
        return err
    }
    defer session.Close()

    // prepareCommand
    err = client.prepareCommand(session, cmd)
    if err != nil{
        return err
    }

    err = session.Run(cmd.Path)
    return err
}

// prepare to run ssh command
func (client *SSHClient) preapareCommand(session *ssh.Session,cmd *SSHCommand) error{
    for _, env := range cmd.Env{
        val := strings.Split(env, "=")
        if len(val) != 2{
            continue
        }

        if err := session.Setenv(val[0]),val[1]); err != nil{
            return err
        }   
    }

    if cmd.Stdin != nil{
        stdin, err := session.StdinPipe()
        if err != nil{
            return fmt.Errorf("Unable to setup stdin for session %v", err)
        }
        go io.Copy(stdin, cmd.Stdin)
    }

    if cmd.Stdout != nil{
        stdout, err := session.StdoutPipe()
        if err != nil{
            return fmt.Errorf("Unable to setup stdout for session %v", err)
        }
        go io.Copy(stdout, cmd.Stdout)
    }

    if cmd.Stderr != nil{
        stderr, err := session.StderrPipe()
        if err != nil{
            return fmt.Errorf("Unable to setup stderr for session %v", err)
        } 
        go io.Copy(stderr, cmd.Stderr)
    }
    return nil
 }


// get new terminal session
func (client *SSHClient) newSession() (*ssh.Session, error){
    var(
        session *ssh.Session
        err     error
    )
    connection, err := client.GetConn()
    if err != nil{
        return  nil, err
    }

    session, err = connection.NewSession()
    if err != nil{
        return nil, fmt.Errorf("Failed to create session: %s", err)
    }

    modes := ssh.TerminalModes{
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err = session.RequestPty("xterm", 80, 40, modes); err != nil{
        session.Close()
        return nil, fmt.Errorf("request terminal failed: %s", err)
    }
    return session, nil
}


func (client *SSHClient)GetConn() (ssh.connection, error){
    connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d",client.Host,client.Port),client.Config)
    return connection, err
}

func SSHAgent() ssh.AuthMethod {
    if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
        return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
    }
    return nil
}
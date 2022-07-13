package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	// example input
	hostInput       = "127.0.0.1"               // host of ssh server
	portInput       = "22"                      // port of sftp server
	userInput       = "anan"                    // name of user on the ssh server
	passInput       = "yourcustompassword"      // password for the user on the ssh server
	pathInput       = "/home/user"              // path of sftp folder
	privateKeyInput = "/Users/user/.ssh/id_rsa" // path to client private key
	publicKeyInput  = "known_hosts"             // path to known_hosts where server public key is stored
)

func main() {

	var hostKey ssh.PublicKey

	hostKey, err := GetHostKey(hostInput)
	if err != nil {
		log.Fatal(err)
	}

	// Read privatKey of SSH client
	// and add it as authentication method
	//
	key, err := ioutil.ReadFile(privateKeyInput)
	if err != nil {
		log.Fatalf("unable to read private key: %v", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatalf("unable to parse private key: %v", err)
	}

	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig,
	// and provide a HostKeyCallback.
	config := &ssh.ClientConfig{
		User: userInput,
		Auth: []ssh.AuthMethod{
			// Auth method 1: password
			ssh.Password(passInput),
			// Auth method 2: Use the PublicKeys
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
		// HostKeyCallback: hostKey,
	}
	sshConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostInput, portInput), config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}
	defer sshConn.Close()

	// open an SFTP session over an existing ssh connection.
	client, err := sftp.NewClient(sshConn)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// walk a directory
	w := client.Walk(pathInput)
	for w.Step() {
		if w.Err() != nil {
			continue
		}
		log.Println(w.Path())
	}

	// leave your mark
	f, err := client.Create("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte("Hello world!")); err != nil {
		log.Fatal(err)
	}
	f.Close()

	// check it's there
	fi, err := client.Lstat("hello.txt")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(fi)
}

func GetHostKey(host string) (ssh.PublicKey, error) {
	file, err := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", publicKeyInput))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var hostKey ssh.PublicKey
	for scanner.Scan() {
		fields := strings.Split(scanner.Text(), " ")
		if len(fields) != 3 {
			continue
		}
		if strings.Contains(fields[0], host) {
			var err error
			hostKey, _, _, _, err = ssh.ParseAuthorizedKey(scanner.Bytes())
			if err != nil {
				return nil, fmt.Errorf("error parsing %q: %v", fields[2], err)
			}
			break
		}
	}

	if hostKey == nil {
		return nil, fmt.Errorf("no hostkey for %s", host)
	}
	return hostKey, nil
}

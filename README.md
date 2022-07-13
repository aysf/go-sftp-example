# Example of SFTP + SSH in Golang


This snippet is for sftp connection combined from the example of Golang package:

- https://pkg.go.dev/golang.org/x/crypto/ssh#example-Dial
- https://pkg.go.dev/golang.org/x/crypto/ssh#example-PublicKeys
- https://pkg.go.dev/github.com/pkg/sftp@v1.12.0#example-package
- https://stackoverflow.com/questions/45441735/ssh-handshake-complains-about-missing-host-key


# Usage

1. Clone this repository
2. Setup your sftp server
3. Adjust these variables with your own environment
```go
	hostInput       = "127.0.0.1"               // host of ssh server
	portInput       = "22"                      // port of sftp server
	userInput       = "anan"                    // name of user on the ssh server
	passInput       = "yourcustompassword"      // password for the user on the ssh server
	pathInput       = "/home/user"              // path of sftp folder
	privateKeyInput = "/Users/user/.ssh/id_rsa" // path to client private key
	publicKeyInput  = "known_hosts"             // path to known_hosts where server public key is stored
```
4. go to project folder and run `go run .`
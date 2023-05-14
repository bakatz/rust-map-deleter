package main

import (
	"context"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type Response struct {
	Message        string `json:"message"`
	Code           string `json:"error_code,omitempty"`
	NumMapsDeleted int    `json:"num_maps_deleted"`
}

func Handler(ctx context.Context) (Response, error) {
	// set up the SFTP server details
	hostPort := os.Getenv("SFTP_HOST_PORT")
	username := os.Getenv("SFTP_USERNAME")
	password := os.Getenv("SFTP_PASSWORD")
	gameServerBasePath := os.Getenv("GAME_SERVER_BASE_PATH")

	timeout := 5 * time.Second

	// set up the SSH client config
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}

	// connect to the SSH server with a timeout
	conn, err := ssh.Dial("tcp", hostPort, sshConfig)
	if err != nil {
		return Response{Message: "Error connecting to SFTP server", Code: "sftp_connection_error"}, err
	}
	defer conn.Close()

	// open an SFTP session on the SSH connection
	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return Response{Message: "Error opening SFTP session", Code: "sftp_session_error"}, err
	}
	defer sftpClient.Close()

	matches, err := sftpClient.Glob(gameServerBasePath + "/*.map")
	if err != nil {
		return Response{Message: "Error searching for map files", Code: "map_search_error"}, err
	}
	println(matches)

	numDeleted := 0
	for _, filePath := range matches {
		sftpClient.Rename(filePath, filePath+".softdeleted")
		numDeleted++
	}

	// return a success message
	return Response{Message: "Function finished without errors", NumMapsDeleted: numDeleted}, nil
}

func main() {
	godotenv.Load("../../.env")
	if os.Getenv("RUN_WITHOUT_LAMBDA") == "true" {
		Handler(context.TODO())
	} else {
		lambda.Start(Handler)
	}
}

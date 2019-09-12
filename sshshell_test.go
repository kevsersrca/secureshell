package Shell

import (
	"testing"
)

func TestClient_Exec(t *testing.T) {
	client, err := ConnectWithKeyFile("localhost:22", "root", "/home/user/.ssh/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	output, err := client.Exec("uptime -p")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}

func TestClient_ExecRoot(t *testing.T) {
	passwd := "root_password"
	client, err := ConnectWithKeyFile("localhost:22", "root", "/home/user/.ssh/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	output, err := client.ExecRoot("cat /root/root.txt", passwd)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}

func TestClient_ExecWithPassword(t *testing.T) {
	client, err := ConnectWithPassword("localhost:22", "root", "ssh_password")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	output, err := client.Exec("cat /root/root.txt")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(output)
}
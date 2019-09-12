# SSH Client Script

    SSH Shell with Go Programming

## Install

```go

    go get github.com/kevsersrca/secureshell
```

## Usage

```go

    package main

    import(
        s"github.com/kevsersrca/secureshell"
    )
    func main() {
            client, err := s.ConnectWithKeyFile("localhost:22", "root", "/home/user/.ssh/id_rsa")
        	if err != nil {
        		panic(err)
        	}
        	defer client.Close()
        
        	output, err := client.Exec("uptime -p")
        	if err != nil {
        		panic(err)
        	}
        	fmt.Println(output)
    }

```

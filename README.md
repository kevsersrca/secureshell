# SSH Client Script
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkevsersrca%2Fsecureshell.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkevsersrca%2Fsecureshell?ref=badge_shield)


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


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fkevsersrca%2Fsecureshell.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fkevsersrca%2Fsecureshell?ref=badge_large)
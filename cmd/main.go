package main

import (
	"bufio"
	"fmt"
	querymanager "jsondb/internal/query_manager"
	"os"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("query> ")

		scanner.Scan()
		src := scanner.Text()

		if src == "q" {
			break
		}

		queryManager := querymanager.New(src)

		err := queryManager.ExecuteQuery()

		if err != nil {
			panic(err)
		}

	}
}

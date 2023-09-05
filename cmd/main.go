package main

import (
	"bufio"
	"fmt"
	dm "jsondb/internal/database_manager"
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

		dbManager := dm.CreateNewDatabaseManager(src, "project1", "/Users/hassanelabdallah/golang/pracs/json_dbms/configs/data-sample")

		err := dbManager.ExecuteQuery()

		if err != nil {
			panic(err)
		}

	}
}

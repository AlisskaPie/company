package main

import (
	"companies/api"

	_ "github.com/lib/pq"
)

func main() {
	api.Run()
}

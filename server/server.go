package main

import (
    "./goserver"
    "./bloghandler"
    "./viewhandler"
)

func main() {
    bloghandler.Register()
    viewhandler.Register()
    goserver.Start(8080)
}

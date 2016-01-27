// [Environment variables](http://en.wikipedia.org/wiki/Environment_variable)
// are a universal mechanism for [conveying configuration
// information to Unix programs](http://www.12factor.net/config).
// Let's look at how to set, get, and list environment variables.

package main

import "fmt"
import "io/ioutil"

func main() {

	dat, err := ioutil.ReadFile("/etc/hosts")
	if err != nil {
		panic(err)
	}
	fmt.Print(string(dat))

}

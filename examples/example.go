// example program

package main

import (
	"fmt"
	"github.com/fzzy/radix/redis"
	"os"
	"time"
)

func errHndlr(err error) {
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}

func main() {
	// time.ParseDuration returns a time.Duration type (go1.2)
	duration, err := time.ParseDuration("10s")
	errHndlr(err)
	
	c, err := redis.DialTimeout("tcp", "127.0.0.1:6379", duration)
	errHndlr(err)
	defer c.Close()

	// select database
	r := c.Cmd("select", 8)
	errHndlr(r.Err)

	r = c.Cmd("flushdb")
	errHndlr(r.Err)

	r = c.Cmd("echo", "Hello world!")
	errHndlr(r.Err)

	s, err := r.Str()
	errHndlr(err)
	fmt.Println("echo:", s)

	//* Strings
	r = c.Cmd("set", "mykey0", "myval0")
	errHndlr(r.Err)

	s, err = c.Cmd("get", "mykey0").Str()
	errHndlr(err)
	fmt.Println("mykey0:", s)

	myhash := map[string]string{
		"mykey1": "myval1",
		"mykey2": "myval2",
		"mykey3": "myval3",
	}

	// Alternatively:
	// c.Cmd("mset", "mykey1", "myval1", "mykey2", "myval2", "mykey3", "myval3")
	r = c.Cmd("mset", myhash)
	errHndlr(r.Err)

	ls, err := c.Cmd("mget", "mykey1", "mykey2", "mykey3").List()
	errHndlr(err)
	fmt.Println("mykeys values:", ls)

	//* List handling
	mylist := []string{"foo", "bar", "qux"}

	// Alternativaly:
	// c.Cmd("rpush", "mylist", "foo", "bar", "qux")
	r = c.Cmd("rpush", "mylist", mylist)
	errHndlr(r.Err)

	mylist, err = c.Cmd("lrange", "mylist", 0, -1).List()
	errHndlr(err)

	fmt.Println("mylist:", mylist)

	//* Hash handling

	// Alternatively:
	// c.Cmd("hmset", "myhash", ""mykey1", "myval1", "mykey2", "myval2", "mykey3", "myval3")
	r = c.Cmd("hmset", "myhash", myhash)
	errHndlr(r.Err)

	myhash, err = c.Cmd("hgetall", "myhash").Hash()
	errHndlr(err)

	fmt.Println("myhash:", myhash)

	//* Pipelining
	c.Append("set", "multikey", "multival")
	c.Append("get", "multikey")

	c.GetReply()     // set
	r = c.GetReply() // get
	errHndlr(r.Err)

	s, err = r.Str()
	errHndlr(err)
	fmt.Println("multikey:", s)
}

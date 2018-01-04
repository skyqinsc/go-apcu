package main

import (
	"apcu"
	"fmt"
	"time"
)

func main() {
	fmt.Println("====")
	apcu.GlobalCache = apcu.NewCache(10*time.Second, 1*time.Second)
	fmt.Println("====")
	apcu.GlobalCache.Store("one", 1, 5*time.Second)
	fmt.Println("====")
	apcu.GlobalCache.Store("two", 2, 10*time.Second)
	apcu.GlobalCache.Store("three", 3, 0*time.Second)
	fmt.Println("====")
	time.Sleep(3 * time.Second)
	fmt.Println("====")
	val, ok := apcu.GlobalCache.Fetch("one")
	fmt.Printf("one-3  %v = %v \n", val, ok)

	time.Sleep(4 * time.Second)
	val, ok = apcu.GlobalCache.Fetch("one")
	fmt.Printf("one-7  %v = %v \n", val, ok)
	val, ok = apcu.GlobalCache.Fetch("two")
	fmt.Printf("two-7  %v = %v \n", val, ok)

	time.Sleep(5 * time.Second)
	val, ok = apcu.GlobalCache.Fetch("two")
	fmt.Printf("two-12  %v = %v \n", val, ok)
	val, ok = apcu.GlobalCache.Fetch("three")
	fmt.Printf("three-12  %v = %v \n", val, ok)
}

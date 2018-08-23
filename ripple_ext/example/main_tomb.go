package main

import (
	"gopkg.in/tomb.v1"
	"time"
	"fmt"
)

var (
	t1 tomb.Tomb
	t2 tomb.Tomb
)

func main()  {
	start()

	fmt.Println("sleep 2 second to exit...")
	time.Sleep(time.Second*2)

	stop()
	time.Sleep(time.Second*4)
	fmt.Println("quit")
}

func start()  {
	fmt.Println("start...")
	go routine1()
	go routine2()
	fmt.Println("start end...")
}

func stop()  {
	fmt.Println("stop...")
	t1.Kill(nil)
	//t2.Kill(nil)
	t1.Wait()
	//t2.Wait()
	fmt.Println("stop end...")
}

func routine1()  {
	fmt.Println("routine1...")
	defer fmt.Println("routine1 end...")
	defer t1.Done()
	for{
		select {
		case i, j := <-t1.Dying():
			fmt.Println("routine1 dying", i, j)
			time.Sleep(time.Second)
		default:
			fmt.Println("routine1 default")
			time.Sleep(time.Second)
		}
	}
}

func routine2()  {
	fmt.Println("routine2...")
	defer fmt.Println("routine2 end...")
	defer t2.Done()
	for{
		select {
		case <-t2.Dying():
			fmt.Println("routine2 dying")
		default:
			fmt.Println("routine2 default")
			time.Sleep(time.Second)
		}
	}
}
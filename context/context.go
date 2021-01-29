package main

import (
	"context"
	"fmt"
	"time"
)

func main ()  {

	WithDemo()
}
func WithTime()  {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	fmt.Println(ctx)
	fmt.Println(ctx.Deadline())
	fmt.Println(ctx.Value("name"))
	select {
	case <-time.After(time.Second * 2):
		//cancel()
		fmt.Println("sleep", ctx)
	case <-ctx.Done():
		fmt.Println("done:", ctx.Err())
	}

}
func WithCancel()  {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	fmt.Println(ctx)
	i := 0
	for{
		select {
		case <-time.After(time.Second *2):
			i++
			fmt.Println(i)
			if i > 5 {
				cancelFunc()
				break
			}
		case <-ctx.Done():
			fmt.Println("err")
			break
		}
	}


	fmt.Println(ctx.Err())
}


func WithDemo()  {
	f := struct {
		Name string
		Value string
	}{
		"eros",
		"1000",
	}
	ctx := context.WithValue(context.Background(), "key", f)

	fmt.Println(ctx.Value("key"))
}

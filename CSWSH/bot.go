package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"

	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/runner"
)

func main() {
	fmt.Println("start bot")
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(15 * time.Second)
		os.Exit(0)
	}()

	// create chrome instance
	c, err := chromedp.New(ctxt,
		chromedp.WithRunnerOptions(
			runner.Path("/usr/bin/chromium"),
			runner.Flag("headless", true),
			runner.Flag("disable-gpu", true),
			runner.Flag("no-sandbox", true),
		))
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	err = c.Run(ctxt, visit())
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func visit() chromedp.Tasks {
	return chromedp.Tasks{

		chromedp.ActionFunc(func(ctxt context.Context, h cdp.Executor) error {
			success, err := network.SetCookie("token", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTI2MDIxNDU2LCJuYW1lIjoiQW15YW5nIn0.JJmKn7DuM1VbriXeG4XqT18ycDdObdaE1fltp2CIGAY").
				WithDomain("127.0.0.1").
				WithHTTPOnly(true).
				Do(ctxt, h)
			if err != nil {
				return err
			}
			if !success {
				return errors.New("could not set cookie")
			}
			return nil
		}),

		chromedp.Navigate(`http://127.0.0.1:8001/admin`),
		chromedp.Sleep(5 * time.Second),
	}
}

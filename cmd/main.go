/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"github.com/recordbase/recordbase"
	"github.com/recordbase/recordbasepb"
	"os"
)

func doMain() error {

	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		Rand:               rand.Reader,
	}

	token := os.Getenv("RECORDBASE_AUTH")

	client, err := recordbase.NewClient(context.Background(), "127.0.0.1:8500", token, tlsConfig)
	if err != nil {
		return err
	}

	fmt.Printf("client: %v\n", client)

	defer client.Destroy()

	resp, err := client.Get(context.Background(), &recordbasepb.GetRequest{
		Tenant:     "w",
		PrimaryKey: "alex",
	})
	if err != nil {
		return err
	}

	fmt.Printf("resp = %s\n", resp.String())

	return nil
}

func main() {

	if err := doMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

}



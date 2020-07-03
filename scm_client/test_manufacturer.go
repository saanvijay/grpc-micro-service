package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"scmpb"
	"time"
)

func AddManufacturer(client scmpb.ScmServiceClient) string {

	fmt.Println("Add Manufacturer.....")
	var req scmpb.ManufacturerRequest
	req.ManufacturerName = "Vijay"
	req.ManufacturerRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.ManufacturerStatus = "Requsting products"

	res, err := client.AddManufacturer(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Manufacturer server %v\n", res)

	return res.GetManufacturerId()
}

func UpdateManufacturer(client scmpb.ScmServiceClient, id string) {
	fmt.Println("UpdateManufacturer ....")
	var req scmpb.ManufacturerRequest
	req.ManufacturerId = id
	req.ManufacturerName = "Prakash"
	req.ManufacturerRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.ManufacturerStatus = "Update Manufacturer"

	res, err := client.UpdateManufacturer(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Manufacturer server %v\n", res)

}

func DeleteManufacturer(client scmpb.ScmServiceClient, id string) {
	fmt.Println("DeleteManufacturer ...")
	var sreq scmpb.ManufacturerRequest
	sreq.ManufacturerId = id
	sres, err := client.DeleteManufacturer(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Manufacturer server %v\n", sres)
}

func GetManufacturer(client scmpb.ScmServiceClient, id string) {
	fmt.Println("GetManufacturer ....")
	var sreq scmpb.ManufacturerRequest
	sreq.ManufacturerId = id
	sres, err := client.GetManufacturer(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Manufacturer server %v\n", sres)
}

func ListAllManufacturers(client scmpb.ScmServiceClient) {
	fmt.Println("List all Manufacturers ...")
	var sreq1 scmpb.ManufacturerRequest
	sres1, err := client.ListAllManufacturers(context.Background(), &sreq1)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	for {
		streamRes, err := sres1.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Response from Manufacturer server streaming %v\n", streamRes)
	}
}

// Testmanufacturer tests all funcs in Manufacturer
func TestManufacturer(client scmpb.ScmServiceClient) {
	id := AddManufacturer(client)
	GetManufacturer(client, id)
	UpdateManufacturer(client, id)
	DeleteManufacturer(client, id)
	ListAllManufacturers(client)
	id = AddManufacturer(client)
	ListAllManufacturers(client)
	UpdateManufacturer(client, id)

}

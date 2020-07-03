package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"scmpb"
	"time"
)

func AddTransporter(client scmpb.ScmServiceClient) string {

	fmt.Println("Add Transporter.....")
	var req scmpb.TransporterRequest
	req.TransporterName = "Vijay"
	req.TransporterRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.TransporterStatus = "Requsting Raw materials"

	res, err := client.AddTransporter(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Transporter server %v\n", res)

	return res.GetTransporterId()
}

func UpdateTransporter(client scmpb.ScmServiceClient, id string) {
	fmt.Println("UpdateTransporter ....")
	var req scmpb.TransporterRequest
	req.TransporterId = id
	req.TransporterName = "DTDC"
	req.TransporterRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.TransporterStatus = "Update Transporter"

	res, err := client.UpdateTransporter(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Transporter server %v\n", res)

}

func DeleteTransporter(client scmpb.ScmServiceClient, id string) {
	fmt.Println("DeleteTransporter ...")
	var sreq scmpb.TransporterRequest
	sreq.TransporterId = id
	sres, err := client.DeleteTransporter(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Transporter server %v\n", sres)
}

func GetTransporter(client scmpb.ScmServiceClient, id string) {
	fmt.Println("GetTransporter ....")
	var sreq scmpb.TransporterRequest
	sreq.TransporterId = id
	sres, err := client.GetTransporter(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Transporter server %v\n", sres)
}

func ListAllTransporters(client scmpb.ScmServiceClient) {
	fmt.Println("List all Transporters ...")
	var sreq1 scmpb.TransporterRequest
	sres1, err := client.ListAllTransporters(context.Background(), &sreq1)
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
		fmt.Printf("Response from Transporter server streaming %v\n", streamRes)
	}
}

// TestTransporter tests all funcs in Transporter
func TestTransporter(client scmpb.ScmServiceClient) {
	id := AddTransporter(client)
	GetTransporter(client, id)
	UpdateTransporter(client, id)
	ListAllTransporters(client)
	DeleteTransporter(client, id)
	ListAllTransporters(client)
	id = AddTransporter(client)
	UpdateTransporter(client, id)
}

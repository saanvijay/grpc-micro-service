package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"supplychainpb"
	"time"
)

func AddSupplier(client supplychainpb.ScmServiceClient) string {

	fmt.Println("Add Supplier.....")
	var req supplychainpb.SupplierRequest
	req.SupplierName = "Vijay"
	req.SupplierRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.SupplierStatus = "Requsting Raw materials"

	res, err := client.AddSupplier(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from supplier server %v\n", res)

	return res.GetSupplierId()
}

func UpdateSupplier(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("UpdateSupplier ....")
	var req supplychainpb.SupplierRequest
	req.SupplierId = id
	req.SupplierName = "Prakash"
	req.SupplierRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.SupplierStatus = "Update supplier"

	res, err := client.UpdateSupplier(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from supplier server %v\n", res)

}

func DeleteSupplier(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("DeleteSupplier ...")
	var sreq supplychainpb.SupplierRequest
	sreq.SupplierId = id
	sres, err := client.DeleteSupplier(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from supplier server %v\n", sres)
}

func GetSupplier(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("GetSupplier ....")
	var sreq supplychainpb.SupplierRequest
	sreq.SupplierId = id
	sres, err := client.GetSupplier(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from supplier server %v\n", sres)
}

func ListAllSuppliers(client supplychainpb.ScmServiceClient) {
	fmt.Println("List all suppliers ...")
	var sreq1 supplychainpb.SupplierRequest
	sres1, err := client.ListAllSuppliers(context.Background(), &sreq1)
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
		fmt.Printf("Response from supplier server streaming %v\n", streamRes)
	}
}

// TestSupplier tests all funcs in supplier
func TestSupplier(client supplychainpb.ScmServiceClient) {
	id := AddSupplier(client)
	GetSupplier(client, id)
	UpdateSupplier(client, id)
	DeleteSupplier(client, id)
	ListAllSuppliers(client)
	DeleteSupplier(client, id)
	ListAllSuppliers(client)
	id = AddSupplier(client)
	UpdateSupplier(client, id)
}

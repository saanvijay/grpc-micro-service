package main

import (
	"context"
	"fmt"
	"log"
	"supplychainpb"
)

func CreateNetwork(client supplychainpb.ScmServiceClient) bool {

	fmt.Println("Create Network.....")
	var req supplychainpb.NetworkRequest
	req.Consortium = "TestConsortium1"
	req.NetworkName = "TestNetwork1"
	req.PeersPerOrg = 2
	req.ChannelName = "TestChannel1"
	org1 := supplychainpb.NetworkOrg{OrgName: "Supplier1", OrgType: "Supplier", OrgProperties: "Supplier1_Address"}
	org2 := supplychainpb.NetworkOrg{OrgName: "Manufacturer1", OrgType: "Manufacturer", OrgProperties: "Manufacturer1_Address"}
	var orgList []*supplychainpb.NetworkOrg
	orgList = append(orgList, &org1)
	orgList = append(orgList, &org2)
	req.OrgList = orgList

	res, err := client.CreateNetwork(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from scm server %v\n", res)

	return res.GetResult()
}

// TestNetwork tests all funcs in network
func TestNetwork(client supplychainpb.ScmServiceClient) {
	CreateNetwork(client)
	//StopBlockchainNetwork(client)
	//StartBlockchainNetwork(client)
	//ListAllNetwork(client)
}

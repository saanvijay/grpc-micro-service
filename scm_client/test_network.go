package main

import (
	"context"
	"fmt"
	"log"
	"scmpb"
)

func CreateNetwork(client scmpb.ScmServiceClient) bool {

	fmt.Println("Create Network.....")
	var req scmpb.NetworkRequest
	req.Consortium = "TestConsortium1"
	req.NetworkName = "TestNetwork1"
	req.PeersPerOrg = 2
	req.ChannelName = "TestChannel1"
	org1 := scmpb.NetworkOrg{OrgName: "Supplier1", OrgType: "Supplier", OrgProperties: "Supplier1_Address"}
	org2 := scmpb.NetworkOrg{OrgName: "Manufacturer1", OrgType: "Manufacturer", OrgProperties: "Manufacturer1_Address"}
	var orgList []*scmpb.NetworkOrg
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
func TestNetwork(client scmpb.ScmServiceClient) {
	CreateNetwork(client)
	//StopBlockchainNetwork(client)
	//StartBlockchainNetwork(client)
	//ListAllNetwork(client)
}

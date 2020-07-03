package main

import (
	"bgclient"
	"context"
	"fmt"
	"log"
	"os"
	"scmpb"
)

func (n *server) CreateNetwork(ctx context.Context, req *scmpb.NetworkRequest) (*scmpb.NetworkResponse, error) {

	var network bgclient.BlockchainNetwork
	network.Consortium = req.GetConsortium()
	network.Name = req.GetNetworkName()
	network.ChannelName = req.GetChannelName()
	network.PeersPerOrg = req.GetPeersPerOrg()
	orgList := req.GetOrgList()
	// Process, check, compute data
	// for each raw material check the stock status
	for _, org := range orgList {
		network.OrgList = append(network.OrgList, bgclient.Organization{org.GetOrgName(), org.GetOrgType(), org.GetOrgProperties()})
	}

	loghandler := log.New(os.Stdout, "bgclient-api ", log.LstdFlags)
	client := bgclient.NewBlockchainGatewayClient(loghandler, "http://localhost:10000")
	fmt.Println(network.Consortium, network.Name, network.ChannelName, network.PeersPerOrg, orgList)
	// Call Blockchain createNetwork
	client.GetBGHttpClient()
	client.CreateBlockchainNetwork(network)

	// Update local database if anything goes as off-chain data

	return &scmpb.NetworkResponse{
		Result: true,
	}, nil
}

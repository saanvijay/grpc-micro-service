package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"supplychainpb"
	"time"
)

func AddAgent(client supplychainpb.ScmServiceClient) string {

	fmt.Println("Add Agent.....")
	var req supplychainpb.AgentRequest
	var agent supplychainpb.Agent
	agent.AgentName = "vijay"
	req.Agent = &agent
	req.AgentRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.AgentStatus = "Requsting products"

	res, err := client.AddAgent(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Agent server %v\n", res)

	return res.Agent.GetAgentId()
}

func UpdateAgent(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("UpdateAgent ....")
	var req supplychainpb.AgentRequest
	var agent supplychainpb.Agent
	agent.AgentId = id
	agent.AgentName = "Prakash"
	req.Agent = &agent
	req.AgentRequestedTime = time.Now().Format("01-02-2006 15:04:05 Monday")
	req.AgentStatus = "Update Agent"

	res, err := client.UpdateAgent(context.Background(), &req)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Agent server %v\n", res)

}

func DeleteAgent(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("DeleteAgent ...")
	var sreq supplychainpb.AgentRequest
	var agent supplychainpb.Agent
	agent.AgentId = id
	sreq.Agent = &agent
	sres, err := client.DeleteAgent(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Agent server %v\n", sres)
}

func GetAgent(client supplychainpb.ScmServiceClient, id string) {
	fmt.Println("GetAgent ....")
	var sreq supplychainpb.AgentRequest
	var agent supplychainpb.Agent
	agent.AgentId = id
	sreq.Agent = &agent
	sres, err := client.GetAgent(context.Background(), &sreq)
	if err != nil {
		log.Fatalf("Request error %v\n", err)
	}
	fmt.Printf("Response from Agent server %v\n", sres)
}

func ListAllAgents(client supplychainpb.ScmServiceClient) {
	fmt.Println("List all Agents ...")
	var sreq1 supplychainpb.AgentRequest
	sres1, err := client.ListAllAgents(context.Background(), &sreq1)
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
		fmt.Printf("Response from Agent server streaming %v\n", streamRes)
	}
}

// TestAgent tests all funcs in Agent
func TestAgent(client supplychainpb.ScmServiceClient) {
	id := AddAgent(client)
	GetAgent(client, id)
	UpdateAgent(client, id)
	DeleteAgent(client, id)
	ListAllAgents(client)
	id = AddAgent(client)
	ListAllAgents(client)
	UpdateAgent(client, id)

}

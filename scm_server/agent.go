package main

import (
	"context"
	"log"
	"scmpb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (agServer *server) GetAgent(ctx context.Context, req *scmpb.AgentRequest) (*scmpb.AgentResponse, error) {

	// Collection for transporter
	collection := mClient.Database("scmdb").Collection("agent")
	// Get data from request
	agentID, err := primitive.ObjectIDFromHex(req.Agent.GetAgentId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := collection.FindOne(context.TODO(), bson.M{"_id": agentID})
	outData := scmpb.AgentResponse{}

	if err := result.Decode(&outData); err != nil {
		return nil, err
	}

	// Get data from request
	agent := req.GetAgent()
	reqProducts := req.GetProducts()
	status := "Agent requested"

	// Send processed, computed data
	res := scmpb.AgentResponse{
		Agent:              agent,
		AgentProperties:    req.GetAgentProperties(),
		Products:           reqProducts,
		AgentRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		AgentStatus:        status,
	}
	return &res, nil
}

func (agServer *server) AddAgent(ctx context.Context, req *scmpb.AgentRequest) (*scmpb.AgentResponse, error) {

	collection := mClient.Database("scmdb").Collection("agent")
	// Get data from request
	agent := req.GetAgent()
	agentProps := req.GetAgentProperties()
	resLogistics := req.GetProducts()

	status := "Agent added"

	// Send processed, computed data
	res := scmpb.AgentResponse{
		Agent:              agent,
		AgentProperties:    agentProps,
		Products:           resLogistics,
		AgentRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		AgentStatus:        status,
	}

	result, err := collection.InsertOne(context.TODO(), res)
	if err != nil {
		log.Fatal(err)
	}
	objID := result.InsertedID.(primitive.ObjectID)
	res.Agent.AgentId = objID.Hex()

	//Update Agent ID, since it is generated from MongoDB
	agServer.UpdateAgent(ctx,
		&scmpb.AgentRequest{
			Agent:              res.Agent,
			AgentProperties:    res.AgentProperties,
			Products:           res.Products,
			AgentRequestedTime: res.AgentRespondedTime,
			AgentStatus:        res.AgentStatus})

	return &res, nil
}

func (agServer *server) DeleteAgent(ctx context.Context, req *scmpb.AgentRequest) (*scmpb.AgentResponse, error) {

	collection := mClient.Database("scmdb").Collection("agent")
	agentID, err := primitive.ObjectIDFromHex(req.Agent.GetAgentId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": agentID})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &scmpb.AgentResponse{
		Agent:       &scmpb.Agent{AgentId: agentID.Hex()},
		AgentStatus: "Agent Deleted",
	}, nil

}

func (agServer *server) UpdateAgent(ctx context.Context, req *scmpb.AgentRequest) (*scmpb.AgentResponse, error) {

	collection := mClient.Database("scmdb").Collection("agent")
	agentID, err := primitive.ObjectIDFromHex(req.Agent.GetAgentId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	agentName := req.Agent.GetAgentName()
	reqProducts := req.GetProducts()
	agentType := req.Agent.GetAgentType()

	// Process, check, compute data
	var resProducts []*scmpb.OutboundLogistics
	// for each raw material check the stock status
	for _, outProduct := range reqProducts {
		// Get the location details from localDB or blockchain and sensor data from MQTT
		outProduct.Product.LogisticInStock = true
		resProducts = append(resProducts, outProduct)

	}
	status := "Agent updated"

	// Send processed, computed data
	update := scmpb.AgentResponse{
		//AgentId:   AgentID,
		Agent:              &scmpb.Agent{AgentId: agentID.Hex(), AgentName: agentName, AgentType: agentType},
		Products:           resProducts,
		AgentRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		AgentStatus:        status,
	}

	filter := bson.M{"_id": agentID}
	_ = collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	return &update, nil
}

func (agServer *server) ListAllAgents(req *scmpb.AgentRequest, stream scmpb.ScmService_ListAllAgentsServer) error {
	collection := mClient.Database("scmdb").Collection("agent")
	data := &scmpb.AgentResponse{}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
		return err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		err := cursor.Decode(data)
		if err != nil {
			log.Fatal(err)
			return err
		}
		stream.Send(&scmpb.AgentResponse{
			Agent:              &scmpb.Agent{AgentId: data.Agent.GetAgentId(), AgentName: data.Agent.GetAgentName()},
			AgentRespondedTime: data.AgentRespondedTime,
			AgentStatus:        data.AgentStatus,
		})
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

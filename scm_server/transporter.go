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

func (transServer *server) GetTransporter(ctx context.Context, req *scmpb.TransporterRequest) (*scmpb.TransporterResponse, error) {

	// Collection for transporter
	collection := mClient.Database("scmdb").Collection("transporter")
	// Get data from request
	transporterID, err := primitive.ObjectIDFromHex(req.GetTransporterId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := collection.FindOne(context.TODO(), bson.M{"_id": transporterID})
	outData := scmpb.TransporterResponse{}

	if err := result.Decode(&outData); err != nil {
		return nil, err
	}

	// Get data from request
	transporterName := req.GetTransporterName()
	reqLogistics := req.GetLogistics()
	status := "Request transporter"

	// Send processed, computed data
	res := &scmpb.TransporterResponse{
		TransporterId:            transporterID.Hex(),
		TransporterName:          transporterName,
		Logistics:                reqLogistics,
		TransporterRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		TransporterStatus:        status,
	}
	return res, nil
}

func (transServer *server) AddTransporter(ctx context.Context, req *scmpb.TransporterRequest) (*scmpb.TransporterResponse, error) {

	collection := mClient.Database("scmdb").Collection("transporter")

	// Get data from request
	//transporterID := req.GetTransporterId()
	transporterName := req.GetTransporterName()
	reqLogistics := req.GetLogistics()

	// Process, check, compute data
	var resLogistics []*scmpb.Logistics
	// for each raw material check the stock status
	for _, logistic := range reqLogistics {
		// Get the location details from localDB or blockchain and sensor data from MQTT
		logistic.LogisticLocation.Latitude = 34235
		logistic.LogisticLocation.Latitude = 4545345
		logistic.LogisticProperties = "Sensor data goes here"

		resLogistics = append(resLogistics, logistic)

	}
	status := "Transporter added"

	// Send processed, computed data
	res := scmpb.TransporterResponse{
		//TransporterId:            transporterID,
		TransporterName:          transporterName,
		Logistics:                resLogistics,
		TransporterRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		TransporterStatus:        status,
	}

	result, err := collection.InsertOne(context.TODO(), res)
	if err != nil {
		log.Fatal(err)
	}
	objID := result.InsertedID.(primitive.ObjectID)
	res.TransporterId = objID.Hex()

	//Update Transporter ID, since it is generated from MongoDB
	transServer.UpdateTransporter(ctx,
		&scmpb.TransporterRequest{
			TransporterId:            res.TransporterId,
			TransporterName:          res.TransporterName,
			Logistics:                res.Logistics,
			TransporterRequestedTime: res.TransporterRespondedTime,
			TransporterStatus:        res.TransporterStatus})

	return &res, nil
}

func (transServer *server) DeleteTransporter(ctx context.Context, req *scmpb.TransporterRequest) (*scmpb.TransporterResponse, error) {

	collection := mClient.Database("scmdb").Collection("transporter")
	TransporterID, err := primitive.ObjectIDFromHex(req.GetTransporterId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": TransporterID})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &scmpb.TransporterResponse{
		TransporterId:     TransporterID.Hex(),
		TransporterStatus: "Deleted",
	}, nil

}

func (transServer *server) UpdateTransporter(ctx context.Context, req *scmpb.TransporterRequest) (*scmpb.TransporterResponse, error) {

	collection := mClient.Database("scmdb").Collection("transporter")
	TransporterID, err := primitive.ObjectIDFromHex(req.GetTransporterId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	//TransporterID := req.GetTransporterId()
	TransporterName := req.GetTransporterName()
	reqLogistics := req.GetLogistics()

	// Process, check, compute data
	var resLogistics []*scmpb.Logistics
	// for each raw material check the stock status
	for _, logistic := range reqLogistics {
		// Get the location details from localDB or blockchain and sensor data from MQTT
		logistic.LogisticLocation.Latitude = 34235
		logistic.LogisticLocation.Latitude = 4545345
		logistic.LogisticType = "fragile"
		logistic.LogisticProperties = "Sensor data goes here"

		resLogistics = append(resLogistics, logistic)
	}
	status := "Transporter updated"

	// Send processed, computed data
	update := scmpb.TransporterResponse{
		TransporterId:            TransporterID.Hex(),
		TransporterName:          TransporterName,
		Logistics:                resLogistics,
		TransporterRespondedTime: time.Now().Format("01-02-2006 15:04:05 Monday"),
		TransporterStatus:        status,
	}

	filter := bson.M{"_id": TransporterID}
	_ = collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	return &update, nil
}

func (transServer *server) ListAllTransporters(req *scmpb.TransporterRequest, stream scmpb.ScmService_ListAllTransportersServer) error {
	collection := mClient.Database("scmdb").Collection("transporter")
	data := &scmpb.TransporterResponse{}

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
		stream.Send(&scmpb.TransporterResponse{
			TransporterId:            data.GetTransporterId(),
			TransporterName:          data.TransporterName,
			TransporterRespondedTime: data.TransporterRespondedTime,
			TransporterStatus:        data.TransporterStatus,
			Logistics:                data.Logistics,
		})
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

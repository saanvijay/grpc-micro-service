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

func (manuServer *server) GetManufacturer(ctx context.Context, req *scmpb.ManufacturerRequest) (*scmpb.ManufacturerResponse, error) {

	// Collection for manufacturer
	collection := mClient.Database("scmdb").Collection("manufacturer")
	// Get data from request
	manufacturerID, err := primitive.ObjectIDFromHex(req.GetManufacturerId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := collection.FindOne(context.TODO(), bson.M{"_id": manufacturerID})
	outData := scmpb.ManufacturerResponse{}

	if err := result.Decode(&outData); err != nil {
		return nil, err
	}

	// Get data from request
	manufacturerName := outData.GetManufacturerName()
	reqProducts := outData.GetManufacturerProducts()
	associatedOwner := outData.GetManufacturerAssociatedOwner()
	status := "Requesting manufacturer"

	// Send processed, computed data
	res := &scmpb.ManufacturerResponse{
		ManufacturerId:              manufacturerID.Hex(),
		ManufacturerName:            manufacturerName,
		ManufacturerProducts:        reqProducts,
		ManufacturerAssociatedOwner: associatedOwner,
		ManufacturerRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		ManufacturerStatus:          status,
	}
	return res, nil
}

func (manuServer *server) AddManufacturer(ctx context.Context, req *scmpb.ManufacturerRequest) (*scmpb.ManufacturerResponse, error) {

	// Collection for manufacturer
	collection := mClient.Database("scmdb").Collection("manufacturer")

	// Get data from request
	manufacturerName := req.GetManufacturerName()
	reqProducts := req.GetManufacturerProducts()
	associatedOwner := req.GetManufacturerAssociatedOwner()

	// Process, check, compute data
	var resProducts []*scmpb.OutboundLogistics
	// for each raw material check the stock status
	for _, prod := range resProducts {
		// Get the In-stock details from localDB or blockchain
		if prod.GetProduct().GetLogisticInStock() {
			// material is in stock change the status and send this material to manufacturer via transporter
			resProducts = append(resProducts, prod)
		}
	}
	status := "Manufacturer added"

	// Send processed, computed data
	res := scmpb.ManufacturerResponse{
		ManufacturerName:            manufacturerName,
		ManufacturerProducts:        reqProducts,
		ManufacturerAssociatedOwner: associatedOwner,
		ManufacturerRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		ManufacturerStatus:          status,
	}

	result, err := collection.InsertOne(context.TODO(), res)
	if err != nil {
		log.Fatal(err)
	}
	objID := result.InsertedID.(primitive.ObjectID)
	res.ManufacturerId = objID.Hex()

	//Update Manufacturer ID, since it is generated from MongoDB
	manuServer.UpdateManufacturer(ctx,
		&scmpb.ManufacturerRequest{
			ManufacturerId:              res.ManufacturerId,
			ManufacturerName:            res.ManufacturerName,
			ManufacturerProducts:        res.ManufacturerProducts,
			ManufacturerAssociatedOwner: res.ManufacturerAssociatedOwner,
			ManufacturerRequestedTime:   res.ManufacturerRespondedTime,
			ManufacturerStatus:          res.ManufacturerStatus})
	return &res, nil

}

func (manuServer *server) DeleteManufacturer(ctx context.Context, req *scmpb.ManufacturerRequest) (*scmpb.ManufacturerResponse, error) {

	collection := mClient.Database("scmdb").Collection("manufacturer")
	manufacturerID, err := primitive.ObjectIDFromHex(req.GetManufacturerId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": manufacturerID})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &scmpb.ManufacturerResponse{
		ManufacturerId:     manufacturerID.Hex(),
		ManufacturerStatus: "Manufacturer Deleted",
	}, nil

}

func (manuServer *server) UpdateManufacturer(ctx context.Context, req *scmpb.ManufacturerRequest) (*scmpb.ManufacturerResponse, error) {

	// Collection for manufacturer
	collection := mClient.Database("scmdb").Collection("manufacturer")
	manufacturerID, err := primitive.ObjectIDFromHex(req.GetManufacturerId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// Get data from request
	manufacturerName := req.GetManufacturerName()
	reqProducts := req.GetManufacturerProducts()
	associatedOwner := req.GetManufacturerAssociatedOwner()

	// Process, check, compute data
	var resProducts []*scmpb.OutboundLogistics
	// for each raw material check the stock status
	for _, prod := range resProducts {
		// Get the In-stock details from localDB or blockchain
		if prod.GetProduct().GetLogisticInStock() {
			// material is in stock change the status and send this material to manufacturer via transporter
			resProducts = append(resProducts, prod)
		}
	}
	status := "Manufacturer updated"

	// Send processed, computed data
	update := scmpb.ManufacturerResponse{
		ManufacturerId:              manufacturerID.Hex(),
		ManufacturerName:            manufacturerName,
		ManufacturerProducts:        reqProducts,
		ManufacturerAssociatedOwner: associatedOwner,
		ManufacturerRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		ManufacturerStatus:          status,
	}

	filter := bson.M{"_id": manufacturerID}
	_ = collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	return &update, nil

}

func (manuServer *server) ListAllManufacturers(req *scmpb.ManufacturerRequest, stream scmpb.ScmService_ListAllManufacturersServer) error {
	collection := mClient.Database("scmdb").Collection("manufacturer")
	data := &scmpb.ManufacturerResponse{}

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
		stream.Send(&scmpb.ManufacturerResponse{
			ManufacturerId:              data.GetManufacturerId(),
			ManufacturerName:            data.GetManufacturerName(),
			ManufacturerAssociatedOwner: data.GetManufacturerAssociatedOwner(),
			ManufacturerRespondedTime:   data.GetManufacturerRespondedTime(),
			ManufacturerStatus:          data.GetManufacturerStatus(),
		})
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

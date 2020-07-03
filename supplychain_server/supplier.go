package main

import (
	"context"
	"log"
	"supplychainpb"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (supplierServer *server) GetSupplier(ctx context.Context, req *supplychainpb.SupplierRequest) (*supplychainpb.SupplierResponse, error) {

	// Collection for supplier
	collection := mClient.Database("scmdb").Collection("supplier")
	// Get data from request
	supplierID, err := primitive.ObjectIDFromHex(req.GetSupplierId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	result := collection.FindOne(context.TODO(), bson.M{"_id": supplierID})
	outData := supplychainpb.SupplierResponse{}

	if err := result.Decode(&outData); err != nil {
		return nil, err
	}

	supplierName := outData.GetSupplierName()
	outRawMaterials := outData.GetRawmaterials()
	associatedOwner := outData.GetSupplierAssociatedOwner()

	status := "Request supplier"

	// Send processed, computed data
	res := &supplychainpb.SupplierResponse{
		SupplierId:              supplierID.Hex(),
		SupplierName:            supplierName,
		Rawmaterials:            outRawMaterials,
		SupplierAssociatedOwner: associatedOwner,
		SupplierRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		SupplierStatus:          status,
	}
	return res, nil
}

func (supplierServer *server) AddSupplier(ctx context.Context, req *supplychainpb.SupplierRequest) (*supplychainpb.SupplierResponse, error) {

	collection := mClient.Database("scmdb").Collection("supplier")
	// Get data from request
	//supplierID := req.GetSupplierId()
	supplierName := req.GetSupplierName()
	//supplierReqTime := req.GetSupplierRequestedTime()
	//supplierStatus := req.GetSupplierStatus()
	reqRawMaterials := req.GetRawmaterials()
	associatedOwner := req.GetSupplierAssociatedOwner()

	// Process, check, compute data
	var resRawmaterials []*supplychainpb.InboundLogistics
	// for each raw material check the stock status
	for _, material := range reqRawMaterials {
		// Get the In-stock details from localDB or blockchain
		if material.GetRawmaterial().GetLogisticInStock() {
			// material is in stock change the status and send this material to manufacturer via transporter
			resRawmaterials = append(resRawmaterials, material)
		}
	}
	status := "Supplier added"

	// Send processed, computed data
	res := supplychainpb.SupplierResponse{
		//SupplierId:   supplierID,
		SupplierName:            supplierName,
		Rawmaterials:            resRawmaterials,
		SupplierAssociatedOwner: associatedOwner,
		SupplierRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		SupplierStatus:          status,
	}

	result, err := collection.InsertOne(context.TODO(), res)
	if err != nil {
		log.Fatal(err)
	}
	objID := result.InsertedID.(primitive.ObjectID)
	res.SupplierId = objID.Hex()

	// Update ObjectID/SupplierID
	supplierServer.UpdateSupplier(ctx,
		&supplychainpb.SupplierRequest{
			SupplierId:              res.SupplierId,
			SupplierName:            res.SupplierName,
			Rawmaterials:            res.Rawmaterials,
			SupplierAssociatedOwner: res.SupplierAssociatedOwner,
			SupplierRequestedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
			SupplierStatus:          res.SupplierStatus,
		})
	return &res, nil
}

func (supplierServer *server) DeleteSupplier(ctx context.Context, req *supplychainpb.SupplierRequest) (*supplychainpb.SupplierResponse, error) {

	collection := mClient.Database("scmdb").Collection("supplier")
	supplierID, err := primitive.ObjectIDFromHex(req.GetSupplierId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	_, err = collection.DeleteOne(context.TODO(), bson.M{"_id": supplierID})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &supplychainpb.SupplierResponse{
		SupplierId:     supplierID.Hex(),
		SupplierStatus: "Supplier Deleted",
	}, nil

}

func (supplierServer *server) UpdateSupplier(ctx context.Context, req *supplychainpb.SupplierRequest) (*supplychainpb.SupplierResponse, error) {

	collection := mClient.Database("scmdb").Collection("supplier")
	supplierID, err := primitive.ObjectIDFromHex(req.GetSupplierId())
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	supplierName := req.GetSupplierName()
	//supplierReqTime := req.GetSupplierRequestedTime()
	//supplierStatus := req.GetSupplierStatus()
	reqRawMaterials := req.GetRawmaterials()
	associatedOwner := req.GetSupplierAssociatedOwner()

	// Process, check, compute data
	var resRawmaterials []*supplychainpb.InboundLogistics
	// for each raw material check the stock status
	for _, material := range reqRawMaterials {
		// Get the In-stock details from localDB or blockchain
		if material.GetRawmaterial().GetLogisticInStock() {
			// material is in stock change the status and send this material to manufacturer via transporter
			resRawmaterials = append(resRawmaterials, material)
		}
	}
	status := "Supplier updated"

	// Send processed, computed data
	update := supplychainpb.SupplierResponse{
		SupplierId:              supplierID.Hex(),
		SupplierName:            supplierName,
		Rawmaterials:            resRawmaterials,
		SupplierAssociatedOwner: associatedOwner,
		SupplierRespondedTime:   time.Now().Format("01-02-2006 15:04:05 Monday"),
		SupplierStatus:          status,
	}

	filter := bson.M{"_id": supplierID}
	_ = collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))

	return &update, nil
}

func (supplierServer *server) ListAllSuppliers(req *supplychainpb.SupplierRequest, stream supplychainpb.ScmService_ListAllSuppliersServer) error {
	collection := mClient.Database("scmdb").Collection("supplier")
	data := &supplychainpb.SupplierResponse{}

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
		stream.Send(&supplychainpb.SupplierResponse{
			SupplierId:              data.SupplierId,
			SupplierName:            data.SupplierName,
			SupplierAssociatedOwner: data.SupplierAssociatedOwner,
			SupplierRespondedTime:   data.SupplierRespondedTime,
			SupplierStatus:          data.SupplierStatus,
			Rawmaterials:            data.GetRawmaterials(),
		})
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

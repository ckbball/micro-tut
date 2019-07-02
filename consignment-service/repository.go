package main

import (
  "context"
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  "go.mongodb.org/mongo-driver/mongo"
)

type repository interface {
  Create(consignment *pb.Consignment) error
  GetAll() ([]*pb.Consignment, error)
}

type MongoRepository struct {
  collection *mongo.Collection
}

// Create a new consignment
func (repository *MongoRepository) Create(consignment *pb.Consignment) error {
  _, err := repository.collection.InsertOne(context.Background(), consignment)
  return err
}

// GetAll consignments
func (repository *MongoRepository) GetAll() ([]*pb.Consignment, error) {
  cur, err := repository.collection.Find(context.Background(), nil, nil)
  var consignments []*pb.Consignment
  for cur.Next(context.Background()) {
    var consignment *pb.Consignment
    if err := cur.Decode(&consignment); err != nil {
      return nil, err
    }
    consignments = append(consignments, consignment)
  }
  return consignments, err
}

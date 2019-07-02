// consignment-service/main.go
package main

import (
  "context"
  "fmt"
  "log"
  // Import the generated protobuf code
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  vesselProto "github.com/ckbball/micro-tut/vessel-service/proto/vessel"
  "github.com/micro/go-micro"
  "os"
)

const (
  defaultHost = "datastore:27017"
)

func main() {

  // Set up micro instance
  srv := micro.NewService(
    micro.Name("consignment.service"),
  )

  srv.Init()

  uri := os.Getenv("DB_HOST")
  if uri == "" {
    uri = defaultHost
  }
  client, err := CreateClient(uri)
  if err != nil {
    log.Panic(err)
  }
  defer client.Disconnect(context.TODO())

  consignmentCollection := client.Database("shippy").Collection("consignments")

  repository := &MongoRepository{consignmentCollection}
  vesselClient := vessel.ProtoNewVesselServiceClient("shippy.service.client", srv.Client())
  h := &handler{repository, vesselClient}

  // Register handlers
  pb.RegisterShippingServiceHandler(srv.Server(), h)

  // Run the server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}

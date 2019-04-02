// consignment-service/main.go
package main

import (
  "fmt"
  "log"

  // Import the generated protobuf code
  vesselProto "github.com/EwanValentine/shippy/vessel-service/proto/vessel"
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  micro "github.com/micro/go-micro"
  "os"
)

const (
  defaultHost = "localhost:27017"
)

func main() {

  // Database host from environment variables
  host := os.Getenv("DB_HOST")

  if host == "" {
    host = defaultHost
  }

  session, err := CreateSession(host)

  // Mgo creates a "master" session, we need to end that session before
  // the main func closes.
  defer session.Close()

  if err != nil {
    // error from CreateSession
    log.Panicf("Could not connect to datastore with host %s - %v", host, err)
  }

  // Create a new service with some options
  srv := micro.NewService(

    // micro.Name must match package name given in protobuf def
    micro.Name("go.micro.srv.consignment"),
    micro.Version("latest"),
  )

  vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

  // Init will parse command line flags
  srv.Init()

  // Register handler
  pb.RegisterShippingServiceHandler(srv.Server(), &service{session, vesselClient})

  // Run the server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}

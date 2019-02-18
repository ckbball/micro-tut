// consignment-service/main.go
package main

import (
  "fmt"

  // Import the generated protobuf code
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  vesselProto "github.com/ckbball/micro-tut/vessel-service/proto/vessel"
  micro "github.com/micro/go-micro"
  "golang.org/x/net/context"
)

type Repository interface {
  Create(*pb.Consignment) (*pb.Consignment, error)
  GetAll() []*pb.Consignment
}

// Repo - Dummy repo, to simulate a datastore
// Will be replaced with real one later
type ConsignmentRepository struct {
  consignments []*pb.Consignment
}

func (repo *ConsignmentRepository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
  updated := append(repo.consignments, consignment)
  repo.consignments = updated
  return consignment, nil
}

func (repo *ConsignmentRepository) GetAll() []*pb.Consignment {
  return repo.consignments
}

// Service needs to implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface generated
// in the code for the exact method signatures
type service struct {
  repo         Repository
  vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - the only method we have on our service,
// which is create, this takes a context and a request as an
// argument, these are handled by gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

  // Here we call a client instance of our vessel service with our consignment weight,
  // and the amount of containers as the capacity value
  vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
    MaxWeight: req.Weight,
    Capacity:  int32(len(req.Containers)),
  })
  log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
  if err != nil {
    return err
  }

  // Set the VesselId as the vessel we got back from our vessel service
  req.VesselId = vesselResponse.Vessel.Id

  // Save our consignment
  consignment, err := s.repo.Create(req)
  if err != nil {
    return nil, err
  }

  // Return matching 'Response' message created in protobuf definition
  res.Created = true
  res.Consignment = consignment
  return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
  consignments := s.repo.GetAll()
  res.Consignments = consignments
  return nil
}

func main() {

  repo := &ConsignmentRepository{}

  srv := micro.NewService(

    // This name must match name given in protobuf definition
    micro.Name("go.micro.srv.consignment"),
    micro.Version("latest"),
  )

  vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

  // Init will parse command line flags
  srv.Init()

  // Register handler
  pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

  // Run the server
  if err := srv.Run(); err != nil {
    fmt.Println(err)
  }
}

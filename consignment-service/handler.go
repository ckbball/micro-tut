package main

import (
  vesselProto "github.com/EwanValentine/shippy/vessel-service/proto/vessel"
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  "golang.org/x/net/context"
  "log"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. Check generated code for the exact
// method names
type handler struct {
  vesselClient vesselProto.VesselServiceClient
}

func (s *handler) GetRepo() Repository {
  return &ConsignmentRepository{s.session.Clone()}
}

// CreateConsignment:
// Args: takes a context and a request handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
  repo := s.GetRepo()
  defer repo.Close()

  // Call a client instance of our vessel service with our consignment weight,
  // and the amount of containers as the capacity value
  vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
    MaxWeight: req.Weight,
    Capacity:  int32(len(req.Containers)),
  })
  log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
  if err != nil {
    return err
  }

  // Set VesselId as the vessel we get back from vessel service
  req.VesselId = vesselResponse.Vessel.Id

  // Save our consignment
  err = repo.Create(req)
  if err != nil {
    return err
  }

  // Return matching the 'Response' message we created in our
  // protobuf definition.
  res.Created = true
  res.Consignment = req
  return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
  repo := s.GetRepo()
  defer repo.Close()
  consignments, err := repo.GetAll()
  if err != nil {
    return err
  }

  res.Consignments = consignments
  return nil
}

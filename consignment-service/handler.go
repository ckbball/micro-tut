package main

import (
  "context"
  vesselProto "github.com/EwanValentine/shippy-service-vessel/proto/vessel"
  pb "github.com/ckbball/micro-tut/consignment-service/proto/consignment"
  "log"
)

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. Check generated code for the exact
// method names
type handler struct {
  repository
  vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment:
// Args: takes a context and a request handled by the gRPC server.
func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

  vesselResponse, err := s.vesselClient.FindAvailable(ctx, &vesselProto.Specification{
    MaxWeight: req.Weight,
    Capacity:  int32(len(req.Containers)),
  })
  log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
  if err != nil {
    return err
  }

  // Set VesselId as vessel we got from service
  req.VesselId = vessel.Response.Vessel.Id

  // Save our consignment
  if err = s.repository.Create(req); err != nil {
    return err
  }

  res.Created = true
  res.Consignment = req
  return nil

}

// GetConsignments -
func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
  consignments, err := s.repository.GetAll()
  if err != nil {
    return err
  }
  res.Consignments = consignments
  return nil
}

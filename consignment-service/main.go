package main

import (
	"context"
	"log"

	"github.com/micro/go-micro"

	pb "github.com/dime2lo/shipper-medium-tutorial/consignment-service/proto/consignment"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

type Repository struct {
	consignments []*pb.Consignment
}

func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	r.consignments = append(r.consignments, consignment)
	return consignment, nil
}

func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

type service struct {
	repo IRepository
}

func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {

	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	res.Created = true
	res.Consignment = consignment
	return nil
}

func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {

	consigments := s.repo.GetAll()
	res.Consignments = consigments
	return nil
}

func main() {
	log.Println("Consignment service starting...")
	repo := &Repository{}

	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	srv.Init()

	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo})

	if err := srv.Run(); err != nil {
		log.Println(err)
	}
}

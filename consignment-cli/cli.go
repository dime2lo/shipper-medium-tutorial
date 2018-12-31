package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	pb "github.com/dime2lo/shipper-medium-tutorial/consignment-service/proto/consignment"
	microclient "github.com/micro/go-micro/client"
	"github.com/micro/go-micro/cmd"
	"github.com/pkg/errors"
)

const (
	address         = "localhost:50051"
	defaultFileName = "consignment.json"
)

func parseFile(file string) (*pb.Consignment, error) {
	var consignment *pb.Consignment
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, errors.Wrap(err, "error reading file:"+file)
	}

	err = json.Unmarshal(data, &consignment)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode json data to consignment")
	}
	return consignment, nil
}

func main() {
	cmd.Init()

	client := pb.NewShippingServiceClient("go.micro.srv.consignment", microclient.DefaultClient)
	createConsignment(client)

	getAll, err := client.GetConsignments(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("error getting all consignments: %v", err)
	}

	for _, c := range getAll.GetConsignments() {
		log.Println(c)
	}
}

func createConsignment(client pb.ShippingServiceClient) {
	file := defaultFileName
	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not parse file: %v", err)
	}

	r, err := client.CreateConsignment(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Failed to create Consignment: %v", err)
	}
	log.Printf("Created: %t", r.Created)
}

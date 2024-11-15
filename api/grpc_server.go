package api

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	pb "github.com/Dendyator/AntiBF/api/proto/pb" //nolint
	"github.com/Dendyator/AntiBF/internal/core"   //nolint
	"github.com/Dendyator/AntiBF/internal/logger" //nolint
)

type server struct {
	pb.UnimplementedAntiBruteForceServer
	logger *logger.Logger
}

func NewServer(log *logger.Logger) *server {
	return &server{logger: log}
}

func (s *server) CheckAuthorization(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debugf("Received CheckAuthorization request for login: %s, IP: %s", req.GetLogin(), req.GetIp())

	if !isValidCIDR(req.GetIp()) {
		s.logger.Warnf("Invalid IP format: %s", req.GetIp())
		return &pb.AuthResponse{Ok: false}, nil
	}

	ok := core.CheckAuthorization(req.GetLogin(), req.GetPassword(), req.GetIp())
	if !ok {
		s.logger.Warnf("Authorization blocked for login: %s, IP: %s", req.GetLogin(), req.GetIp())
	}
	return &pb.AuthResponse{Ok: ok}, nil
}

func (s *server) ResetBucket(ctx context.Context, req *pb.ResetRequest) (*pb.ResetResponse, error) {
	s.logger.Debugf("Received ResetBucket request for login: %s, IP: %s", req.GetLogin(), req.GetIp())

	if !isValidCIDR(req.GetIp()) {
		s.logger.Warnf("Invalid IP format: %s", req.GetIp())
		return &pb.ResetResponse{Success: false}, nil
	}

	success := core.ResetBucket(req.GetLogin(), req.GetIp())
	if !success {
		s.logger.Warnf("Reset bucket failed for login: %s, IP: %s", req.GetLogin(), req.GetIp())
	} else {
		s.logger.Infof("Bucket successfully reset for login: %s, IP: %s", req.GetLogin(), req.GetIp())
	}
	return &pb.ResetResponse{Success: success}, nil
}

func (s *server) AddToBlacklist(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debugf("Received AddToBlacklist request for subnet: %s", req.GetSubnet())

	if !isValidCIDR(req.GetSubnet()) {
		s.logger.Warnf("Invalid subnet format: %s", req.GetSubnet())
		return &pb.ListResponse{Success: false}, nil
	}

	success := core.ManageList(req.GetSubnet(), core.Blacklist, true)
	if !success {
		s.logger.Warnf("Failed to add to blacklist: %s", req.GetSubnet())
	} else {
		s.logger.Infof("Successfully added to blacklist: %s", req.GetSubnet())
	}
	return &pb.ListResponse{Success: success}, nil
}

func (s *server) RemoveFromBlacklist(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debugf("Received RemoveFromBlacklist request for subnet: %s", req.GetSubnet())

	if !isValidCIDR(req.GetSubnet()) {
		s.logger.Warnf("Invalid subnet format: %s", req.GetSubnet())
		return &pb.ListResponse{Success: false}, nil
	}

	success := core.ManageList(req.GetSubnet(), core.Blacklist, false)
	if !success {
		s.logger.Warnf("Failed to remove from blacklist: %s", req.GetSubnet())
	} else {
		s.logger.Infof("Successfully removed from blacklist: %s", req.GetSubnet())
	}
	return &pb.ListResponse{Success: success}, nil
}

func (s *server) AddToWhitelist(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debugf("Received AddToWhitelist request for subnet: %s", req.GetSubnet())

	if !isValidCIDR(req.GetSubnet()) {
		s.logger.Warnf("Invalid subnet format: %s", req.GetSubnet())
		return &pb.ListResponse{Success: false}, nil
	}

	success := core.ManageList(req.GetSubnet(), core.Whitelist, true)
	if !success {
		s.logger.Warnf("Failed to add to whitelist: %s", req.GetSubnet())
	} else {
		s.logger.Infof("Successfully added to whitelist: %s", req.GetSubnet())
	}
	return &pb.ListResponse{Success: success}, nil
}

func (s *server) RemoveFromWhitelist(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.logger.Debugf("Received RemoveFromWhitelist request for subnet: %s", req.GetSubnet())

	if !isValidCIDR(req.GetSubnet()) {
		s.logger.Warnf("Invalid subnet format: %s", req.GetSubnet())
		return &pb.ListResponse{Success: false}, nil
	}

	success := core.ManageList(req.GetSubnet(), core.Whitelist, false)
	if !success {
		s.logger.Warnf("Failed to remove from whitelist: %s", req.GetSubnet())
	} else {
		s.logger.Infof("Successfully removed from whitelist: %s", req.GetSubnet())
	}
	return &pb.ListResponse{Success: success}, nil
}

func RunGRPCServer(log *logger.Logger) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := NewServer(log)
	pb.RegisterAntiBruteForceServer(grpcServer, srv)
	reflection.Register(grpcServer)

	log.Infof("gRPC server is listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func isValidCIDR(cidr string) bool {
	_, _, err := net.ParseCIDR(cidr)
	return err == nil
}

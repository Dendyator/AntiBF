package grpc

import (
	"context"
	"github.com/Dendyator/AntiBF/internal/delivery/grpc/proto/pb"
	"github.com/Dendyator/AntiBF/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"

	"github.com/Dendyator/AntiBF/internal/usecase" //nolint
)

type server struct {
	pb.UnimplementedAntiBruteForceServer
	rateLimiter *usecase.RateLimiter
	logger      *logger.Logger
}

func NewServer(rateLimiter *usecase.RateLimiter, log *logger.Logger) *server {
	return &server{
		rateLimiter: rateLimiter,
		logger:      log,
	}
}

func (s *server) CheckAuthorization(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	s.logger.Debugf("Received CheckAuthorization request for login: %s, IP: %s", req.GetLogin(), req.GetIp())

	if !isValidCIDR(req.GetIp()) {
		s.logger.Warnf("Invalid IP format: %s", req.GetIp())
		return &pb.AuthResponse{Ok: false}, nil
	}

	ok := s.rateLimiter.CheckAuthorization(req.GetLogin(), req.GetPassword(), req.GetIp())
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

	success := s.rateLimiter.ResetBucket(req.GetLogin(), req.GetIp())
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

	success := s.rateLimiter.ManageList(req.GetSubnet(), usecase.Blacklist, true)
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

	success := s.rateLimiter.ManageList(req.GetSubnet(), usecase.Blacklist, false)
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

	success := s.rateLimiter.ManageList(req.GetSubnet(), usecase.Whitelist, true)
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

	success := s.rateLimiter.ManageList(req.GetSubnet(), usecase.Whitelist, false)
	if !success {
		s.logger.Warnf("Failed to remove from whitelist: %s", req.GetSubnet())
	} else {
		s.logger.Infof("Successfully removed from whitelist: %s", req.GetSubnet())
	}
	return &pb.ListResponse{Success: success}, nil
}

func RunGRPCServer(rateLimiter *usecase.RateLimiter, log *logger.Logger) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	srv := NewServer(rateLimiter, log)
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

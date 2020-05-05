package main

import (
	"context"
	"net"

	"github.com/noptics/golog"
	"github.com/noptics/services/cmd/registry/data"
	"github.com/noptics/services/pkg/nproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	db      data.Store
	grpcs   *grpc.Server
	errChan chan error
}

func NewGRPCServer(db data.Store, port string, errChan chan error) (*GRPCServer, error) {
	gs := &GRPCServer{
		db:      db,
		errChan: errChan,
		grpcs:   grpc.NewServer(),
	}

	nproto.RegisterProtoRegistryServer(gs.grpcs, gs)

	golog.Infow("starting grpc service", "port", port)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return nil, err
	}

	go func() {
		if err = gs.grpcs.Serve(lis); err != nil {
			gs.errChan <- err
		}
	}()

	return gs, nil
}

func (s *GRPCServer) Stop() {
	s.grpcs.GracefulStop()
}

func (s *GRPCServer) GetFiles(ctx context.Context, in *nproto.GetFilesRequest) (*nproto.GetFilesReply, error) {
	_, f, err := s.db.GetChannelData(in.ClusterID, in.Channel)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetFilesReply{Files: f}, nil
}

func (s *GRPCServer) SaveFiles(ctx context.Context, in *nproto.SaveFilesRequest) (*nproto.SaveFilesReply, error) {
	err := s.db.SaveFiles(in.ClusterID, in.Channel, in.Files)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.SaveFilesReply{}, nil
}

func (s *GRPCServer) SetMessage(ctx context.Context, in *nproto.SetMessageRequest) (*nproto.SetMessageReply, error) {
	err := s.db.SetChannelMessage(in.ClusterID, in.Channel, in.Name)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.SetMessageReply{}, nil
}

func (s *GRPCServer) GetMessage(ctx context.Context, in *nproto.GetMessageRequest) (*nproto.GetMessageReply, error) {
	m, _, err := s.db.GetChannelData(in.ClusterID, in.Channel)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetMessageReply{Name: m}, nil
}

func (s *GRPCServer) GetChannelData(ctx context.Context, in *nproto.GetChannelDataRequest) (*nproto.GetChannelDataReply, error) {
	m, files, err := s.db.GetChannelData(in.ClusterID, in.Channel)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetChannelDataReply{ClusterID: in.ClusterID, Channel: in.Channel, Files: files, Message: m}, nil
}

func (s *GRPCServer) SetChannelData(ctx context.Context, in *nproto.SetChannelDataRequest) (*nproto.SetChannelDataReply, error) {
	err := s.db.SaveChannelData(in.ClusterID, in.Channel, in.Message, in.Files)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.SetChannelDataReply{}, nil
}

func (s *GRPCServer) GetChannels(ctx context.Context, in *nproto.GetChannelsRequest) (*nproto.GetChannelsReply, error) {
	channels, err := s.db.GetChannels(in.ClusterID)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetChannelsReply{Channels: channels}, nil
}

func (s *GRPCServer) SaveCluster(ctx context.Context, in *nproto.SaveClusterRequest) (*nproto.SaveClusterReply, error) {
	id, err := s.db.SaveCluster(in.Cluster)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.SaveClusterReply{Id: id}, nil
}

func (s *GRPCServer) GetCluster(ctx context.Context, in *nproto.GetClusterRequest) (*nproto.GetClusterReply, error) {
	cluster, err := s.db.GetCluster(in.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetClusterReply{Cluster: cluster}, nil
}

func (s *GRPCServer) GetClusters(ctx context.Context, in *nproto.GetClustersRequest) (*nproto.GetClustersReply, error) {
	clusters, err := s.db.GetClusters()

	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return &nproto.GetClustersReply{Clusters: clusters}, nil
}

package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/nebula-vpn/backend/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	server *grpc.Server
	cfg    *config.ServerConfig
	logger *zap.Logger
}

func NewServer(cfg *config.ServerConfig, logger *zap.Logger) *Server {
	s := &Server{
		server: grpc.NewServer(
			grpc.UnaryInterceptor(LoggingInterceptor(logger)),
			grpc.MaxRecvMsgSize(10 * 1024 * 1024),
			grpc.MaxSendMsgSize(10 * 1024 * 1024),
		),
		cfg:    cfg,
		logger: logger,
	}

	reflection.Register(s.server)
	return s
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC port: %w", err)
	}

	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
	}()

	s.logger.Info("gRPC server starting", zap.Int("port", s.cfg.GRPCPort))
	return s.server.Serve(lis)
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.server.RegisterService(desc, impl)
}

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Debug("gRPC call", zap.String("method", info.FullMethod))
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("gRPC error", zap.String("method", info.FullMethod), zap.Error(err))
		}
		return resp, err
	}
}

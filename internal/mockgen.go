package internal

//go:generate mockgen -source=services/grpc/goph_keeper/v1/goph_keeper_grpc.pb.go -destination=mocks/goph_keeper_grpc_mock.go -package=mocks
//go:generate mockgen -source=services/entities/file_helper.go -destination=mocks/file_helper_mock.go -package=mocks
//go:generate mockgen -source=server/grpc_server.go -destination=mocks/grpc_server_mock.go -package=mocks



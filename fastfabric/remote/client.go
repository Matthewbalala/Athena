package remote

import (
	"google.golang.org/grpc"
)

var storageClient StoragePeerClient

func StartStoragePeerClient(address string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return err
	}
	storageClient = NewStoragePeerClient(conn)
	return nil
}

func GetStoragePeerClient() StoragePeerClient {
	return storageClient
}

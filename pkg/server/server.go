package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aramase/azure-appconfig-csi-provider/pkg/appconfig"
	"github.com/aramase/azure-appconfig-csi-provider/pkg/types"

	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sigs.k8s.io/secrets-store-csi-driver/provider/v1alpha1"
)

const (
	connectionStringKey = "connectionstring"
	kvsStringKey        = "kvs"
)

type Server struct {
	log logr.Logger
}

var _ v1alpha1.CSIDriverProviderServer = &Server{}

// Mount implements the provider gRPC method
func (s *Server) Mount(ctx context.Context, req *v1alpha1.MountRequest) (*v1alpha1.MountResponse, error) {
	var attrib, secrets map[string]string
	var filePermission os.FileMode
	var err error

	// only connection string is supported for now
	// TODO(aramase): support fetching using managed identity
	if req.GetSecrets() == "" {
		return nil, status.Error(codes.InvalidArgument, "secrets cannot be empty")
	}

	err = json.Unmarshal([]byte(req.GetAttributes()), &attrib)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid attributes")
	}
	err = json.Unmarshal([]byte(req.GetSecrets()), &secrets)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid secrets")
	}
	err = json.Unmarshal([]byte(req.GetPermission()), &filePermission)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid permission")
	}

	client := appconfig.New(secrets[connectionStringKey])
	kvList, err := parseMountKVs(attrib[kvsStringKey])
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid kvs")
	}

	out := &v1alpha1.MountResponse{}
	ovs := []*v1alpha1.ObjectVersion{}

	for _, kv := range kvList {
		result, err := client.GetKV(kv.Key, kv.Label)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("failed to get key %s: %v", kv.Key, err))
		}

		// a single key with no label can return multiple results if the key is a prefix
		for _, res := range result {
			path := res.Key
			if kv.Label != "" {
				path = fmt.Sprintf("%s/%s", path, kv.Label)
			}
			out.Files = append(out.Files, &v1alpha1.File{
				Path:     path,
				Mode:     int32(filePermission),
				Contents: []byte(res.Value),
			})
			s.log.Info("added kv to response", "key", kv.Key, "label", kv.Label, "path", path)
			// using the etag as the object version
			ovs = append(ovs, &v1alpha1.ObjectVersion{Id: path, Version: res.ETag})
		}
	}

	out.ObjectVersion = ovs
	return out, nil
}

// Version implements the provider gRPC method
func (s *Server) Version(ctx context.Context, req *v1alpha1.VersionRequest) (*v1alpha1.VersionResponse, error) {
	return nil, status.Error(codes.Unimplemented, "version not implemented")
}

func parseMountKVs(kvs string) ([]types.KV, error) {
	var kvList []types.KV
	err := json.Unmarshal([]byte(kvs), &kvList)
	if err != nil {
		return nil, fmt.Errorf("failed to parse kvs: %v", err)
	}
	return kvList, nil
}

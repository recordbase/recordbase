/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package recordbase

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sprintframework/raftpb"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/emptypb"
	"strings"
     _ "github.com/codeallergy/grpc-multi-resolver"
)

func NewClient(ctx context.Context, commaSeparatedEndpoints, token string, tlsConfigOpt *tls.Config) (Client, error) {

	endpoints := splitAndTrim(commaSeparatedEndpoints, ",")

	conf, err := findClusterConfiguration(ctx, endpoints, token, tlsConfigOpt)
	if err != nil {
		return nil, err
	}

	fmt.Printf("conf = %s\n", conf.String())

	if len(conf.ServerList) == 0 {
		return nil, errors.New("no raft servers found")
	}

	switch len(conf.ServerList) {

	case 0:
		return nil, errors.New("no raft servers found")

	case 1:

		conn, err := Dial(ctx, conf.ServerList[0].ApiAddr, token, tlsConfigOpt)
		if err != nil {
			return nil, err
		}

		return Create(conn), nil

	default:
		multipointEndpoint := formatMultipointEndpoint(conf)

		conn, err := DialWithLoadBalancer(ctx, multipointEndpoint, ServiceName, token, tlsConfigOpt)
		if err != nil {
			return nil, err
		}

		return Create(conn), nil
	}

}

func findClusterConfiguration(ctx context.Context, endpoints []string, token string, tlsConfigOpt *tls.Config) (conf *raftpb.RaftConfiguration, err error) {

	for _, endpoint := range endpoints {

		conf, err = requestClusterConfiguration(ctx, endpoint, token, tlsConfigOpt)
		if err == nil {
			break
		}

	}
	return
}

func requestClusterConfiguration(ctx context.Context, publicEndpoint, token string, tlsConfigOpt *tls.Config) (*raftpb.RaftConfiguration, error) {

	conn, err := Dial(ctx, publicEndpoint, token, tlsConfigOpt)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	raftService := raftpb.NewRaftServiceClient(conn)
	return raftService.GetConfiguration(ctx, &emptypb.Empty{})
}


func splitAndTrim(str string, sep string) []string {
	var arr []string
	for _, s := range strings.Split(str, sep) {
		s = strings.TrimSpace(s)
		if len(s) > 0 {
			arr = append(arr, s)
		}
	}
	return arr
}

func formatMultipointEndpoint(conf *raftpb.RaftConfiguration) string {
	var multipoint strings.Builder
	multipoint.WriteString("multi:///")
	for i, server := range conf.ServerList {
		if i > 0 {
			multipoint.WriteRune(',')
		}
		multipoint.WriteString(server.ApiAddr)
	}
	return multipoint.String()
}
/*
 * Copyright (c) 2023 Zander Schwid & Co. LLC.
 * SPDX-License-Identifier: BUSL-1.1
 */

package recordbase

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

var ClientMaxReceiveMessageSize = 1024 * 1024 * 100

func Dial(ctx context.Context, endpoint, token string, tlsConfigOpt *tls.Config) (*grpc.ClientConn, error) {

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	fmt.Printf("endpoint = %s\n", endpoint)

	return grpc.DialContext(ctx, endpoint,
		//grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithPerRPCCredentials(TokenAuth(token)),
		grpc.WithTransportCredentials(getCredentials(tlsConfigOpt)))

}

func DialWithLoadBalancer(ctx context.Context, multipointEndpoint, serviceName, token string, tlsConfigOpt *tls.Config) (*grpc.ClientConn, error) {

	// https://github.com/grpc/grpc/blob/master/doc/service_config.md
	serviceConfig := `{"healthCheckConfig": {"serviceName": "%s"}, "loadBalancingConfig": [ { "round_robin": {} } ]}`

	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(100 * time.Millisecond)),
		grpc_retry.WithMax(5),
	}

	return grpc.DialContext(ctx, multipointEndpoint,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(serviceConfig, serviceName)),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(ClientMaxReceiveMessageSize)),
		grpc.WithDefaultCallOptions(grpc.WaitForReady(true)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithPerRPCCredentials(TokenAuth(token)),
		grpc.WithTransportCredentials(getCredentials(tlsConfigOpt)))

}

func getCredentials(tlsConfigOptional *tls.Config) credentials.TransportCredentials {
	if tlsConfigOptional != nil {
		return credentials.NewTLS(tlsConfigOptional)
	} else {
		return insecure.NewCredentials()
	}
}
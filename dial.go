/*
 * Copyright (c) 2022-2023 Zander Schwid & Co. LLC.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
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

	return grpc.DialContext(ctx, endpoint,
		grpc.WithBlock(),
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
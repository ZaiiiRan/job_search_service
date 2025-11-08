package client

import (
	"context"

	"github.com/ZaiiiRan/job_search_service/common/pkg/ctxmetadata"
	"google.golang.org/grpc"
)

func PropagateClientMetaUnary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		ctx = ctxmetadata.ForwardReqIdToOutgoingContext(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

package api

import (
	"context"
	"errors"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/gempellm/golang-chan/internal/model"
	"github.com/gempellm/golang-chan/internal/repo"

	pb "github.com/gempellm/golang-chan/pkg/go_chan_api"
)

var (
	totalPostNotFound = promauto.NewCounter(prometheus.CounterOpts{
		Name: "logistic_Post_api_Post_not_found_total",
		Help: "Total number of Posts that were not found",
	})
)

type PostAPI struct {
	pb.UnimplementedGoChanApiServiceServer
	repo repo.Repo
}

// NewPostAPI returns api of go-chan-api service
func NewPostAPI(r repo.Repo) pb.GoChanApiServiceServer {
	return &PostAPI{repo: r}
}

func (o *PostAPI) DescribePostV1(
	ctx context.Context,
	req *pb.DescribePostV1Request,
) (*pb.DescribePostV1Response, error) {

	// if err := req.Validate(); err != nil {
	// 	log.Error().Err(err).Msg("DescribePostV1 - invalid argument")

	// 	return nil, status.Error(codes.InvalidArgument, err.Error())
	// }

	Post, err := o.repo.DescribePost(ctx, req.PostId)
	if err != nil {
		log.Error().Err(err).Msg("DescribePostV1 -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if Post == nil {
		log.Debug().Uint64("PostId", req.PostId).Msg("Post not found")
		totalPostNotFound.Inc()

		return nil, status.Error(codes.NotFound, "Post not found")
	}

	log.Debug().Msg("DescribePostV1 - success")

	return &pb.DescribePostV1Response{
		Post: &pb.Post{
			Id: Post.ID,
		},
	}, nil
}

func (o *PostAPI) CreatePost(ctx context.Context, req *pb.CreatePostV1Request) (*pb.CreatePostV1Response, error) {

	mockMedia := make([]*model.Media, 0)

	_, err := o.repo.CreatePost(ctx, req.Title, req.Msg, req.ThreadId, mockMedia)
	if err != nil {
		log.Error().Err(err).Msg("CreatePost -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	log.Debug().Msg("CreatePost - success")

	return &pb.CreatePostV1Response{
		PostId: 42,
	}, nil
}

func (o *PostAPI) DescribePost(ctx context.Context, req *pb.DescribePostV1Request) (*pb.DescribePostV1Response, error) {
	// if err := req.Validate(); err != nil {
	// 	log.Error().Err(err).Msg("DescribePost - invalid argument")

	// 	return nil, status.Error(codes.InvalidArgument, err.Error())
	// }

	post, err := o.repo.DescribePost(ctx, req.PostId)
	if err != nil {
		if errors.Is(err, repo.ErrPostNotFound) {
			return nil, status.Error(codes.NotFound, "Post not found")
		}

		log.Error().Err(err).Msg("DescribePost -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	if post == nil {
		log.Debug().Uint64("PostId", req.PostId).Msg("Post not found")
		totalPostNotFound.Inc()

		return nil, status.Error(codes.NotFound, "Post not found")
	}

	log.Debug().Msg("DescribePost - success")

	return &pb.DescribePostV1Response{
		Post: &pb.Post{
			Id:       post.ID,
			ThreadId: post.ThreadID,
			Title:    post.Title,
			Msg:      post.Message,
			Created:  timestamppb.Now(),
		},
	}, nil
}

func (o *PostAPI) ListPosts(ctx context.Context, req *pb.ListPostsV1Request) (*pb.ListPostsV1Response, error) {
	// if err := req.Validate(); err != nil {
	// 	log.Error().Err(err).Msg("ListPosts - invalid argument")

	// 	return nil, status.Error(codes.InvalidArgument, err.Error())
	// }

	PostsFromRepo, err := o.repo.ListPosts(ctx, req.ThreadId)
	if err != nil {
		log.Error().Err(err).Msg("ListPosts -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	PostsToPb := []*pb.Post{}

	for _, v := range PostsFromRepo {
		p := &pb.Post{Id: v.ID, ThreadId: v.ThreadID, Title: v.Title, Msg: v.Message, Created: v.Created}
		PostsToPb = append(PostsToPb, p)
	}

	return &pb.ListPostsV1Response{
		Posts: PostsToPb,
	}, nil
}

func (o *PostAPI) RemovePost(ctx context.Context, req *pb.RemovePostV1Request) (*pb.RemovePostV1Response, error) {
	// if err := req.Validate(); err != nil {
	// 	log.Error().Err(err).Msg("RemovePost - invalid argument")

	// 	return nil, status.Error(codes.InvalidArgument, err.Error())
	// }

	ok, err := o.repo.RemovePost(ctx, req.PostId)
	if err != nil {
		log.Error().Err(err).Msg("RemovePost -- failed")

		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemovePostV1Response{
		Found: ok,
	}, nil
}

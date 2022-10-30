package repo

import (
	"context"
	"errors"

	"github.com/gempellm/golang-chan/internal/model"
	"github.com/jmoiron/sqlx"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ErrPostNotFound error = errors.New("post not found")

var Posts []*model.Post = []*model.Post{
	{ID: 1, Message: "first Post"},
	{ID: 2, Message: "second Post"},
	{ID: 3, Message: "third Post"},
	{ID: 4, Message: "fourth Post"},
	{ID: 5, Message: "fifth Post"},
}

//go:generate mockgen -package=mocks -destination=./internal/mocks/repo_mock.go -self_package=github.com/gempellm/golang-chan/internal/repo . Repo

// Repo is DAO for Post
type Repo interface {
	CreatePost(ctx context.Context, title, message string, threadID uint64, attachedMedia []*model.Media) (*model.Post, error)
	DescribePost(ctx context.Context, postID uint64) (*model.Post, error)
	ListPosts(ctx context.Context, threadID uint64) ([]*model.Post, error)
	RemovePost(ctx context.Context, postID uint64) (bool, error)
}

type repo struct {
	db        *sqlx.DB
	batchSize uint
}

// NewRepo returns Repo interface
func NewRepo(db *sqlx.DB, batchSize uint) Repo {
	return &repo{db: db, batchSize: batchSize}
}

func (r *repo) DescribePost(ctx context.Context, PostID uint64) (*model.Post, error) {
	for _, v := range Posts {
		if v.ID == PostID {
			return v, nil
		}
	}
	return nil, ErrPostNotFound
}

func (r *repo) CreatePost(ctx context.Context, title, message string, threadID uint64, attachedMedia []*model.Media) (*model.Post, error) {
	var p *model.Post

	if len(Posts) == 0 {
		p = &model.Post{ID: 0, Message: message, Created: timestamppb.Now()}
	} else {
		p = &model.Post{ID: (Posts[len(Posts)-1].ID) + 1, Message: message, Created: timestamppb.Now()}
	}

	Posts = append(Posts, p)

	return p, nil
}

func (r *repo) ListPosts(ctx context.Context, threadID uint64) ([]*model.Post, error) {
	PostsLength := uint64(len(Posts))

	if PostsLength == 0 {
		return nil, nil
	}

	// lastID := uint64(PostsLength - 1)
	// if cursor > lastID {
	// 	return nil, nil
	// }

	var result []*model.Post = make([]*model.Post, 0)

	// for i := cursor; ; i++ {
	// 	if i > lastID {
	// 		break
	// 	}

	// 	result = append(result, Posts[i])
	// }

	return result, nil
}

func (r *repo) RemovePost(ctx context.Context, postID uint64) (bool, error) {
	var ok bool

	for i, v := range Posts {
		if v.ID == postID {
			newPosts := make([]*model.Post, 0)
			newPosts = append(newPosts, Posts[:i]...)
			newPosts = append(newPosts, Posts[i+1:]...)

			Posts = newPosts
			ok = true
			break
		}
	}

	return ok, nil
}

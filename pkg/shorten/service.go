package shorten

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/aladhims/shortener/pkg/base"
	pb "github.com/aladhims/shortener/pkg/shorten/proto"
)

var slugAlreadyExists = errors.New("The slug is already exist")
var requestedSlugDoesNotExist = errors.New("The requested slug doesn't exist")

type service struct {
	repository Repository
}

func NewService(r Repository) pb.ServiceServer {
	return &service{
		repository: r,
	}
}

func (s *service) Shorten(ctx context.Context, req *pb.ShortURL) (*pb.ShortenResponse, error) {
	rand.Seed(time.Now().UnixNano())
	res := &pb.ShortenResponse{}
	existingShortURL := &pb.ShortURL{}
	if req.UrlType == pb.URLType_DEFINED {
		existingShortURL = s.repository.GetByOrigin(req.Origin)

		if isSame(req, existingShortURL) {
			res.Slug = existingShortURL.Slug
			res.Status = pb.ShortenResponseStatus_SAME_ORIGIN
			return res, nil
		}

		existingShortURL = s.repository.GetBySlug(req.Slug)
		if existingShortURL != nil {
			return &pb.ShortenResponse{
				Status: pb.ShortenResponseStatus_SLUG_ALREADY_EXISTS,
			}, nil
		}

		req.Base = base.Decode(req.Slug)

		err := s.repository.Create(req)
		if err != nil {
			return &pb.ShortenResponse{
				Status: pb.ShortenResponseStatus_FAILED_SHORTEN,
			}, nil
		}

		res.Slug = req.Slug
		res.Status = pb.ShortenResponseStatus_SUCCESS_SHORTEN

		return res, nil
	}

	existingShortURL = s.repository.GetByOrigin(req.Origin)

	if isSame(req, existingShortURL) {
		// there's an existing url uses a user defined slug
		res.Slug = existingShortURL.Slug
		res.Status = pb.ShortenResponseStatus_SAME_ORIGIN
		return res, nil
	}

	for {
		req.Base = rand.Uint64()
		req.Slug = base.Encode(req.Base)
		existingShortURL = s.repository.GetBySlug(req.Slug)
		if existingShortURL == nil {
			break
		}
	}

	err := s.repository.Create(req)
	if err != nil {
		return &pb.ShortenResponse{
			Status: pb.ShortenResponseStatus_FAILED_SHORTEN,
		}, nil
	}

	res.Slug = req.Slug

	return res, nil
}

func (s *service) Expand(ctx context.Context, req *pb.ExpandRequest) (*pb.ExpandResponse, error) {
	origin := base.Decode(req.Slug)

	shortURL := s.repository.GetByBase(origin)
	if shortURL == nil {
		return &pb.ExpandResponse{
			Status: pb.ExpandResponseStatus_NOT_FOUND,
		}, nil
	}

	return &pb.ExpandResponse{
		ShortURL: shortURL,
		Status:   pb.ExpandResponseStatus_SUCCESS_EXPAND,
	}, nil
}

func (s *service) Check(ctx context.Context, req *pb.HealthCheckRequest) (*pb.HealthCheckResponse, error) {
	if s.repository == nil {
		return &pb.HealthCheckResponse{
			Status: pb.HealthCheckResponse_NOT_SERVING,
		}, nil
	}

	return &pb.HealthCheckResponse{
		Status: pb.HealthCheckResponse_SERVING,
	}, nil
}

func isSame(a, b *pb.ShortURL) bool {
	if b == nil {
		return false
	}
	if a.Origin == b.Origin || a.Slug == b.Slug {
		return true
	}

	return false
}

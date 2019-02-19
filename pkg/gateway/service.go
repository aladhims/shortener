package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	notificationpb "github.com/aladhims/shortener/pkg/notification/proto"
	shortenpb "github.com/aladhims/shortener/pkg/shorten/proto"
	userpb "github.com/aladhims/shortener/pkg/user/proto"
	"github.com/gorilla/mux"
)

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILED  = "failed"
	API_VERSION    = "v1"
)

type response struct {
	Data    interface{} `json:"data,omitempty"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
}

type ApiGateway struct {
	ctx                context.Context
	shortenClient      shortenpb.ServiceClient
	userClient         userpb.ServiceClient
	notificationClient notificationpb.ServiceClient
}

func NewService(sc shortenpb.ServiceClient, uc userpb.ServiceClient, nc notificationpb.ServiceClient) ApiGateway {
	return ApiGateway{
		ctx:                context.Background(),
		shortenClient:      sc,
		userClient:         uc,
		notificationClient: nc,
	}
}

func (s *ApiGateway) handleShorten(w http.ResponseWriter, r *http.Request) {
	var shortURLReq shortenpb.ShortURL

	err := r.ParseForm()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fullname := r.FormValue("fullname")
	email := r.FormValue("email")

	userRes, err := s.userClient.Create(s.ctx, &userpb.User{
		Fullname: fullname,
		Email:    email,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	origin := r.FormValue("origin")

	shortURLReq.Origin = origin
	shortURLReq.UrlType = shortenpb.URLType_RANDOM
	shortURLReq.UserId = userRes.Id
	slug := r.FormValue("slug")

	if slug != "" {
		shortURLReq.UrlType = shortenpb.URLType_DEFINED
		shortURLReq.Slug = slug
	}

	shortenRes, err := s.shortenClient.Shorten(s.ctx, &shortURLReq)

	encoder := json.NewEncoder(w)

	switch shortenRes.Status {
	case shortenpb.ShortenResponseStatus_SUCCESS_SHORTEN:
		encoder.Encode(response{
			Data:    shortenRes,
			Status:  STATUS_SUCCESS,
			Message: "success shortened your url",
		})
		_, err = s.notificationClient.Notify(s.ctx, &notificationpb.NotifyRequest{
			Email:    email,
			Fullname: fullname,
			Origin:   origin,
			Slug:     shortenRes.Slug,
		})
		break
	case shortenpb.ShortenResponseStatus_FAILED_SHORTEN:
		http.Error(w, "Something bad happened", http.StatusInternalServerError)
		break
	case shortenpb.ShortenResponseStatus_SLUG_ALREADY_EXISTS:
		encoder.Encode(response{
			Status:  STATUS_FAILED,
			Message: "Can't use the defined slug",
		})
		break
	case shortenpb.ShortenResponseStatus_SAME_ORIGIN:
		encoder.Encode(response{
			Data:    shortenRes,
			Status:  STATUS_SUCCESS,
			Message: "this origin url has already been registered",
		})
		break
	default:
		http.Error(w, "Something bad happened", http.StatusInternalServerError)
	}
}

func (s *ApiGateway) handleExpand(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	slug := params["slug"]

	res, _ := s.shortenClient.Expand(s.ctx, &shortenpb.ExpandRequest{
		Slug: slug,
	})

	switch res.Status {
	case shortenpb.ExpandResponseStatus_SUCCESS_EXPAND:
		http.Redirect(w, r, res.ShortURL.Origin, http.StatusMovedPermanently)
		break
	case shortenpb.ExpandResponseStatus_NOT_FOUND:
		http.NotFound(w, r)
		break
	default:
		http.Error(w, "something bad happened", http.StatusInternalServerError)
		break
	}
}

func (s *ApiGateway) handleHealth(w http.ResponseWriter, r *http.Request) {
	if s.shortenClient == nil || s.userClient == nil || s.notificationClient == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "NOT_HEALTHY")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (s *ApiGateway) Run(port string) {
	router := mux.NewRouter()

	router.HandleFunc("/api/"+API_VERSION+"/shorten", s.handleShorten).Methods("POST")
	router.HandleFunc("/api/"+API_VERSION+"/health", s.handleHealth).Methods("GET")
	router.HandleFunc("/{slug}", s.handleExpand).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))
}

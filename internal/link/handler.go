package link

import (
	"fmt"
	"go/adv-demo/configs"
	"go/adv-demo/pkg/event"
	"go/adv-demo/pkg/middleware"
	"go/adv-demo/pkg/request"
	"go/adv-demo/pkg/response"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
	Config         *configs.Config
	EventBus       *event.EventBus
}

type LinkHandler struct {
	LinkRepository *LinkRepository
	EventBus       *event.EventBus
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		EventBus:       deps.EventBus,
	}
	router.HandleFunc("POST /link", handler.Create())
	router.Handle("PATCH /link/{id}", middleware.IsAuthed(handler.Update(), deps.Config))
	router.HandleFunc("DELETE /link/{id}", handler.Delete())
	router.HandleFunc("GET /{hash}", handler.GoTo())
	router.HandleFunc("GET /link", handler.GetAll())
}

func (handler *LinkHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := request.HandleBody[LinkCreateRequest](&w, req)
		if err != nil {
			return
		}
		link := NewLink(body.Url)
		for {
			existedLink, _ := handler.LinkRepository.GetByHash(link.Hash)
			if existedLink == nil {
				break
			}
			link.GenerateHash()
		}
		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Response(w, createdLink, http.StatusCreated)
	}
}

func (handler *LinkHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		email, ok := req.Context().Value(middleware.ContextEmailKey).(string)
		if ok {
			fmt.Println(email)
		}
		body, err := request.HandleBody[LinkUpdateRequest](&w, req)
		if err != nil {
			return
		}
		idStr := req.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   body.Url,
			Hash:  body.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		response.Response(w, link, http.StatusCreated)
	}
}

func (handler *LinkHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idStr := req.PathValue("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = handler.LinkRepository.GetById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		response.Response(w, nil, http.StatusOK)
	}
}

func (handler *LinkHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		hash := req.PathValue("hash")
		linkFounded, err := handler.LinkRepository.GetByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//handler.StatRepository.AddClick(linkFounded.ID)
		go handler.EventBus.Publish(event.Event{
			Type: event.EventLinkVisited,
			Data: linkFounded.ID,
		})
		http.Redirect(w, req, linkFounded.Url, http.StatusTemporaryRedirect)
	}
}

func (handler *LinkHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		limit, err := strconv.Atoi(req.URL.Query().Get("limit"))
		if err != nil {
			http.Error(w, "Invalid limit", http.StatusResetContent)
			return
		}
		offset, err := strconv.Atoi(req.URL.Query().Get("offset"))
		if err != nil {
			http.Error(w, "Invalid offset", http.StatusResetContent)
			return
		}
		links := handler.LinkRepository.GetAll(limit, offset)
		count := handler.LinkRepository.Count()
		data := &GetAllLinksResponso{
			Links: links,
			Count: count,
		}
		response.Response(w, data, http.StatusOK)
	}
}

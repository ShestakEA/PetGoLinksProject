package stat

import (
	"go/adv-demo/pkg/response"
	"net/http"
	"time"
)

const (
	GroupByDay   = "day"
	GroupByMonth = "month"
)

type StatHandler struct {
	StatRepository *StatRepository
}

type StatHandlerDeps struct {
	StatRepository *StatRepository
}

func NewStatHandler(router *http.ServeMux, deps StatHandlerDeps) {
	handler := &StatHandler{
		StatRepository: deps.StatRepository,
	}
	router.HandleFunc("GET /stat", handler.GetStat())
}

func (handler *StatHandler) GetStat() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		from, err := time.Parse("2006-01-02", req.URL.Query().Get("from"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		to, err := time.Parse("2006-01-02", req.URL.Query().Get("to"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		by := req.URL.Query().Get("by")
		if by != GroupByDay && by != GroupByMonth {
			http.Error(w, "Invalid params", http.StatusBadRequest)
			return
		}
		stats := handler.StatRepository.GetStats(by, from, to)
		response.Response(w, stats, 200)
	}
}

package web

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (s *Server) HandleOrganizationList(w http.ResponseWriter, r *http.Request) {
	h := "HandleOrganizationList "
	s.handleHttpPostOrGetWithLog(h, w, r, func(ctx context.Context, w http.ResponseWriter, r *http.Request, clog *log.Entry) {
		items, err := s.c.OrganizationList(ctx)
		if err != nil {
			clog.WithError(err).Error("error in api.OrganizationList()")
			s.sendResponseByError(w, err, clog)
			return
		}
		s.sendResponseOKWithData(w, items, clog)
	})
}

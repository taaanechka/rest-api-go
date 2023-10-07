package users

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	userservice "github.com/taaanechka/rest-api-go/internal/api-server/services"
	"github.com/taaanechka/rest-api-go/internal/apperror"
	"github.com/taaanechka/rest-api-go/internal/handlers"
	"github.com/taaanechka/rest-api-go/pkg/logging"
)

const (
	usersURL = "/users"
	userURL  = "/users/:uuid"
)

type handler struct {
	service *userservice.Service
	lg      *logging.Logger
}

func NewHandler(lg *logging.Logger, service *userservice.Service) handlers.Handler {
	return &handler{
		service: service,
		lg:      lg,
	}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, usersURL, apperror.Middleware(h.GetList))
	router.HandlerFunc(http.MethodPost, usersURL, apperror.Middleware(h.CreateUser))
	router.HandlerFunc(http.MethodGet, userURL, apperror.Middleware(h.GetUserByUUID))
	router.HandlerFunc(http.MethodPut, userURL, apperror.Middleware(h.UpdateUser))
	router.HandlerFunc(http.MethodPatch, userURL, apperror.Middleware(h.PartiallyUpdateUser))
	router.HandlerFunc(http.MethodDelete, userURL, apperror.Middleware(h.DeleteUser))
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var res CreateUserReq
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	oid, err := h.service.Create(context.Background(), convertCreateUserReqToBL(&res))
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return fmt.Errorf("failed to create user: %w", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("id:%v", oid)))

	return nil
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request) error {
	us, err := h.service.GetList(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	res := make([]GetUserResp, 0, len(us))
	for _, u := range us {
		res = append(res, convertBLToGetUserResp(&u))
	}

	resBytes, err := json.Marshal(res)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)

	return nil
}

func (h *handler) GetUserByUUID(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	res, err := h.service.GetByUUID(context.Background(), id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	resBytes, err := json.Marshal(convertBLToGetUserResp(&res))
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resBytes)

	return nil
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var res UserUpdReq
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	err, found := h.service.Update(context.Background(), id, convertUpdToBL(&res))
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return fmt.Errorf("failed to create user: %w", err)
	}
	if !found {
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusNoContent)
	}

	return nil
}

func (h *handler) PartiallyUpdateUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	var res UserPatchReq
	err = json.Unmarshal(bodyBytes, &res)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	fmt.Printf("data: %#v\n", res)

	err = h.service.Patch(context.Background(), id, convertPatchToBL(&res))
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return fmt.Errorf("failed to create user")
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) error {
	params := httprouter.ParamsFromContext(r.Context())
	id := params.ByName("uuid")

	err := h.service.Delete(context.Background(), id)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusNoContent)
	return nil
}

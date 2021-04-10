package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/pkg/httpcommon"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	logger    *logrus.Logger
	bin       domain.BinUsecase
	validator *validator.Validate
}

func NewHTTPDelivery(router *mux.Router, logger *logrus.Logger, bin domain.BinUsecase) {
	httpInstance := &httpDelivery{
		logger:    logger,
		bin:       bin,
		validator: validator.New(),
	}

	// Bind with given router
	router.HandleFunc("/bin", httpInstance.Select).Methods("GET")
	router.HandleFunc("/bin", httpInstance.Create).Methods("POST")
	router.HandleFunc("/bin/{id}", httpInstance.Get).Methods("GET")
	router.HandleFunc("/bin/{id}", httpInstance.Update).Methods("PUT")
	router.HandleFunc("/bin/{id}", httpInstance.Delete).Methods("DELETE")
}

func (h *httpDelivery) Get(w http.ResponseWriter, r *http.Request) {
	var (
		binID int64
	)

	vars := mux.Vars(r)
	if id, ok := vars["id"]; !ok {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	} else {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			httpcommon.ResponseJSONError(w, http.StatusBadRequest, "ID Must be a number")
			return
		}
		binID = id
	}

	response, err := h.bin.Get(binID)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot find Bin, Make sure you find correct Bin")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusOK, response)
}

func (h *httpDelivery) Select(w http.ResponseWriter, r *http.Request) {
	var (
		queryParam domain.BinQueryParameter
	)

	// Parse Query Parameter
	if err := queryParam.Parse(r.URL.Query()); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Invalid Query")
		return
	}

	responses, err := h.bin.Select(queryParam)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot Query Bin")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusOK, responses)
}

func (h *httpDelivery) Create(w http.ResponseWriter, r *http.Request) {
	var (
		createData domain.BinDataParameter
	)

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot Read Body")
		return
	}

	if err := json.Unmarshal(bodyData, &createData); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Unable to Unmarshal JSON")
		return
	}

	if err := h.validator.Struct(&createData); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Validation Failure, Try Again")
		return
	}

	response, err := h.bin.Create(createData)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "An Error Occured When Creating Bin")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusCreated, response)
}

func (h *httpDelivery) Update(w http.ResponseWriter, r *http.Request) {
	var (
		binID      int64
		updateData domain.BinDataParameter
	)

	vars := mux.Vars(r)
	if id, ok := vars["id"]; !ok {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	} else {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			httpcommon.ResponseJSONError(w, http.StatusBadRequest, "ID Must be a number")
			return
		}
		binID = id
	}

	bodyData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot Read Body")
		return
	}

	if err := json.Unmarshal(bodyData, &updateData); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Unable to Unmarshal JSON")
		return
	}

	if err := h.validator.Struct(&updateData); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Validation Failure, Try Again")
		return
	}

	response, err := h.bin.Update(binID, updateData)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "An Error Occured When Updating Bin")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusCreated, response)
}

func (h *httpDelivery) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		binID int64
	)

	vars := mux.Vars(r)
	if id, ok := vars["id"]; !ok {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	} else {
		id, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			httpcommon.ResponseJSONError(w, http.StatusBadRequest, "ID Must be a number")
			return
		}
		binID = id
	}

	if resp, err := h.bin.Delete(binID); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Unable to Delete Bin")
	} else {
		httpcommon.ResponseJSON(w, http.StatusCreated, resp)
	}
}

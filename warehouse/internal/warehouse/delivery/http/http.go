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
	warehouse domain.WarehouseUsecase
	validator *validator.Validate
}

func NewHTTPDelivery(router *mux.Router, logger *logrus.Logger, warehouse domain.WarehouseUsecase) {
	httpInstance := &httpDelivery{
		logger:    logger,
		warehouse: warehouse,
		validator: validator.New(),
	}

	// Bind with given router
	router.HandleFunc("/warehouse", httpInstance.Select).Methods("GET")
	router.HandleFunc("/warehouse", httpInstance.Create).Methods("POST")
	router.HandleFunc("/warehouse/{id}", httpInstance.Get).Methods("GET")
	router.HandleFunc("/warehouse/{id}", httpInstance.Update).Methods("PUT")
	router.HandleFunc("/warehouse/{id}", httpInstance.Delete).Methods("DELETE")
}

func (h *httpDelivery) Get(w http.ResponseWriter, r *http.Request) {
	var (
		warehouseID int64
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
		warehouseID = id
	}

	response, err := h.warehouse.Get(warehouseID)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot find warehouse, Make sure you find correct warehouse")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusOK, response)
}

func (h *httpDelivery) Select(w http.ResponseWriter, r *http.Request) {
	var (
		queryParam domain.WarehouseQueryParameter
	)

	// Parse Query Parameter
	if err := queryParam.Parse(r.URL.Query()); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Invalid Query")
		return
	}

	responses, err := h.warehouse.Select(queryParam)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Cannot Query Warehouses")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusOK, responses)
}

func (h *httpDelivery) Create(w http.ResponseWriter, r *http.Request) {
	var (
		createData domain.WarehouseDataParameter
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

	response, err := h.warehouse.Create(createData)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "An Error Occured When Creating Warehouse")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusCreated, response)
}

func (h *httpDelivery) Update(w http.ResponseWriter, r *http.Request) {
	var (
		warehouseID int64
		updateData  domain.WarehouseDataParameter
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
		warehouseID = id
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

	response, err := h.warehouse.Update(warehouseID, updateData)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "An Error Occured When Updating Warehouse")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusCreated, response)
}

func (h *httpDelivery) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		warehouseID int64
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
		warehouseID = id
	}

	if resp, err := h.warehouse.Delete(warehouseID); err != nil {
		httpcommon.ResponseJSONError(w, http.StatusBadRequest, "Unable to Delete Warehouse")
	} else {
		httpcommon.ResponseJSON(w, http.StatusCreated, resp)
	}
}

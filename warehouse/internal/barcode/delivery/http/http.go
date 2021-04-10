package http

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/alvinradeka/jamblang-hakenton/warehouse/internal/domain"
	"github.com/alvinradeka/jamblang-hakenton/warehouse/pkg/httpcommon"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type httpDelivery struct {
	logger    *logrus.Logger
	sku       domain.SKUUsecase
	warehouse domain.WarehouseUsecase
	barcode   domain.BarcodeUsecase
	validator *validator.Validate
}

func NewHTTPDelivery(router *mux.Router, logger *logrus.Logger, sku domain.SKUUsecase, warehouse domain.WarehouseUsecase, barcode domain.BarcodeUsecase) {
	httpInstance := &httpDelivery{
		logger:    logger,
		sku:       sku,
		warehouse: warehouse,
		barcode:   barcode,
		validator: validator.New(),
	}

	// Bind with given router
	router.HandleFunc("/barcode/upload", httpInstance.BarcodeUpload).Methods("POST")
}

func (h *httpDelivery) BarcodeUpload(w http.ResponseWriter, r *http.Request) {
	// Read File
	file, _, err := r.FormFile("barcode_image")
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusInternalServerError, "Cannot Read Barcode Image")
		return
	}

	// Open File temporarily
	tempFile, err := ioutil.TempFile("/tmp", "uploadedfile")
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusInternalServerError, "Cannot Open Temporary Image")
		return
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	io.Copy(tempFile, file)

	resp, err := h.barcode.ParseBarcodeFromFileToLambda(tempFile)
	if err != nil {
		httpcommon.ResponseJSONError(w, http.StatusInternalServerError, "Cannot Process Temporary Image")
		return
	}

	httpcommon.ResponseJSON(w, http.StatusOK, resp)
}

package handler

import (
	"HRD/internal/service"
	"HRD/model"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type KendaraanHandler struct {
    service *service.KendaraanService // Langsung struct service
}

func NewKendaraanHandler(service *service.KendaraanService) *KendaraanHandler {
    return &KendaraanHandler{service: service}
}

// POST /api/kendaraan
func (h *KendaraanHandler) Create(c echo.Context) error {
    // 1. Ambil data form
    nama := c.FormValue("nama_pengemudi")
    mobil := c.FormValue("model_mobil")
    lokasi := c.FormValue("lokasi_sekarang")
    bbm := c.FormValue("bbm")

    // 2. Ambil file
    file, _ := c.FormFile("gambar")

    input := model.Kendaraan{
        NamaPengemudi:  nama,
        ModelMobil:     mobil,
        LokasiSekarang: lokasi,
        Bbm:      bbm,
    }
	fmt.Println(input)

    // 3. Panggil Service (Tanpa Context)
    if err := h.service.CreateKendaraan(input, file); err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
    }

    return c.JSON(http.StatusCreated, map[string]string{"message": "Berhasil disimpan"})
}

// GET /api/kendaraan
func (h *KendaraanHandler) GetAll(c echo.Context) error {
    data, err := h.service.GetAllKendaraan()
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "data": data,
    })
}
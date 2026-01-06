package service

import (
	"HRD/internal/repository"
	"HRD/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

type KendaraanService struct {
    repo *repository.KendaraanRepository // Langsung tembak ke struct repository
}

func NewKendaraanService(repo *repository.KendaraanRepository) *KendaraanService {
    return &KendaraanService{repo: repo}
}

func (s *KendaraanService) CreateKendaraan(input model.Kendaraan, file *multipart.FileHeader) error {
    // 1. Logic Upload File
    if file != nil {
        fileURL, err := s.saveFile(file, "vehicle")
        if err != nil {
            return err
        }
        input.GambarURL = fileURL
    } else {
        input.GambarURL = ""
    }

    // 2. Panggil Repo (Tanpa Context)
    return s.repo.Save(&input)
}

func (s *KendaraanService) GetAllKendaraan() ([]model.Kendaraan, error) {
    data, err := s.repo.FindAll()
    if err != nil {
        return nil, err
    }

    loc, _ := time.LoadLocation("Asia/Jakarta")

    for i := range data {
        t := data[i].WaktuInput

        // 1️⃣ Anggap waktu ini adalah WIB
        tWIB := time.Date(
            t.Year(), t.Month(), t.Day(),
            t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
            loc,
        )

        // 2️⃣ Konversi ke UTC (INI BARU BENAR)
        data[i].WaktuInput = tWIB.UTC()
    }

    return data, nil
}

// --- Helper Simpan File (Logic Anti-Blob) ---
func (s *KendaraanService) saveFile(fileHeader *multipart.FileHeader, prefix string) (string, error) {
    uploadDir := filepath.Join("uploads", "kendaraan")
    
    if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
        os.MkdirAll(uploadDir, os.ModePerm)
    }

    src, err := fileHeader.Open()
    if err != nil {
        return "", err
    }
    defer src.Close()

    // Deteksi Ekstensi
    ext := filepath.Ext(fileHeader.Filename)
    if ext == "" || ext == ".blob" {
        contentType := fileHeader.Header.Get("Content-Type")
        switch contentType {
        case "image/jpeg":
            ext = ".jpg"
        case "image/png":
            ext = ".png"
        default:
            ext = ".jpg"
        }
    }

    filename := fmt.Sprintf("%s_%d%s", prefix, time.Now().Unix(), ext)
    path := filepath.Join(uploadDir, filename)

    dst, err := os.Create(path)
    if err != nil {
        return "", err
    }
    defer dst.Close()

    if _, err = io.Copy(dst, src); err != nil {
        return "", err
    }

    return path, nil
}
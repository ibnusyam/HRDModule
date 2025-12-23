package model

import "time"

type Kendaraan struct {
    ID            int    `json:"id"`
    NamaPengemudi string `json:"nama_pengemudi"`
    ModelMobil    string `json:"model_mobil"`
    LokasiSekarang string `json:"lokasi_sekarang"`
    Bbm    string `json:"bbm"`
    GambarURL     string `json:"gambar_url"`
    WaktuInput     time.Time `json:"waktu_input"`
}
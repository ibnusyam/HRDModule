package repository

import (
	"HRD/model"
	"database/sql"
)

// Langsung struct konkrit
type KendaraanRepository struct {
    db *sql.DB
}

func NewKendaraanRepository(db *sql.DB) *KendaraanRepository {
    return &KendaraanRepository{db: db}
}

func (r *KendaraanRepository) Save(k *model.Kendaraan) error {
    query := `
        INSERT INTO kendaraan (nama_pengemudi, model_mobil, lokasi_sekarang, bbm, gambar_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, waktu_input
    `
    // Scan id DAN waktu_input yang digenerate oleh database
    err := r.db.QueryRow(query, 
        k.NamaPengemudi, 
        k.ModelMobil, 
        k.LokasiSekarang, 
        k.Bbm, 
        k.GambarURL,
    ).Scan(&k.ID, &k.WaktuInput)

    return err
}

func (r *KendaraanRepository) FindAll() ([]model.Kendaraan, error) {
    // Tambahkan waktu_input dalam SELECT
    query := `SELECT id, nama_pengemudi, model_mobil, lokasi_sekarang, bbm, gambar_url, waktu_input FROM kendaraan ORDER BY id DESC`
    
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.Kendaraan
    for rows.Next() {
        var k model.Kendaraan
        // Scan field waktu_input
        if err := rows.Scan(&k.ID, &k.NamaPengemudi, &k.ModelMobil, &k.LokasiSekarang, &k.Bbm, &k.GambarURL, &k.WaktuInput); err != nil {
            return nil, err
        }
        results = append(results, k)
    }
    return results, nil
}
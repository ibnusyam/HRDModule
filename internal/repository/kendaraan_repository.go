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
        INSERT INTO kendaraan (nama_pengemudi, model_mobil, lokasi_sekarang, bbm ,gambar_url)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `
    // Pakai QueryRow biasa (tanpa Context)
    err := r.db.QueryRow(query, 
        k.NamaPengemudi, 
        k.ModelMobil, 
        k.LokasiSekarang, 
        k.Bbm, 
        k.GambarURL,
    ).Scan(&k.ID)

    return err
}

func (r *KendaraanRepository) FindAll() ([]model.Kendaraan, error) {
    query := `SELECT id, nama_pengemudi, model_mobil, lokasi_sekarang, bbm, gambar_url FROM kendaraan ORDER BY id DESC`
    
    // Pakai Query biasa (tanpa Context)
    rows, err := r.db.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var results []model.Kendaraan
    for rows.Next() {
        var k model.Kendaraan
        if err := rows.Scan(&k.ID, &k.NamaPengemudi, &k.ModelMobil, &k.LokasiSekarang, &k.Bbm, &k.GambarURL); err != nil {
            return nil, err
        }
        results = append(results, k)
    }
    return results, nil
}
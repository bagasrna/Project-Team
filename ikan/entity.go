package ikan

type Ikan struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Kategori    string `json:"kategori"`
	Jenis_Ikan  string `json:"jenis_ikan"`
	Harga       string `json:"harga"`
	TokoID      uint   `json:"toko_id"`
	Provinsi    string `json:"provinsi"`
	Kota        string `json:"kota"`
	Bulan_Panen string `json:"bulan_panen"`
	Deskripsi_Produk string `json:"deskripsi_produk"`
}
package ikan

type Service interface {
	TambahIkan(ikan Ikan) (Ikan, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) TambahIkan(input Ikan) (Ikan, error) {

	ikan := Ikan{
		Kategori:         input.Kategori,
		Jenis_Ikan:       input.Jenis_Ikan,
		Harga:            input.Harga,
		TokoID:           input.TokoID,
		Provinsi:         input.Provinsi,
		Kota:             input.Kota,
		Bulan_Panen:      input.Bulan_Panen,
		Deskripsi_Produk: input.Deskripsi_Produk,
	}

	registerdIkan, err := s.repository.TambahIkan(ikan)
	if err != nil {
		return registerdIkan, err
	}
	return registerdIkan, nil
}

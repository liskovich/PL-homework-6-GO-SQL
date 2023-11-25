package db

import (
	"database/sql"

	"example.com/api/model"
)

type BeerRepository interface {
	CreateBeer(beer model.Beer) error
	UpdateBeer(beerID uint, beer model.Beer) (*model.Beer, error)
	GetBeerById(beerID uint) (*model.Beer, error)
	GetAllBeers() ([]*model.Beer, error)
}

type beerRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) BeerRepository {
	return &beerRepo{db: db}
}

func (bRepo *beerRepo) CreateBeer(beer model.Beer) error {
	_, err := bRepo.db.Exec(InsertBeerQuery, beer.Name, beer.Description, beer.Thumbnail)
	return err
}

func (bRepo *beerRepo) GetAllBeers() ([]*model.Beer, error) {
	rows, err := bRepo.db.Query(SelectAllBeersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beers := []*model.Beer{}
	for rows.Next() {
		var beer model.Beer
		if err := rows.Scan(&beer.ID, &beer.Name, &beer.Description, &beer.Thumbnail); err != nil {
			return nil, err
		}
		beers = append(beers, &beer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return beers, nil
}

func (bRepo *beerRepo) GetBeerById(beerID uint) (*model.Beer, error) {
	row := bRepo.db.QueryRow(SelectBeerByIdQuery, beerID)
	var beer model.Beer
	err := row.Scan(&beer.ID, &beer.Name, &beer.Description, &beer.Thumbnail)
	if err != nil {
		return nil, err
	}
	return &beer, nil
}

func (bRepo *beerRepo) UpdateBeer(beerID uint, beer model.Beer) (*model.Beer, error) {
	_, err := bRepo.db.Exec(UpdateBeerQuery, beer.Name, beer.Description, beer.Thumbnail, beerID)
	if err != nil {
		return nil, err
	}

	updatedBeer, err := bRepo.GetBeerById(beerID)
	if err != nil {
		return nil, err
	}
	return updatedBeer, nil
}

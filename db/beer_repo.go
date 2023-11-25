package db

import (
	"database/sql"

	"example.com/api/model"
)

type BeerRepository interface {
	CreateBeer(beer model.BeerMutate) error
	UpdateBeer(beerID uint, beer model.BeerMutate) (*model.BeerCompact, error)
	DeleteBeer(beerID uint) error
	GetBeerById(beerID uint) (*model.BeerCompact, error)
	GetAllBeers() ([]*model.BeerCompact, error)
}

type beerRepo struct {
	db *sql.DB
}

func NewBeerRepository(db *sql.DB) BeerRepository {
	return &beerRepo{db: db}
}

func (bRepo *beerRepo) CreateBeer(beer model.BeerMutate) error {
	_, err := bRepo.db.Exec(InsertBeerQuery, beer.Name, beer.Description, beer.Thumbnail)
	return err
}

func (bRepo *beerRepo) DeleteBeer(beerID uint) error {
	_, err := bRepo.db.Exec(DeleteBeerQuery, beerID)
	return err
}

func (bRepo *beerRepo) GetAllBeers() ([]*model.BeerCompact, error) {
	rows, err := bRepo.db.Query(SelectAllBeersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beers := []*model.BeerCompact{}
	for rows.Next() {
		var beer model.BeerCompact
		// TODO: add the "upvotes_count" and "comments_count"
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

func (bRepo *beerRepo) GetBeerById(beerID uint) (*model.BeerCompact, error) {
	row := bRepo.db.QueryRow(SelectBeerByIdQuery, beerID)
	var beer model.BeerCompact
	// TODO: add the "upvotes_count", "comments_count" and "comments"
	err := row.Scan(&beer.ID, &beer.Name, &beer.Description, &beer.Thumbnail)
	if err != nil {
		return nil, err
	}
	return &beer, nil
}

func (bRepo *beerRepo) UpdateBeer(beerID uint, beer model.BeerMutate) (*model.BeerCompact, error) {
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

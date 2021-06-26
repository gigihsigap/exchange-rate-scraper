package repository

import (
	"backend/app/form"
	"backend/app/helpers"
	"backend/app/model"
	"backend/app/scraper"
	"backend/db"
	"errors"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var EntryEntity IEntry

type entryEntity struct {
	resource *db.Resource
	repo     *mongo.Collection
}

type IEntry interface {
	// GetAll() ([]model.Entry, int, error)
	GetByDate(date string) (*model.Entry, int, error)
	GetEntryListByDate(startDate, endDate string) ([]model.Entry, int, error)
	GetEntryListBySymbol(startDate, endDate, symbol string) ([]model.Entry, int, error)
	CreateOneEntry(entryForm form.Entry) (*model.Entry, int, error)

	UpdateEntry(date string, entryForm form.Entry) (model.Entry, int, error)
	DeleteEntry(date string) (model.Entry, int, error)
	Indexing() (model.Entry, int, error)
}

//func NewEntryEntity
func NewEntryEntity(resource *db.Resource) IEntry {
	entryRepo := resource.DB.Collection("entry")
	EntryEntity = &entryEntity{resource: resource, repo: entryRepo}
	return EntryEntity
}

func (entity *entryEntity) GetByDate(date string) (*model.Entry, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	dateTime, err := helpers.YYYYMMDDToISODate(date)

	var entry model.Entry
	err = entity.repo.FindOne(ctx, bson.M{"date": dateTime}).Decode(&entry)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &entry, http.StatusOK, nil
}

func (entity *entryEntity) CreateOneEntry(entryForm form.Entry) (*model.Entry, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	dateTime, err := helpers.YYYYMMDDToISODate(entryForm.Date)

	newEntry := model.Entry{
		Id:   primitive.NewObjectID(),
		Date: dateTime,
		Exc_Rates: []model.Exc_rates{
			{
				Symbol: entryForm.Symbol,
				Rates: model.Rates{
					ER: model.Buysell{
						Buy:  entryForm.ER.Buy,
						Sell: entryForm.ER.Sell,
					},
					TT: model.Buysell{
						Buy:  entryForm.TT.Buy,
						Sell: entryForm.TT.Sell,
					},
					BN: model.Buysell{
						Buy:  entryForm.BN.Buy,
						Sell: entryForm.BN.Sell,
					},
				},
			},
		},
		Created_At: time.Now(),
		Updated_At: time.Now(),
	}

	found, _, _ := entity.GetByDate(entryForm.Date)
	if found != nil {
		return nil, http.StatusBadRequest, errors.New("Entry at that date already exists")
	}
	_, err = entity.repo.InsertOne(ctx, newEntry)

	if err != nil {
		logrus.Print(err)
		return nil, 400, err
	}

	return &newEntry, http.StatusOK, nil
}

func (entity *entryEntity) UpdateEntry(date string, entryForm form.Entry) (model.Entry, int, error) {
	ctx, cancel := initContext()
	defer cancel()

	entry, _, err := entity.GetByDate(date)
	if err != nil {
		return model.Entry{}, http.StatusNotFound, nil
	}

	dateTime, err := helpers.YYYYMMDDToISODate(date)

	newEntry := model.Entry{
		Id:         entry.Id,
		Date:       dateTime,
		Exc_Rates:  entry.Exc_Rates,
		Updated_At: time.Now(),
	}

	inputRate := model.Exc_rates{
		Symbol: entryForm.Symbol,
		Rates: model.Rates{
			ER: model.Buysell{
				Buy:  entryForm.ER.Buy,
				Sell: entryForm.ER.Sell,
			},
			TT: model.Buysell{
				Buy:  entryForm.TT.Buy,
				Sell: entryForm.TT.Sell,
			},
			BN: model.Buysell{
				Buy:  entryForm.BN.Buy,
				Sell: entryForm.BN.Sell,
			},
		},
	}

	// If currency already exists, update the rates
	// If currency data hasn't existed yet, append the slice
	index, isNewCurrency := helpers.Contains(entry.Exc_Rates, entryForm.Symbol)

	if isNewCurrency {
		newEntry.Exc_Rates = append(entry.Exc_Rates, inputRate)
	} else {
		newEntry.Exc_Rates[index] = inputRate
	}

	isReturnNewDoc := options.After
	opts := &options.FindOneAndUpdateOptions{
		ReturnDocument: &isReturnNewDoc,
	}

	err = entity.repo.FindOneAndUpdate(ctx, bson.M{"date": dateTime}, bson.M{"$set": newEntry}, opts).Decode(&entry)

	if err != nil {
		logrus.Error(err)
		return model.Entry{}, getHTTPCode(err), err
	}

	return newEntry, http.StatusOK, nil
}

func (entity *entryEntity) DeleteEntry(date string) (model.Entry, int, error) {
	var entry *model.Entry
	ctx, cancel := initContext()

	defer cancel()

	entry, _, err := entity.GetByDate(date)
	if err != nil {
		return model.Entry{}, http.StatusNotFound, nil
	}

	dateTime, err := helpers.YYYYMMDDToISODate(date)

	err = entity.repo.FindOneAndDelete(ctx, bson.M{"date": dateTime}).Decode(&entry)

	return *entry, http.StatusOK, nil
}

func (entity *entryEntity) GetEntryListByDate(startDate, endDate string) ([]model.Entry, int, error) {
	entryList := []model.Entry{}
	ctx, cancel := initContext()
	defer cancel()

	startTime, err := helpers.YYYYMMDDToISODate(startDate)
	endTime, err := helpers.YYYYMMDDToISODate(endDate)

	cursor, err := entity.repo.Find(
		ctx,
		bson.M{
			"date": bson.M{"$gte": startTime, "$lte": endTime},
		},
	)

	if err != nil {
		logrus.Print(err)
		return []model.Entry{}, 400, err
	}

	for cursor.Next(ctx) {
		var entry model.Entry
		err = cursor.Decode(&entry)
		if err != nil {
			logrus.Print(err)
		}
		entryList = append(entryList, entry)
	}
	return entryList, http.StatusOK, nil
}

func (entity *entryEntity) GetEntryListBySymbol(startDate, endDate, symbol string) ([]model.Entry, int, error) {
	entryList := []model.Entry{}
	ctx, cancel := initContext()
	defer cancel()

	startTime, err := helpers.YYYYMMDDToISODate(startDate)
	endTime, err := helpers.YYYYMMDDToISODate(endDate)

	// Queries entry data that includes matching currency
	// which can also include different currency on the same date
	cursor, err := entity.repo.Find(ctx,
		bson.M{
			"exc_rates": bson.M{"$elemMatch": bson.M{"symbol": symbol}},
			"date":      bson.M{"$gte": startTime, "$lte": endTime},
		},
	)

	if err != nil {
		logrus.Print(err)
		return []model.Entry{}, 400, err
	}

	for cursor.Next(ctx) {
		var entry model.Entry
		err = cursor.Decode(&entry)
		if err != nil {
			logrus.Print(err)
		}
		entryList = append(entryList, entry)
	}

	// Since we only want one currency type,
	// the next step is to clean unwanted currency data
	for i, a := range entryList {
		for j, b := range a.Exc_Rates {
			if b.Symbol == symbol {
				entryList[i].Exc_Rates = []model.Exc_rates{a.Exc_Rates[j]}
			}
		}
	}

	return entryList, http.StatusOK, nil
}

func (entity *entryEntity) Indexing() (model.Entry, int, error) {
	var entry *model.Entry
	ctx, cancel := initContext()
	defer cancel()

	// Check if entry data at current date already exists
	currentTime := time.Now().Format("2006-01-02")

	entry, _, err := entity.GetByDate(currentTime)

	if entry != nil {
		return model.Entry{}, http.StatusNotFound, errors.New("Scrapping failed because data at current date already exists. Maybe you already ran it once?")
	} else {
		entry, err = scraper.Run()

		if err != nil {
			logrus.Print(err)
			// return nil, 400, err
		}

		_, err = entity.repo.InsertOne(ctx, entry)

		if err != nil {
			logrus.Print(err)
			// return nil, 400, err
		}
	}

	if err != nil {
		return model.Entry{}, http.StatusNotFound, nil
	}

	if err != nil {
		logrus.Print(err)
		// return nil, 400, err
	}

	return *entry, http.StatusOK, nil
}

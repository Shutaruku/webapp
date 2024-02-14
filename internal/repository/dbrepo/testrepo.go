package dbrepo

import (
	"errors"
	"log"
	"time"

	"github.com/YuanData/webapp/internal/models"
)

func (m *testDBRepo) AllUsers() bool {
	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.BungalowID == 99 {
		return 0, errors.New("some error")
	}

	return 1, nil
}

func (m *testDBRepo) InsertBungalowRestriction(r models.BungalowRestriction) error {
	if r.BungalowID == 999 {
		return errors.New("just because")
	}

	return nil
}

func (m *testDBRepo) SearchAvailabilityByDatesByBungalowID(start, end time.Time, bungalowID int) (bool, error) {
	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return false, errors.New("some error")
	}

	if start.After(t) {
		return false, nil
	}

	return true, nil
}

func (m *testDBRepo) SearchAvailabilityByDatesForAllBungalows(start, end time.Time) ([]models.Bungalow, error) {
	var bungalows []models.Bungalow

	layout := "2006-01-02"
	str := "2036-12-31"
	t, err := time.Parse(layout, str)
	if err != nil {
		log.Println(err)
	}

	testDateToFail, err := time.Parse(layout, "2038-01-01")
	if err != nil {
		log.Println(err)
	}

	if start == testDateToFail {
		return bungalows, errors.New("some error")
	}

	if start.After(t) {
		return bungalows, nil
	}

	bungalow := models.Bungalow{
		ID: 1,
	}
	bungalows = append(bungalows, bungalow)

	return bungalows, nil
}

func (m *testDBRepo) GetBungalowByID(id int) (models.Bungalow, error) {
	var bungalow models.Bungalow
	if id > 3 {
		return bungalow, errors.New("an error occured")
	}

	return bungalow, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}

func (m *testDBRepo) UpdateUser(u models.User) error {
	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "", nil
}

func (m *testDBRepo) AllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var res models.Reservation

	return res, nil
}

func (m *testDBRepo) UpdateReservation(r models.Reservation) error {

	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

func (m *testDBRepo) UpdateStatusOfReservation(id, status int) error {

	return nil
}

func (m *testDBRepo) AllBungalows() ([]models.Bungalow, error) {

	var bungalows []models.Bungalow

	return bungalows, nil
}

func (m *testDBRepo) GetRestrictionsForBungalowByDate(bungalowID int, start, end time.Time) ([]models.BungalowRestriction, error) {

	var restrictions []models.BungalowRestriction

	return restrictions, nil
}

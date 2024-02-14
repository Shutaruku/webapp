package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/YuanData/webapp/internal/driver"
	"github.com/YuanData/webapp/internal/models"
)

type postData struct {
	key   string
	value string
}

var allTheHandlerTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"eremite", "/eremite", "GET", http.StatusOK},
	{"couple", "/couple", "GET", http.StatusOK},
	{"family", "/family", "GET", http.StatusOK},
	{"reservation", "/reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"not-existing-route", "/not-existing-dummy", "GET", http.StatusNotFound},
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	testRepo := NewRepo(&app, &db)

	if reflect.TypeOf(testRepo).String() != "*handlers.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(testRepo).String())
	}
}

func TestAllTheHandlers(t *testing.T) {

	routes := getRoutes()

	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range allTheHandlerTests {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("%s: expected %d, got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

func TestRepository_PostReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")

	req, _ := http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post availability when no bungalows available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	postedData = url.Values{}
	postedData.Add("start", "2036-01-01")
	postedData.Add("end", "2036-01-02")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Post availability when bungalows are available gave wrong status code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("POST", "/reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with empty request body (nil) gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2037-01-02")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid start date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "invalid")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability with invalid end date gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "2038-01-01")
	postedData.Add("end", "2038-01-02")

	req, _ = http.NewRequest("POST", "/reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Post availability when database query fails gave wrong status code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_MakeReservation(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.BungalowID = 99
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("handler MakeReservation failed: unexpected response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

}

func TestRepository_PostMakeReservation(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	layout := "2006-01-02"
	sd, _ := time.Parse(layout, "2037-01-01")
	ed, _ := time.Parse(layout, "2037-01-02")
	bungalowId, _ := strconv.Atoi("1")

	reservation := models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx := getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	req, _ = http.NewRequest("POST", "/make-reservation", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler returned wrong response code for missing session data: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("full_name", "P")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("1")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("PostMakeReservation handler returned wrong response code invalid/insufficient data: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("99")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("full_name", "Peter Griffin")
	postedData.Add("email", "peter@griffin.family")
	postedData.Add("phone", "1234567890")

	layout = "2006-01-02"
	sd, _ = time.Parse(layout, "2037-01-01")
	ed, _ = time.Parse(layout, "2037-01-02")
	bungalowId, _ = strconv.Atoi("999")

	reservation = models.Reservation{
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: bungalowId,
		Bungalow: models.Bungalow{
			BungalowName: "some bungalow name for tests",
		},
	}

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostMakeReservation)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostMakeReservation handler failed when trying to inserting a reservation into the database: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_ReservationJSON(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "invalid")

	req, _ := http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid bungalow ID: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "invalid")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "1")

	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid start date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "invalid")
	postedData.Add("bungalow_id", "1")

	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationJSON handler unexpectedly seems able parsing invalid end date: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	postedData = url.Values{}
	postedData.Add("start", "2037-01-01")
	postedData.Add("end", "2037-01-02")
	postedData.Add("bungalow_id", "1")

	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("can't parse json")
	}

	if j.OK {
		t.Error("got availability, unexpected none")
	}

	postedData = url.Values{}
	postedData.Add("start", "2036-01-01")
	postedData.Add("end", "2036-01-02")
	postedData.Add("bungalow_id", "1")

	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if !j.OK {
		t.Error("got no availability but expected in handler ReservationJSON")
	}

	req, _ = http.NewRequest("POST", "/reservation-json", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if j.OK || j.Message != "Internal server error" {
		t.Error("got availability with empty request body")
	}

	postedData = url.Values{}
	postedData.Add("start", "2038-01-01")
	postedData.Add("end", "2038-01-02")
	postedData.Add("bungalow_id", "1")

	req, _ = http.NewRequest("POST", "/reservation-json", strings.NewReader(postedData.Encode()))

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationJSON)

	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json!")
	}

	if j.OK || j.Message != "Error querying database" {
		t.Error("got availability, unexpected because database returned an error")
	}
}

func TestRepository_ReservationOverview(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/reservation-overview", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	reservation = models.Reservation{
		BungalowID: 99,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ = http.NewRequest("GET", "/reservation-overview", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/reservation-overview", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ReservationOverview)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ReservationOverview handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}
}

func TestRepository_ChooseBungalow(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-bungalow/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.RequestURI = "/choose-bungalow/1"

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	req, _ = http.NewRequest("GET", "/choose-bungalow/1", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-bungalow/1"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/choose-bungalow/fish", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	req.RequestURI = "/choose-bungalow/fish"

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("ChooseBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_BookBungalow(t *testing.T) {

	reservation := models.Reservation{
		BungalowID: 1,
		Bungalow: models.Bungalow{
			ID:           1,
			BungalowName: "The Solitude Shack",
		},
	}

	req, _ := http.NewRequest("GET", "/book-bungalow?s=2038-01-01&e=2038-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	req, _ = http.NewRequest("GET", "/book-bungalow?s=2036-01-01&e=2036-01-02&id=99", nil)

	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookBungalow)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("BookBungalow handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}

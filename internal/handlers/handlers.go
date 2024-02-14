package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YuanData/webapp/internal/config"
	"github.com/YuanData/webapp/internal/driver"
	"github.com/YuanData/webapp/internal/forms"
	"github.com/YuanData/webapp/internal/models"
	"github.com/YuanData/webapp/internal/render"
	"github.com/YuanData/webapp/internal/repository"
	"github.com/YuanData/webapp/internal/repository/dbrepo"
)

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

var Repo *Repository

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home-page.tpml", &models.TemplateData{})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about-page.tpml", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact-page.tpml", &models.TemplateData{})
}

func (m *Repository) Eremite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "eremite-page.tpml", &models.TemplateData{})
}

func (m *Repository) Couple(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "couple-page.tpml", &models.TemplateData{})
}

func (m *Repository) Family(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "family-page.tpml", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "check-availability-page.tpml", &models.TemplateData{})
}

func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	bungalows, err := m.DB.SearchAvailabilityByDatesForAllBungalows(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if len(bungalows) == 0 {
		m.App.Session.Put(r.Context(), "error", ":( No holiday home is available at that time.")
		http.Redirect(w, r, "/reservation", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["bungalows"] = bungalows

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-bungalow-page.tpml", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK         bool   `json:"ok"`
	Message    string `json:"message"`
	BungalowID string `json:"bungalow_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

func (m *Repository) ReservationJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		output, _ := json.MarshalIndent(resp, "", "    ")

		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
		return
	}

	bungalowID, err := strconv.Atoi(r.Form.Get("bungalow_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get data from form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	available, err := m.DB.SearchAvailabilityByDatesByBungalowID(startDate, endDate, bungalowID)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error querying database",
		}

		output, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(output)
		return
	}

	resp := jsonResponse{
		OK:         available,
		Message:    "",
		StartDate:  sd,
		EndDate:    ed,
		BungalowID: strconv.Itoa(bungalowID),
	}

	output, _ := json.MarshalIndent(resp, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}

func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation back from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	bungalow, err := m.DB.GetBungalowByID(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find bungalow!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation-page.tpml", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) PostMakeReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation back from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FullName:   r.Form.Get("full_name"),
		Email:      r.Form.Get("email"),
		Phone:      r.Form.Get("phone"),
		StartDate:  res.StartDate,
		EndDate:    res.EndDate,
		BungalowID: res.BungalowID,
		Bungalow: models.Bungalow{
			BungalowName: res.Bungalow.BungalowName,
		},
	}

	//validate form data
	form := forms.New(r.PostForm)

	form.Required("full_name", "email")
	form.MinLength("full_name", 2)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		sd := res.StartDate.Format("2006-01-02")
		ed := res.EndDate.Format("2006-01-02")

		stringMap := make(map[string]string)
		stringMap["start_date"] = sd
		stringMap["end_date"] = ed

		m.App.Session.Put(r.Context(), "reservation", reservation)

		render.Template(w, r, "make-reservation-page.tpml", &models.TemplateData{
			Form:      form,
			Data:      data,
			StringMap: stringMap,
		})
		return
	}

	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't write reservation to database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.BungalowRestriction{
		StartDate:     res.StartDate,
		EndDate:       res.EndDate,
		BungalowID:    res.BungalowID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertBungalowRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't reserve bungalow in database")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	htmlMessage := fmt.Sprintf(`
	<strong>Receipt of a request for a reservation</strong><br><br>
	Dear %s: <br>
	we received your reservation request to rent the our bungalow "%s" from %s to %s.
	`, reservation.FullName, res.Bungalow.BungalowName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:      reservation.Email,
		From:    "noreply@sys.com",
		Subject: "Receipt of a request for a reservation",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	htmlMessage = fmt.Sprintf(`
		<strong>New Reservation Request</strong><br>
		we received a new reservation request to rent the bungalow "%s" from %s to %s.
		`, res.Bungalow.BungalowName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))

	msg = models.MailData{
		To:      "ta@target.com",
		From:    "noreply@sys.com",
		Subject: "New Reservation Request",
		Content: htmlMessage,
	}
	m.App.MailChan <- msg

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservation-overview", http.StatusSeeOther)
}

func (m *Repository) ReservationOverview(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation back from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	m.App.Session.Remove(r.Context(), "reservation")

	bungalow, err := m.DB.GetBungalowByID(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't find bungalow!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName

	data := make(map[string]interface{})
	data["reservation"] = res

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-overview-page.tpml", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) ChooseBungalow(w http.ResponseWriter, r *http.Request) {

	exploded := strings.Split(r.RequestURI, "/")
	bungalowID, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Missing parameter from URL")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "can't get reservation back from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.BungalowID = bungalowID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) BookBungalow(w http.ResponseWriter, r *http.Request) {

	bungalowID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	var res models.Reservation

	bungalow, err := m.DB.GetBungalowByID(bungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find bungalow!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Bungalow.BungalowName = bungalow.BungalowName
	res.BungalowID = bungalowID
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login-page.tpml", &models.TemplateData{
		Form: forms.New(nil),
	})
}

func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login-page.tpml", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Invalid credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "success", "Successfully logged in")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)

}

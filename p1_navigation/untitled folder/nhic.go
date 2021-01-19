package nhic

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"gitlab.lean/leandevclan/nhic/config"
	"gitlab.lean/leandevclan/nhic/nic"
	"gitlab.lean/leandevclan/nhic/oauth"
	"gitlab.lean/leandevclan/nhic/scfhs"
	"gitlab.lean/leandevclan/nhic/store"
	"gitlab.lean/leandevclan/nhic/yakeen"
)

const (
	hoursPerYear float64 = 8760
)

var (
	ErrBadArgs            = errors.New("missing individual arguments")
	ErrBadNationalID      = errors.New("malformed national_id")
	ErrBadIqamaID         = errors.New("malformed iqama_id")
	ErrBadBirthDate       = errors.New("malformed birth_date")
	ErrBadExpiryDate      = errors.New("malformed expiry_date")
	ErrUnknownPatientType = errors.New("patient type is unknown")
	ErrFetchingInfo       = errors.New("encountered error while fetch information")  // yakeen
	ErrLookingUpInfo      = errors.New("encountered error while lookup information") // store
	ErrUpdateInfo         = errors.New("encountered error while update information")
	ErrSearchInput        = errors.New("search input error")
	ErrNotFound           = errors.New("no info found")
)

// PatientKind defines the type of the patient we're handling
// Valid values are:
// Expat == KindExpat
// Citizen == KindCitizen
type PatientKind int

const (
	KindUnknown PatientKind = iota + 1
	KindExpat
	KindCitizen
)

const (
	citizenPrefix = "1"
	expatPrefix   = "2"
)

type PatientQuery struct {
	ID        string
	BirthDate string
}

// Kind returns the patient type
func (pq *PatientQuery) Kind() PatientKind {
	if len(pq.ID) == 10 && strings.HasPrefix(pq.ID, citizenPrefix) {
		return KindCitizen
	}
	if len(pq.ID) > 4 && strings.HasPrefix(pq.ID, expatPrefix) {
		return KindExpat
	}
	return KindUnknown
}

// Parse creates a populate a PatientQuery struct from query parameters
func (pq *PatientQuery) Parse(u *url.URL) {
	q := u.Query()
	pq.ID = q.Get("id")
	pq.BirthDate = q.Get("birth_date")
}

// Validate checks if PatientQuery fields are valid, depending on the type
func (pq *PatientQuery) Validate() error {
	switch pq.Kind() {
	case KindCitizen:
		if len(pq.ID) != 10 && strings.HasPrefix(pq.ID, citizenPrefix) {
			return ErrBadNationalID
		}
		if len(pq.BirthDate) <= 4 {
			return ErrBadBirthDate
		}
	case KindExpat:
		// @TODO(kl): find the proper number
		if len(pq.ID) <= 3 && strings.HasPrefix(pq.ID, expatPrefix) {
			return ErrBadIqamaID
		}
		if len(pq.BirthDate) <= 4 {
			return ErrBadBirthDate
		}
	default:
		return ErrUnknownPatientType
	}

	return nil
}

// Validate checks if ID is valid , depending on the type
func (pq *PatientQuery) ValidateID() error {
	switch pq.Kind() {
	case KindCitizen:
		if len(pq.ID) != 10 && strings.HasPrefix(pq.ID, citizenPrefix) {
			return ErrBadNationalID
		}
	case KindExpat:
		// @TODO(kl): find the proper number
		if len(pq.ID) <= 3 && strings.HasPrefix(pq.ID, expatPrefix) {
			return ErrBadIqamaID
		}
	default:
		return ErrUnknownPatientType
	}

	return nil
}

// Controller handles the logic of:
// - Getting Patient and iff missing call Yakeen, then add it to our db
// - Getting Establishments
// - Getting Practitioners
//
// The goal from having the logic here,
// is to keep the business logic away from implementation details like http routes
type Controller struct {
	store    *store.Store
	yakeen   *yakeen.Yakeen
	nic      *nic.Nic
	Sc       *scfhs.Scfhs
	features []string
}

// New returns an instance of Controller
// configs are define in package config
func New(s *store.Store, conf *config.Config) (*Controller, error) {
	// init yakeen
	yak, err := yakeen.New(conf.Gateway.Token, conf.Gateway.URL)
	if err != nil {
		return nil, err
	}

	//init oauth
	oauth, err := oauth.New(conf.Oauth.Consumers, conf.Oauth.DBPath, conf.Gateway.URL)
	if err != nil {
		return nil, err
	}

	//init nic
	n, err := nic.New(conf.Nic.CallerID, conf.Gateway.URL, oauth)
	if err != nil {
		return nil, err
	}

	//init Scfhs
	sc, err := scfhs.New(conf.Scfhs.URL, conf.Scfhs.Token)
	if err != nil {
		return nil, err
	}

	cont := &Controller{
		store:    s,
		yakeen:   yak,
		nic:      n,
		Sc:       sc,
		features: conf.Features,
	}
	return cont, nil
}

// GetPatient talks to store.GetPatient if not found it then calls Yakeen
// it stores the results in the downstream db
// then returns
func (c *Controller) GetPatient(pq *PatientQuery) (*store.Patient, error) {
	id := pq.ID
	pnt, err := c.store.GetPatient(id, pq.BirthDate)
	if err != nil && err == store.ErrSearch {
		log.Println(err)
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}

	// patient found nothing to do
	if pnt != nil && err != store.ErrNotFound {
		return pnt, nil
	}

	if c.FeatureIsEnabled("disable-yakeen") {
		return nil, store.ErrNotFound
	}

	// prepare birthDate based on patient type
	// hijri for citizens, gregorian for expats
	switch pq.Kind() {
	case KindCitizen:
		if pnt.DateH != nil {
			pq.BirthDate = *pnt.DateH
		}
		fmt.Println(pq)
	case KindExpat:
		if pnt.DateG != nil {
			pq.BirthDate = *pnt.DateG
		}
		fmt.Println(pq)
	}

	err = c.getPnt(pq, pnt)
	if err != nil && (err == yakeen.ErrBadDOB || err == yakeen.ErrBadID) {
		log.Println(err)
		return nil, ErrSearchInput
	} else if err != nil {
		log.Println(err)
		// avoid leaking sensitive info
		return nil, ErrFetchingInfo
	}

	// compute patient age
	pnt.Age = c.calcAge(pnt.DateG)

	// get nationality iso code
	country, _ := c.store.GetCountryIsoCode(pnt.Nationality)
	if country != nil {
		pnt.NationalityCode = country.Code
		pnt.Nationality = country.CountryNameEn
	}

	// add to db
	go c.addPatient(pnt)

	return pnt, nil
}

//GetPatientByID get patient from db
func (c *Controller) GetPatientByID(id string) (*store.Patient, error) {
	pnt, err := c.store.GetPatientByID(id)
	if err != nil && err == store.ErrSearch {
		log.Println(err)
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}

	// patient not found
	if pnt == nil && err == store.ErrNotFound {
		return nil, ErrNotFound
	}

	return pnt, nil
}

// UpdatePatient calls Yakeen to updated info and update it in
// in the downstream db
// then returns
func (c *Controller) UpdatePatient(pq *PatientQuery) (*store.Patient, error) {
	id := pq.ID
	pnt, err := c.store.GetPatientByID(id)
	if err != nil && err == store.ErrSearch {
		log.Println(err)
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}

	// Getting the date from the database instead of user input,,, Caused an issue with some formatting and mismatching dates
	// prepare birthDate based on patient type
	// hijri for citizens, gregorian for expats
	// switch pq.Kind() {
	// case KindCitizen:
	// 	// if the record is found "DateOfBirthG" will be returned
	// 	// if the record was found "DateG" will be returned
	// 	if pnt.DateH != nil {
	// 		pq.BirthDate = *pnt.DateH
	// 	} else if pnt.DateOfBirthH != nil {
	// 		pq.BirthDate = *pnt.DateOfBirthH
	// 	}
	// case KindExpat:
	// 	if pnt.DateG != nil {
	// 		pq.BirthDate = *pnt.DateG
	// 	} else if pnt.DateOfBirthG != nil {
	// 		pq.BirthDate = *pnt.DateOfBirthG
	// 	}
	// }

	err = c.getPnt(pq, pnt)
	if err != nil && (err == yakeen.ErrBadDOB || err == yakeen.ErrBadID) {
		log.Println(err)
		return nil, ErrSearchInput
	} else if err != nil {
		log.Println(err)
		// avoid leaking sensitive info
		return nil, ErrFetchingInfo
	}

	// compute patient age
	pnt.Age = c.calcAge(pnt.DateG)

	// add to db
	err = c.store.UpdatesPatient(pnt)
	if err != nil {
		log.Println(err)
		return nil, ErrUpdateInfo
	}

	return nil, nil
}

func (c *Controller) GetFullPatientInfo(pq *PatientQuery) (*store.Patient, error) {
	pnt := &store.Patient{}
	err := c.getFullPnt(pq, pnt)
	if err != nil && err == nic.ErrValidation {
		return nil, ErrSearchInput
	} else if err != nil {
		log.Println(err)
		// avoid leaking sensitive info
		return nil, ErrFetchingInfo
	}
	// compute patient age
	pnt.Age = c.calcAge(pnt.DateOfBirthG)
	go c.addPatient(pnt)
	return pnt, nil
}

func (c *Controller) calcAge(birthDate *string) *string {
	if birthDate == nil {
		return birthDate
	}
	bdate := *birthDate
	if len(bdate) <= 0 {
		return birthDate
	}
	t, err := time.Parse("02-01-2006", bdate)
	if err != nil {
		log.Println(err)
		return nil
	}

	diff := time.Since(t).Hours() / hoursPerYear
	age := strconv.Itoa(int(diff))
	return &age
}

// format birthDateH to yakeen from yyyy-mm-dd to dd-mm-yyyy
func (c *Controller) formatBirthDate(birthDate string) string {
	d, err := time.Parse("2006-01-02", birthDate)
	if err != nil {
		fmt.Println(err)

		if strings.Count(birthDate, "/") >= 2 {
			res := strings.Replace(birthDate, "/", "-", 2)
			birthDate = res
		}

		return birthDate
	}

	return d.Format("02-01-2006")
}
func (c *Controller) getPnt(pq *PatientQuery, pnt *store.Patient) error {
	// fetch patient from yakeen since it's not found
	switch pq.Kind() {
	case KindCitizen:
		ctzn, err := c.yakeen.GetCitizen(pq.ID, c.formatBirthDate(pq.BirthDate))
		if err != nil {
			return err
		}
		err = c.ctznToPnt(ctzn, pnt)
		if err != nil {
			return err
		}
	case KindExpat:
		exp, err := c.yakeen.GetExpat(pq.ID, c.formatBirthDate(pq.BirthDate))
		if err != nil {
			return err
		}
		err = c.expatToPtnt(exp, pnt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) getFullPnt(pq *PatientQuery, pnt *store.Patient) error {
	// fetch patient from nic since it's not found
	p, err := c.nic.GetPatient(pq.ID)
	if err != nil {
		return err
	}
	pnt.IDNumber = &p.ID
	pnt.Gender = &p.Gender
	pnt.MobileNumber = &p.MobileNumber

	pnt.NationalityCode = &p.NationalityCode
	pnt.Nationality = &p.NationalityDescAr

	pnt.OccupationCode = &p.OccupationCode
	pnt.Occupation = &p.OccupationDescAr

	pnt.FirstNameEn = &p.FirstNameEn
	pnt.SecondNameEn = &p.SecondNameEn
	pnt.ThirdNameEn = &p.ThirdNameEn
	pnt.LastNameEn = &p.LastNameEn

	pnt.FirstNameAr = &p.FirstNameAr
	pnt.SecondNameAr = &p.SecondNameAr
	pnt.ThirdNameAr = &p.ThirdNameAr
	pnt.LastNameAr = &p.LastNameAr

	// format the data return from nic 1963-07-21T00:00:00 to date
	d, err := time.Parse("2006-01-02T15:04:05", p.BirthDateG)
	if err != nil {
		return err
	}
	date := d.Format("02-01-2006")
	pnt.DateOfBirthG = &date

	// add the patient's gregorian birth date from SQL Server
	pnt.DateOfBirthH = pnt.DateH

	switch pq.Kind() {
	case KindCitizen:
		ty := "NationalId"
		pnt.IDType = &ty
	case KindExpat:
		ty := "Iqama"
		pnt.IDType = &ty
	}

	return nil
}

func (c *Controller) Convert(pq *PatientQuery, pnt *store.Patient) ([]byte, error) {
	//convert age to string
	age := 0
	if pnt.Age != nil {
		a, err := strconv.Atoi(*pnt.Age)
		if err != nil {
			log.Println(err)
		}
		age = a
	}

	switch pq.Kind() {
	case KindCitizen:
		cc := CompatCitizen{
			HealthID:          c.SetDefaultValue(pnt.HealthID, nil),
			IDType:            c.SetDefaultValue(pnt.IDType, nil),
			IDExpiryDate:      c.SetDefaultValue(pnt.IDExpiryDate, nil),
			IDNumber:          c.SetDefaultValue(&pq.ID, nil),
			DateOfBirth:       c.SetDefaultValue(&pq.BirthDate, nil),
			Age:               &age,
			PlaceOfBirth:      c.SetDefaultValue(pnt.PlaceOfBirth, nil),
			EnglishFirstName:  c.SetDefaultValue(pnt.FirstNameEn, nil),
			EnglishSecondName: c.SetDefaultValue(pnt.SecondNameEn, nil),
			EnglishThirdName:  c.SetDefaultValue(pnt.ThirdNameEn, nil),
			EnglishLastName:   c.SetDefaultValue(pnt.LastNameEn, nil),
			FirstName:         c.SetDefaultValue(pnt.FirstNameAr, nil),
			FatherName:        c.SetDefaultValue(pnt.SecondNameAr, nil),
			GrandFatherName:   c.SetDefaultValue(pnt.ThirdNameAr, nil),
			FamlyName:         c.SetDefaultValue(pnt.LastNameAr, nil),
			SubtribeName:      c.SetDefaultValue(pnt.SubtribeName, nil),
			Gender:            c.SetDefaultValue(pnt.Gender, nil),
			Nationality:       c.SetDefaultValue(pnt.Nationality, &SANationality),
			NationalityCode:   c.SetDefaultValue(pnt.NationalityCode, &SACountryCode),
			OccupationCode:    c.SetDefaultValue(pnt.OccupationCode, nil),
			MaritalStatus:     c.SetDefaultValue(pnt.MaritalStatus, &UnknownStatus),
			MaritalStatusCode: c.SetDefaultValue(pnt.MaritalStatusCode, &UnknownStatusCode),
			PatientStatus:     c.SetDefaultValue(pnt.PatientStatus, nil),
		}
		return json.Marshal(cc)
	case KindExpat:
		ity := "Iqama"
		ce := CompatExpat{
			HealthID:     c.SetDefaultValue(pnt.HealthID, nil),
			IDType:       c.SetDefaultValue(&ity, nil),
			IDNumber:     c.SetDefaultValue(&pq.ID, nil),
			DateOfBirth:  c.SetDefaultValue(&pq.BirthDate, nil),
			PlaceOfBirth: c.SetDefaultValue(pnt.PlaceOfBirth, nil),

			IDExpiryDate:    c.SetDefaultValue(pnt.IDExpiryDate, nil),
			FirstName:       c.SetDefaultValue(pnt.FirstNameAr, nil),
			FatherName:      c.SetDefaultValue(pnt.SecondNameAr, nil),
			GrandFatherName: c.SetDefaultValue(pnt.ThirdNameAr, nil),
			FamlyName:       c.SetDefaultValue(pnt.LastNameAr, nil),

			EnglishFirstName:  c.SetDefaultValue(pnt.FirstNameEn, nil),
			EnglishSecondName: c.SetDefaultValue(pnt.SecondNameEn, nil),
			EnglishThirdName:  c.SetDefaultValue(pnt.ThirdNameEn, nil),
			EnglishLastName:   c.SetDefaultValue(pnt.LastNameEn, nil),

			Gender:          c.SetDefaultValue(pnt.Gender, nil),
			Nationality:     c.SetDefaultValue(pnt.Nationality, nil),
			NationalityCode: c.SetDefaultValue(pnt.NationalityCode, nil),
			OccupationCode:  c.SetDefaultValue(pnt.Occupation, nil),
			Age:             &age,
		}
		return json.Marshal(ce)
	}
	return nil, ErrUnknownPatientType
}

func (c *Controller) Convertv2(pq *PatientQuery, pnt *store.Patient) ([]byte, error) {
	//convert age to string
	age := 0
	if pnt.Age != nil {
		a, err := strconv.Atoi(*pnt.Age)
		if err != nil {
			log.Println(err)
		}
		age = a
	}
	switch pq.Kind() {
	case KindCitizen:
		cc := CompatCitizenv2{
			HealthID:          c.SetDefaultValue(pnt.HealthID, nil),
			IDType:            c.SetDefaultValue(pnt.IDType, nil),
			IDNumber:          c.SetDefaultValue(&pq.ID, nil),
			DateOfBirth:       c.SetDefaultValue(pnt.DateOfBirthG, nil),
			Age:               &age,
			PlaceOfBirth:      c.SetDefaultValue(pnt.PlaceOfBirth, nil),
			EnglishFirstName:  c.SetDefaultValue(pnt.FirstNameEn, nil),
			EnglishSecondName: c.SetDefaultValue(pnt.SecondNameEn, nil),
			EnglishThirdName:  c.SetDefaultValue(pnt.ThirdNameEn, nil),
			EnglishLastName:   c.SetDefaultValue(pnt.LastNameEn, nil),
			FirstNameAr:       c.SetDefaultValue(pnt.FirstNameAr, nil),
			SecondNameAr:      c.SetDefaultValue(pnt.SecondNameAr, nil),
			ThirdNameAr:       c.SetDefaultValue(pnt.ThirdNameAr, nil),
			LastNameAr:        c.SetDefaultValue(pnt.LastNameAr, nil),
			Gender:            c.SetDefaultValue(pnt.Gender, nil),
			Nationality:       c.SetDefaultValue(pnt.Nationality, &SANationality),
			NationalityCode:   c.SetDefaultValue(pnt.NationalityCode, &SACountryCode),
			Occupation:        c.SetDefaultValue(pnt.Occupation, nil),
			MaritalStatus:     c.SetDefaultValue(pnt.MaritalStatus, &UnknownStatus),
			MaritalStatusCode: c.SetDefaultValue(pnt.MaritalStatusCode, &UnknownStatusCode),
			PatientStatus:     c.SetDefaultValue(pnt.PatientStatus, nil),
		}
		return json.Marshal(cc)
	case KindExpat:
		ity := "Iqama"
		ce := CompatExpatv2{
			HealthID:     c.SetDefaultValue(pnt.HealthID, nil),
			IDType:       c.SetDefaultValue(&ity, nil),
			IDNumber:     c.SetDefaultValue(&pq.ID, nil),
			DateOfBirth:  c.SetDefaultValue(pnt.DateOfBirthG, nil),
			PlaceOfBirth: c.SetDefaultValue(pnt.PlaceOfBirth, nil),

			FirstNameAr:  c.SetDefaultValue(pnt.FirstNameAr, nil),
			SecondNameAr: c.SetDefaultValue(pnt.SecondNameAr, nil),
			ThirdNameAr:  c.SetDefaultValue(pnt.ThirdNameAr, nil),
			LastNameAr:   c.SetDefaultValue(pnt.LastNameAr, nil),

			EnglishFirstName:  c.SetDefaultValue(pnt.FirstNameEn, nil),
			EnglishSecondName: c.SetDefaultValue(pnt.SecondNameEn, nil),
			EnglishThirdName:  c.SetDefaultValue(pnt.ThirdNameEn, nil),
			EnglishLastName:   c.SetDefaultValue(pnt.LastNameEn, nil),

			Gender:            c.SetDefaultValue(pnt.Gender, nil),
			Nationality:       c.SetDefaultValue(pnt.Nationality, nil),
			NationalityCode:   c.SetDefaultValue(pnt.NationalityCode, nil),
			Occupation:        c.SetDefaultValue(pnt.Occupation, nil),
			MaritalStatus:     c.SetDefaultValue(pnt.MaritalStatus, &UnknownStatus),
			MaritalStatusCode: c.SetDefaultValue(pnt.MaritalStatusCode, &UnknownStatusCode),
			PatientStatus:     c.SetDefaultValue(pnt.PatientStatus, nil),
			Age:               &age,
		}
		return json.Marshal(ce)
	}
	return nil, ErrUnknownPatientType
}

func (c *Controller) ctznToPnt(ctzn *yakeen.Citizen, pnt *store.Patient) error {
	ty := "NationalId"
	pnt.IDType = &ty
	pnt.IDNumber = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.NationalID
	pnt.IDIssuePlace = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.IDIssuePlace
	pnt.IDIssueDate = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.IDIssueDate
	pnt.IDExpiryDate = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.IDExpiryDate
	pnt.DateOfBirthH = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.BirthDate

	// add the patient's gregorian birth date from SQL Server
	pnt.DateOfBirthG = pnt.DateG

	// pnt.OccupationCode = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.

	pnt.FirstNameEn = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.EnglishFirstName
	pnt.SecondNameEn = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.EnglishSecondName
	pnt.ThirdNameEn = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.EnglishThirdName
	pnt.LastNameEn = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.EnglishLastName

	pnt.FirstNameAr = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.FirstName
	pnt.SecondNameAr = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.FatherName
	pnt.ThirdNameAr = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.GrandFatherName
	pnt.SubtribeName = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.SubtribeName
	pnt.LastNameAr = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.FamilyName

	pnt.Gender = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.Gender
	pnt.PlaceOfBirth = &ctzn.GetCitizenInfoResponse.CitizenInfoResult.PlaceOfBirth

	return nil
}

func (c *Controller) expatToPtnt(expt *yakeen.Expat, pnt *store.Patient) error {
	ty := "Iqama"
	pnt.IDType = &ty

	pnt.IDNumber = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.IqamaID
	pnt.IDIssuePlace = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.IqamaIssuePlaceDesc
	pnt.IDIssueDate = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.IqamaIssueDateH
	pnt.IDExpiryDate = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.IqamaExpiryDateH
	pnt.DateOfBirthH = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.BirthDate

	// add the patient's gregorian birth date from SQL Server
	pnt.DateOfBirthG = pnt.DateG

	pnt.Nationality = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.NationalityDesc
	pnt.Occupation = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.OccupationDesc
	pnt.Gender = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.Gender

	pnt.FirstNameEn = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.EnglishFirstName
	pnt.SecondNameEn = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.EnglishSecondName
	pnt.ThirdNameEn = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.EnglishThirdName
	pnt.LastNameEn = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.EnglishLastName

	pnt.FirstNameAr = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.FirstName
	pnt.SecondNameAr = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.SecondName
	pnt.ThirdNameAr = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.ThirdName
	pnt.LastNameAr = &expt.GetAlienInfoByIqamaResponse.AlienInfoByIqamaResult.LastName

	return nil
}

func (c *Controller) addPatient(pnt *store.Patient) {
	log.Println(c.store.AddPatient(pnt))
}

// GetEstablishment searches for a *store.Establishment by id and returns it
func (c *Controller) GetEstablishment(id string) (*store.Establishment, error) {
	est, err := c.store.GetEstablishment(id)
	if err != nil && err == store.ErrSearch {
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}
	return est, nil
}

// GetEstablishment searches for a *store.Establishment by id and returns it
func (c *Controller) GetEstablishmentV2(id string) (*store.EstablishmentV2, error) {
	est, err := c.store.GetEstablishmentV2(id)
	if err != nil && err == store.ErrSearch {
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}
	return est, nil
}

// GetEstablishments get full tEstablishments list and returns it
func (c *Controller) GetEstablishments() (*[]store.Establishments, error) {
	est, err := c.store.GetEstablishments()
	if err != nil && err == store.ErrSearch {
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}
	return est, nil
}

// GetEstablishments get full tEstablishments list and returns it
func (c *Controller) GetEstablishmentsV2() (*[]store.EstablishmentV2, error) {
	est, err := c.store.GetEstablishmentsV2()
	if err != nil && err == store.ErrSearch {
		return nil, ErrSearchInput
	} else if err != nil && err != store.ErrNotFound {
		// avoid leaking sensitive info
		log.Println(err)
		return nil, ErrLookingUpInfo
	}
	return est, nil
}

func (c *Controller) GetPractitioner(id string) (*store.Practitioner, error) {
	pract, err := c.store.GetPractitioner(id)
	if err != nil && err != store.ErrNotFound {
		log.Println(err)
		return nil, ErrLookingUpInfo
	}

	// Practitioner found nothing to do
	if pract != nil && err != store.ErrNotFound {
		return pract, nil
	}
	// Uncomment this to disable SCHFS
	// if c.FeatureIsEnabled("disable-scfhs") {
	// 	return nil, store.ErrNotFound
	// }

	if err := c.getPract(id, pract); err != nil {
		return nil, err
	}

	// Adds the record to DB from SCHFS
	c.store.AddPractitioner(pract)

	// Get the values from the DB after it process the data.
	pract, err = c.store.GetPractitioner(id)
	if err != nil && err != store.ErrNotFound {
		log.Println(err)
		return nil, ErrLookingUpInfo
	}

	return pract, nil
}

// fetch Practitioner from Scfhs since it's not found
func (c *Controller) getPract(id string, pract *store.Practitioner) error {
	p, err := c.Sc.GetPractitioner(id)
	if err != nil {
		return err
	}

	pract.SCFHSRegistrationNumber = &p.Response.Info.Profile.RegistrationNumber
	pract.PractitionerID = pract.HealthID
	pract.IDNumber = &id

	pract.FirstNameAr = &p.Response.Info.Profile.Ar.FirstName
	pract.SecondNameAr = &p.Response.Info.Profile.Ar.SecondName
	pract.LastNameAr = &p.Response.Info.Profile.Ar.LastName

	pract.FirstNameEn = &p.Response.Info.Profile.En.FirstName
	pract.SecondNameEn = &p.Response.Info.Profile.En.SecondName
	pract.LastNameEn = &p.Response.Info.Profile.En.LastName

	pract.Gender_code = &p.Response.Info.Profile.Gender.Code
	pract.Gender_ar = &p.Response.Info.Profile.Gender.NameAr
	pract.Gender_en = &p.Response.Info.Profile.Gender.NameEn

	pract.SCFHSCategoryCode = &p.Response.Info.Professionality.Category.Code
	pract.SCFHSCategoryAr = &p.Response.Info.Professionality.Category.NameAr
	pract.SCFHSCategoryEn = &p.Response.Info.Professionality.Category.NameEn

	pract.SCFHSSpecialityCode = &p.Response.Info.Professionality.Specialty.Code
	pract.SCFHSSpecialityAr = &p.Response.Info.Professionality.Specialty.NameAr
	pract.SCFHSSpecialityEn = &p.Response.Info.Professionality.Specialty.NameEn

	pract.SCFHSPractitionerStatus = &p.Response.Info.Status.Code
	pract.SCFHSPractitionerStatusCode = &p.Response.Info.Status.DescAr

	pract.SCFHSRegistrationIssueDate = &p.Response.Info.Status.License.IssuedDate
	pract.SCFHSRegistrationExpiryDate = &p.Response.Info.Status.License.ExpiryDate

	// @TODO: change this
	pq := PatientQuery{ID: id}
	switch pq.Kind() {
	case KindCitizen:
		ty := "NationalId"
		pract.IDType = &ty
	case KindExpat:
		ty := "Iqama"
		pract.IDType = &ty
	}
	return nil
}

// SetDefaultValue take var and it is default value
func (c *Controller) SetDefaultValue(s *string, d *string) *string {
	if s == nil && d != nil {
		return d
	} else if s == nil {
		// if default value is null
		return &generalDefaultValue
	}
	return s
}

func (c *Controller) UpdateEstablishment(est *store.Establishment) error {
	if err := c.store.UpdateGovEstablishment(est); err != nil {
		log.Println("establishmentUpdate error:", err)
		return ErrUpdateInfo
	}
	return nil
}

// return true if feature in features list
func (c *Controller) FeatureIsEnabled(feature string) bool {
	sort.Strings(c.features)
	i := sort.SearchStrings(c.features, feature)
	if i < len(c.features) && c.features[i] == feature {
		return true
	}
	return false
}

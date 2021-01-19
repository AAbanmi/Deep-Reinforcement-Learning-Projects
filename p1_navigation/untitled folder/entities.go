package store

import "errors"

const ErrorMsg = "NotFound"

var (
	ErrNotFound = errors.New("row not found")
	ErrSearch   = errors.New("search input error")
	ErrUpdate   = errors.New("encountered error while updating")
)

type Practitioner struct {
	// these fields are returned when there's no match for the patient on our db
	LogId            *string `json:"log_id,omitempty"`
	ErrorMsg         *string `json:"error_msg,omitempty"`
	Msg              *string `json:"msg,omitempty" db:"Msg"`
	ReservedHealthID *string `json:"reserved_health_id,omitempty" db:"ReservedHealthID"`
	HealthID         *string `json:"health_id,omitempty" db:"Practitioner_id"`
	SearchID         *string `json:"search_id,omitempty" db:"SearchID"`

	ID             int         `json:"-" db:"id"`
	PractitionerID *string     `json:"-" db:"Practitioner_id"`
	Address        interface{} `db:"address" json:"address,omitempty"`
	RowInseartedAt *string     `json:"-" db:"RowInseartedAt"`
	RowUpdatedAt   *string     `json:"row_updated_at" db:"RowUpdatedAt"`
	RowDeletedAt   *string     `json:"row_deleted_at,omitempty" db:"RowDeletedAt"`
	IsDelted       *string     `json:"-" db:"IsDeleted"`

	FirstNameAr  *string `json:"first_name_ar,omitempty" db:"FirstNameAr"`
	SecondNameAr *string `json:"second_name_ar,omitempty"  db:"SecondNameAr"`
	ThirdNameAr  *string `json:"third_name_ar,omitempty"  db:"ThirdNameAr"`
	LastNameAr   *string `json:"last_name_ar,omitempty"  db:"LastNameAr"`
	FullNameAr   *string `json:"full_name_ar,omitempty"`

	FirstNameEn  *string `json:"first_name_en,omitempty" db:"FirstNameEn"`
	SecondNameEn *string `json:"second_name_en,omitempty" db:"SecondNameEn"`
	ThirdNameEn  *string `json:"third_name_en,omitempty" db:"ThirdNameEn"`
	LastNameEn   *string `json:"last_name_en,omitempty" db:"LastNameEn"`
	FullNameEn   *string `json:"full_name_en,omitempty"`

	IDIssueDate        *string `json:"id_issue_date,omitempty"`
	BirthDate_original *string `json:"birth_date_original"`
	BirthDate_G        *string `json:"birth_date_gregorian"`
	BirthDate_H        *string `json:"birth_date_hirji"`

	Gender_ar   *string `json:"gender_ar,omitempty" db:"Gender_ar"`
	Gender_en   *string `json:"gender_en,omitempty" db:"Gender_en"`
	Gender_code *string `json:"gender_code,omitempty" db:"Gender_code"`

	Nationality_ar   *string `json:"nationality_ar"`
	Nationality_en   *string `json:"nationality_en"`
	Nationality_code *string `json:"nationality_code"`

	Contract_Type               *string `json:"contract_type,omitempty"`
	Phone                       *string `json:"phone,omitempty"`
	Email                       *string `json:"email,omitempty"`
	Department                  *string `json:"department,omitempty"`
	InsuranceCompany            *string `json:"insurance_company,omitempty"`
	InsuranceExpiryDate         *string `json:"insurance_expiry_date,omitempty"`
	EstablishmentName           *string `json:"establishment_name" db:"establishment_name"`
	EstablishmentType           *string `json:"establishment_type,omitempty" db:"establishment_type"`
	EstablishmentSector         *string `json:"establishment_sector,omitempty" db:"establishment_sector"`
	EstablishmentOrgID          *string `json:"establishment_org_id" db:"establishment_org_id"`
	HLSPractitionersID          *string `json:"-" db:"hls_practitioners_id"`
	HLSPractitionerLicensesId   *string `json:"-" db:"hls_practitioner_licenses_id"`
	HLSEstablishmentID          *string `json:"-" db:"hls_establishment_id"`
	HLSEstablishmentLicenseID   *string `json:"-" db:"hls_establishment_license_id"`
	SourceSystem                *string `json:"-" db:"SourceSystem"`
	Gender                      *string `json:"gender,omitempty"`
	Nationality                 *string `json:"nationality,omitempty"`
	Religion                    *string `json:"religion"`
	BirthDateH                  *string `json:"birth_date_h,omitempty"`
	BirthDateG                  *string `json:"birth_date_g,omitempty"`
	License                     *string `json:"license,omitempty"`
	LicenseIssueDate            *string `json:"license_issue_date,omitempty" db:"LicenseIssueDate"`
	LicenseExpiryDate           *string `json:"license_expiry_date,omitempty" db:"LicenseExpiryDate"`
	SCFHSRegistrationIssueDate  *string `json:"scfhs_registration_issue_date,omitempty" db:"SCFHSRegistrationIssueDate"`
	SCFHSRegistrationExpiryDate *string `json:"scfhs_registration_expiry_date,omitempty" db:"SCFHSRegistrationExpiryDate"`
	SCFHSRegistrationNumber     *string `json:"scfhs_registration_number,omitempty" db:"SCFHSRegistrationNumber"`

	IDType           *string `json:"id_type,omitempty" db:"IDType"`
	IDNumber         *string `json:"id_number,omitempty" db:"IDNumber"`
	ExpirationStatus *string `json:"expiration_status,omitempty" db:"ExpirationStatus"`
	IDExpiryDate     *string `json:"id_expiry_date,omitempty"`

	LicenseNumber *string `json:"license_number"`
	Legacy_job    *string `json:"-"`

	SCFHSCategoryCode *string `json:"scfhs_category_code,omitempty" db:"SCFHSCategoryCode"`
	SCFHSCategoryAr   *string `json:"scfhs_category_ar,omitempty" db:"SCFHSCategoryAr"`
	SCFHSCategoryEn   *string `json:"scfhs_category_en,omitempty" db:"SCFHSCategoryEn"`

	SCFHSSpecialityCode *string `json:"scfhs_speciality_code,omitempty" db:"SCFHSSpecialityCode"`
	SCFHSSpecialityAr   *string `json:"scfhs_speciality_ar,omitempty" db:"SCFHSSpecialityAr"`
	SCFHSSpecialityEn   *string `json:"scfhs_speciality_en,omitempty" db:"SCFHSSpecialityEn"`

	SCFHSPractitionerStatus     *string `json:"scfhs_practitioner_status,omitempty" db:"SCFHSPractitionerStatus"`
	SCFHSPractitionerStatusCode *string `json:"scfhs_practitioner_status_code,omitempty"  db:"SCFHSPractitionerStatusCode"`

	LegacyQualification *string `json:"legacy_qualification,omitempty"`
	InsuranceNumber     *string `json:"insurance_number,omitempty"`
	UniversityFaculty   *string `json:"university_faculty"`
}

type Patient struct {
	// these fields are returned when there's no match for the patient on our db
	LogId            *string `json:"log_id,omitempty"`
	ErrorMsg         *string `json:"error_msg,omitempty"`
	Msg              *string `json:"msg,omitempty" db:"Msg"`
	ReservedHealthID *string `json:"reserved_health_id,omitempty"`
	HealthID         *string `json:"health_id,omitempty" db:"HealthId"`
	SearchID         *string `json:"search_id,omitempty" db:"SearchID"`
	DateG            *string `json:"date_g,omitempty" db:"DateG"`
	DateH            *string `json:"date_h,omitempty" db:"DateH"`

	//CamelCase is fine
	ClientIdentifierId *string `json:"ClientIdentifierId,omitempty" db:"ClientIdentifierId"`
	IDType             *string `json:"id_type,omitempty" db:"IdType"`
	IDNumber           *string `json:"id_number,omitempty" db:"IdNumber"`
	IDExpiryDate       *string `json:"id_expiry_date,omitempty" db:"IDExpiryDate"`
	IDIssueDate        *string `json:"id_issue_date,omitempty" db:"IDIssueDate"`
	IDIssuePlace       *string `json:"id_issue_place,omitempty" db:"ID_Place"`

	BloodType *string `json:"blood_type,omitempty"`

	Age *string `json:"age,omitempty" db:"Age"`

	// used for citizens only
	OccupationCode *string `json:"occupation_code,omitempty" db:"OccupationCode"`
	// used for expats
	Occupation *string `json:"occupation,omitempty" db:"Occupation"`

	Nationality     *string `json:"nationality_ar,omitempty" db:"Nationality"`
	NationalityCode *string `json:"nationality_en,omitempty" db:"NationalityCode"`

	MaritalStatus     *string `json:"marital_status,omitempty" db:"MaritalStatus"`
	MaritalStatusCode *string `json:"marital_status_code,omitempty" db:"MaritalStatusCode"`

	PatientStatus   *string `json:"patient_status,omitempty" db:"PatientStatus"`
	TransactionID   *string `json:"transaction_id,omitempty" db:"TransactionID"`
	GenderSpecified *string `json:"gender_specified,omitempty" db:"GenderSpecified"`

	HifizaIssueDate *string `json:"hifiza_issue_date,omitempty" db:"HifizaIssueDate"`
	HifizaNumber    *string `json:"hifiza_number,omitempty" db:"HifizaNumber"`

	PlaceOfBirth *string `json:"place_of_birth,omitempty" db:"PlaceOfBirth"`
	DateOfBirthG *string `json:"date_of_birth_g,omitempty" db:"DateOfBirthG"`
	DateOfBirthH *string `json:"date_of_birth_h,omitempty" db:"DateOfBirthH"`

	Gender *string `json:"gender,omitempty" db:"Gender"`

	FirstNameAr  *string `json:"first_name_ar,omitempty" db:"FirstName"`
	SecondNameAr *string `json:"second_name_ar,omitempty" db:"FatherName"`
	ThirdNameAr  *string `json:"third_name_ar,omitempty" db:"GrandFatherName"`
	LastNameAr   *string `json:"last_name_ar,omitempty" db:"FamilyName"`
	SubtribeName *string `json:"subtribe_name,omitempty" db:"SubtribeName"`

	FirstNameEn  *string `json:"first_name_en,omitempty" db:"EnglishFirstName"`
	SecondNameEn *string `json:"second_name_en,omitempty" db:"EnglishSecondName"`
	ThirdNameEn  *string `json:"third_name_en,omitempty" db:"EnglishThirdName"`
	LastNameEn   *string `json:"last_name_en,omitempty" db:"EnglishLastName"`

	IsDead        *string `json:"is_dead,omitempty" db:"IsDead"`
	SponsorNumber *string `json:"sponsor_number,omitempty" db:"SponsorNumber"`
	MobileNumber  *string `json:"mobile_number,omitempty" db:"MobileNumber"`
	PhoneNumber   *string `json:"phone_number,omitempty"`
	EmailAddress  *string `json:"email_address,omitempty"`
}

type Establishment struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
	Msg      *string `json:"msg,omitempty" db:"Msg"`

	EntityType       *string `json:"entity_type,omitempty"`
	EntitySpeciality *string `json:"entity_speciality,omitempty"`

	ID             *string `json:"id" db:"id"`
	Code           *string `json:"code" db:"Code"`
	MOHID          *string `json:"moh_id" db:"MOHID"`
	OrganizationID *string `json:"organization_id" db:"OrganizationId"`
	LegacyEntityID *string `json:"legacy_entity_id" db:"legacy_entityId"`

	RowInsertedAt *string `json:"created_at" db:"RowInseartedAt"`
	RowUpdatedAt  *string `json:"updated_at" db:"RowUpdatedAt"`
	RowDeletedAt  *string `json:"deleted_at" db:"RowDeletedAt"`
	IsMigrated    *string `json:"migrated" db:"isMigrated"`
	IsDeleted     *string `json:"deleted" db:"IsDeleted"`

	LicenseNumber         *string `json:"license_number,omitempty" db:"LicenseNumber"`
	IssueDate             *string `json:"issue_date,omitempty" db:"Issue_Date"`
	ExpiryDate            *string `json:"expiry_date,omitempty" db:"Expiry_Date"`
	NameAr                *string `json:"name_ar,omitempty"`
	NameEn                *string `json:"name_en,omitempty"`
	SehaID                *string `json:"seha_id,omitempty"`
	BedsCount             *string `json:"beds_count,omitempty"`
	Longitude             *string `json:"longitude,omitempty"`
	Latitude              *string `json:"latitude,omitempty"`
	MapURL                *string `json:"map_url,omitempty" db:"MapUrl"`
	SehaHealthDirectory   *string `json:"seha_health_directory,omitempty" db:"SehaHealth_Directory"`
	LevelOfCare           *string `json:"level_of_care" db:"level_of_care"`
	TypeOfCare            *string `json:"type_of_care" db:"type_of_care"`
	TypeOfCareCode        *string `json:"type_of_care_code" db:"type_of_care_code"`
	HealthDirectoryAr     *string `json:"health_directory_ar" db:"Health_Directory_ar"`
	HealthDirectoryEn     *string `json:"health_directory_en" db:"Health_Directory_en"`
	HealthDirectorySehaID *string `json:"health_directory_seha_id" db:"Health_Directory_SehaID"`
	SectorAr              *string `json:"sector_ar" db:"sector_ar"`
	SectorEn              *string `json:"sector_en" db:"sector_en"`
	SectorCode            *string `json:"sector_code" db:"sector_code"`
	EntityTypeAr          *string `json:"entity_type_ar" db:"entity_type_ar"`
	EntityTypeEn          *string `json:"entity_type_en" db:"entity_type_en"`
	EntityTypeCode        *string `json:"entity_type_code" db:"entity_type_code"`
	NotificationEmail     *string `json:"notification_email,omitempty"`
	Website               *string `json:"website" db:"webSite"`
	PhoneNumber           *string `json:"phone_number" db:"phone_number"`
	Email                 *string `json:"email" db:"email"`

	OldHLSEntityType     *string `json:"old_hls_entity_type" db:"OldHlsEntityType"`
	OldHLSEntityTypeID   *string `json:"old_hls_entity_id" db:"OldHlsEntityTypeId"`
	OldHLSSpeciality     *string `json:"old_hls_speciality" db:"OldHlsSpeciality"`
	OldHLSSpecialityID   *string `json:"old_hls_speciality_id" db:"OldHlsSpecialityId"`
	NewHLSEntityType     *string `json:"new_hls_entity_type" db:"NewHlsEntityType"`
	NewHLSEntityTypeID   *string `json:"new_hls_entity_id" db:"NewHlsEntityTypeId"`
	NewHLSEntityTypeCode *string `json:"new_hls_entity_type_code" db:"NewHlsEntityTypeCode"`
	NewHLSSpeciality     *string `json:"new_hls_speciality" db:"NewHlsSpeciality"`
	NewHLSSpecialityID   *string `json:"new_hls_speciality_id" db:"NewHlsSpecialityId"`

	HLSEstablishmentID                   *string `json:"hls_establishment_id,omitempty" db:"hls_establishment_id"`
	HLSEstablishmentLicenceID            *string `json:"hls_establishment_licence_id,omitempty" db:"hls_establishment_licence_id"`
	CityAr                               *string `json:"city_ar" db:"city_ar"`
	CityEn                               *string `json:"city_en" db:"city_en"`
	CityCode                             *string `json:"city_code" db:"city_code"`
	RegionAr                             *string `json:"region_ar" db:"region_ar"`
	RegionEn                             *string `json:"region_en" db:"region_en"`
	RegionCode                           *string `json:"region_code" db:"region_code"`
	LocationCode                         *string `json:"location_code" db:"location_code"`
	FullAddress                          *string `json:"full_address" db:"FullAddress"`
	Address                              *string `json:"address,omitempty"`
	AddressBuildingNumber                *string `json:"address_building_number,omitempty" db:"Address_BuildingNumber"`
	AddressStreetName                    *string `json:"address_street_name,omitempty" db:"Address_StreetName"`
	AddressBlockNumber                   *string `json:"address_block_number,omitempty" db:"Address_BlockNumber"`
	AddressDistrictName                  *string `json:"address_district_name,omitempty" db:"Address_DistrictName"`
	AddressPostalCode                    *string `json:"address_postal_code,omitempty" db:"Address_PostalCode"`
	AddressCity                          *string `json:"address_city,omitempty" db:"Address_City"`
	OwnerName                            *string `json:"owner_name,omitempty" db:"OwnerName"`
	SystemManagerNameAr                  *string `json:"system_manager_name_ar,omitempty" db:"SystemManagerNameAr"`
	SystemManagerIDNumber                *string `json:"system_manager_id_number,omitempty" db:"SystemManagerIDNumber"`
	SystemManagerMobileNumber            *string `json:"system_manager_mobile_number,omitempty" db:"SystemManagerMobileNumber"`
	SystemManagerEmail                   *string `json:"system_manager_email,omitempty" db:"SystemManagerEmail"`
	TechnicalSupervisorName              *string `json:"technical_supervisor_name,omitempty" db:"TechnicalSupervisorName"`
	TechnicalSupervisorCategory          *string `json:"technical_supervisor_category" db:"TechnicalSupervisorCategory"`
	TechnicalSupervisorSpeciality        *string `json:"technical_supervisor_speciality" db:"TechnicalSupervisorSpeciality"`
	TechnicalSupervisorLicenseExpiryDate *string `json:"technical_supervisor_license_expiry_date,omitempty" db:"TechnicalSupervisorLicenseExpiryDate"`
	AdministrativeDirectorName           *string `json:"administrative_director_name,omitempty" db:"AdministrativeDirectorName"`
}

type EstablishmentV2 struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
	Msg      *string `json:"msg,omitempty" db:"Msg"`

	FacilityAddress *string `json:"facility_address" db:"FacilityAddress"`
	PostalAddress   *string `json:"postal_address" db:"PostalAddress"`

	OrganizationName        *string `json:"organization_name" db:"OrganizationName"`
	OrganizationID          *string `json:"organization_id" db:"OrganizationId"`
	OrganizationSpeciality  *string `json:"organization_speciality" db:"OrganizationSpeciality"`
	OrganizationCredentials *string `json:"organization_credentials" db:"OrgCredentials"`
	OrganizationContact     *string `json:"organization_contact" db:"OrganizationContact"`

	FacilityUniqueIdentifier      *string `json:"facility_unique_identifier" db:"UniqueHealthcareFacilityIdentifier"`
	FacilitySectoridentifier      *string `json:"facility_sector_identifier" db:"HealthcareFacilitySectoridentifier"`
	FacilityLevelOfCareIdentifier *string `json:"facility_level_of_care_identifier" db:"HealthcareFacilityLevelOfCareIdentifier"`
	FacilityTypeOfCare            *string `json:"facility_type_of_care" db:"HealthcareFacilityTypeOfCare"`
	FacilityOperationStatus       *string `json:"facility_operation_status" db:"FacilityOperationStatus"`

	ElectronicServiceURI        *string `json:"electronic_service_url" db:"ElectronicServiceURI"`
	MedicalRecordsDeliveryEmail *string `json:"medical_records_delivery_email" db:"MedicalRecordsDeliveryEmailAddress"`
	LastUpdatedTime             *string `json:"last_updated_time" db:"DateTimeOfLastUpdate"`

	ProviderLanguageSupported *string `json:"provider_language_supported" db:"ProviderLanguageSupported"`
	ProviderRelationship      *string `json:"provider_relationship" db:"ProviderRelationship"`

	AvailableHospitalBedsAdmittedPatients             *string `json:"available_hospital_beds_admitted_Patients" db:"NumberOfAvailableHospitalBedsForAdmittedPatients"`
	AvailableOperatingRooms                           *string `json:"available_operating_rooms" db:"AvailableOperatingRooms"`
	AvailableEmergencyBeds                            *string `json:"available_emergency_beds" db:"AvailableEmergencyBeds"`
	AvailableIntensiveCareUnitSpecialAreasAverageBeds *string `json:"available_intensive_dare_unit_areas_average_beds" db:"AvailableIntensiveCareUnit_SpecialAreasAverageBeds"`

	TotalHospitalBeds *string `json:"total_hospital_beds" db:"TotalHospitalBeds"`
	TeachingStatus    *string `json:"teaching_status" db:"TeachingStatus"`
	The700Number      *string `json:"the_700_number" db:"The700Number"`
}
type Establishments struct {
	ErrorMsg *string `json:"error_msg,omitempty"`
	Msg      *string `json:"msg,omitempty" db:"Msg"`

	EntityType       *string `json:"entity_type,omitempty"`
	EntitySpeciality *string `json:"entity_speciality,omitempty"`

	ID             *string `json:"-" db:"id"`
	Code           *string `json:"code" db:"Code"`
	MOHID          *string `json:"-" db:"MOHID"`
	OrganizationID *string `json:"organization_id" db:"OrganizationId"`
	LegacyEntityID *string `json:"-" db:"legacy_entityId"`

	RowInsertedAt *string `json:"-" db:"RowInseartedAt"`
	RowUpdatedAt  *string `json:"-" db:"RowUpdatedAt"`
	RowDeletedAt  *string `json:"-" db:"RowDeletedAt"`
	IsMigrated    *string `json:"-" db:"isMigrated"`
	IsDeleted     *string `json:"-" db:"IsDeleted"`

	CRNumber            *string `json:"-" db:"CR_Number"`
	CREstablishmentName *string `json:"-" db:"CR_EstablishmentName"`
	SourceSystem        *string `json:"-" db:"SourceSystem"`

	NotificationEmail     *string `json:"-" db:"NotificationEmail"`
	LicenseNumber         *string `json:"-" db:"LicenseNumber"`
	IssueDate             *string `json:"-" db:"Issue_Date"`
	ExpiryDate            *string `json:"-" db:"Expiry_Date"`
	NameAr                *string `json:"name_ar,omitempty"`
	NameEn                *string `json:"name_en,omitempty"`
	SehaID                *string `json:"seha_id,omitempty"`
	BedsCount             *string `json:"beds_count,omitempty"`
	Longitude             *string `json:"longitude,omitempty"`
	Latitude              *string `json:"latitude,omitempty"`
	MapURL                *string `json:"-" db:"MapUrl"`
	SehaHealthDirectory   *string `json:"-" db:"SehaHealth_Directory"`
	LevelOfCare           *string `json:"level_of_care" db:"level_of_care"`
	TypeOfCare            *string `json:"type_of_care" db:"type_of_care"`
	TypeOfCareCode        *string `json:"type_of_care_code" db:"type_of_care_code"`
	HealthDirectoryAr     *string `json:"health_directory_ar" db:"Health_Directory_ar"`
	HealthDirectoryEn     *string `json:"health_directory_en" db:"Health_Directory_en"`
	HealthDirectorySehaID *string `json:"health_directory_seha_id" db:"Health_Directory_SehaID"`
	SectorAr              *string `json:"sector_ar" db:"sector_ar"`
	SectorEn              *string `json:"sector_en" db:"sector_en"`
	SectorCode            *string `json:"sector_code" db:"sector_code"`
	EntityTypeAr          *string `json:"-" db:"entity_type_ar"`
	EntityTypeEn          *string `json:"-" db:"entity_type_en"`
	EntityTypeCode        *string `json:"-" db:"entity_type_code"`
	Website               *string `json:"-" db:"website"`
	PhoneNumber           *string `json:"-" db:"phone_number"`
	Email                 *string `json:"-" db:"email"`

	OldHLSEntityType     *string `json:"-" db:"OldHlsEntityType"`
	OldHLSEntityTypeID   *string `json:"-" db:"OldHlsEntityTypeId"`
	OldHLSSpeciality     *string `json:"-" db:"OldHlsSpeciality"`
	OldHLSSpecialityID   *string `json:"-" db:"OldHlsSpecialityId"`
	NewHLSEntityType     *string `json:"-" db:"NewHlsEntityType"`
	NewHLSEntityTypeID   *string `json:"-" db:"NewHlsEntityTypeId"`
	NewHLSEntityTypeCode *string `json:"-" db:"NewHlsEntityTypeCode"`
	NewHLSSpeciality     *string `json:"-" db:"NewHlsSpeciality"`
	NewHLSSpecialityID   *string `json:"-" db:"NewHlsSpecialityId"`

	HLSEstablishmentID                   *string `json:"-" db:"hls_establishment_id"`
	HLSEstablishmentLicenceID            *string `json:"-" db:"hls_establishment_licence_id"`
	CityAr                               *string `json:"city_ar" db:"city_ar"`
	CityEn                               *string `json:"city_en" db:"city_en"`
	CityCode                             *string `json:"city_code" db:"city_code"`
	RegionAr                             *string `json:"region_ar" db:"region_ar"`
	RegionEn                             *string `json:"region_en" db:"region_en"`
	RegionCode                           *string `json:"region_code" db:"region_code"`
	LocationCode                         *string `json:"location_code" db:"location_code"`
	FullAddress                          *string `json:"-"`
	Address                              *string `json:"-"`
	AddressBuildingNumber                *string `json:"-" db:"Address_BuildingNumber"`
	AddressStreetName                    *string `json:"-" db:"Address_StreetName"`
	AddressBlockNumber                   *string `json:"-" db:"Address_BlockNumber"`
	AddressDistrictName                  *string `json:"-" db:"Address_DistrictName"`
	AddressPostalCode                    *string `json:"-" db:"Address_PostalCode"`
	AddressCity                          *string `json:"-" db:"Address_City"`
	OwnerName                            *string `json:"-" db:"OwnerName"`
	SystemManagerNameAr                  *string `json:"-" db:"SystemManagerNameAr"`
	SystemManagerIDNumber                *string `json:"-" db:"SystemManagerIDNumber"`
	SystemManagerMobileNumber            *string `json:"-" db:"SystemManagerMobileNumber"`
	SystemManagerEmail                   *string `json:"-" db:"SystemManagerEmail"`
	TechnicalSupervisorName              *string `json:"-" db:"TechnicalSupervisorName"`
	TechnicalSupervisorCategory          *string `json:"-" db:"TechnicalSupervisorCategory"`
	TechnicalSupervisorSpeciality        *string `json:"-" db:"TechnicalSupervisorSpeciality"`
	TechnicalSupervisorLicenseExpiryDate *string `json:"-" db:"TechnicalSupervisorLicenseExpiryDate"`
	AdministrativeDirectorName           *string `json:"-" db:"AdministrativeDirectorName"`
}

// Errors related to Establishment
var (
	ErrEmptyCode           = errors.New("code can't be empty in establishment")
	ErrEmptyOrganizationID = errors.New("organization_id can't be empty in establishment")
)

func (e *Establishment) Validate() error {
	if e.OrganizationID == nil {
		return ErrEmptyOrganizationID
	}
	if e.Code == nil {
		return ErrEmptyCode
	}
	return nil
}

type ISOCode struct {
	Code          *string `json:"error_msg,omitempty" db:"ISOAlpha3Code"`
	ErrorMsg      *string `json:"error_msg,omitempty"`
	CountryNameEn *string `db:"CountryNameEn"`
}

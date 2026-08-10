package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issueTracking/internal/db"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInsertDeathReportBadRequest(t *testing.T) {
	db.TruncateTables(t, testPool)
	payload := map[string]any{
		"reportedDate":            "22/07/2026",
		"incidentDate":            "21/07/2026",
		"incidentTime":            "14:30",
		"department":              "ICU",
		"location":                "Main Hospital",
		"category":                "Mortality",
		"subCategory":             "Unexpected death",
		"description":             "Patient experienced acute cardiac arrest following surgical procedure.",
		"actionTaken":             "CPR initiated immediately; resuscitation team responded. Pronounced dead at 15:05.",
		"openedDate":              "22/07/2026",
		"submittedTime":           "08:00",
		"handler":                 "Dr. John Doe",
		"manager":                 "Jane Smith",
		"specialty":               "Cardiology",
		"exactLocation":           "Bed 4, ICU Ward 2",
		"coding":                  "ICD-10-I46.9",
		"type":                    "Clinical Incident",
		"riskGrading":             "High",
		"result":                  "Fatal",
		"actualHarm":              "Severe / Death",
		"potentialHarm":           "Severe",
		"details":                 "Patient was undergoing routine post-op monitoring.",
		"patientInvolved":         true,
		"patientTold":             false,
		"familyTold":              true,
		"whatFamilyTold":          "Family was informed about cardiac complications and unsuccessful resuscitation efforts.",
		"incidentInvestigation":   "Internal review initiated by QA panel.",
		"reviewMeetingDate":       "25/07/2026",
		"qualityAssuranceLead":    "Dr. Alice Johnson",
		"docNotified":             true,
		"meetingDiscussionPoints": "Reviewed timeline of medication administration and monitoring telemetry logs.",
		"meetingActionPoints":     "Audit telemetry equipment calibration and update post-op cardiac monitoring protocol.",
		"levelOfInvestigation":    "Level 3",
	}
	jsonBody, err := json.Marshal(&payload)
	assert.NoError(t, err)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	r := gin.Default()
	r.POST("/api/v1/deathreport", a.deathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/deathreport", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestInsertDeathReportSuccess(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	payload := map[string]any{
		"ref":                     "DR-2026-001",
		"reportedDate":            "22/07/2026",
		"incidentDate":            "21/07/2026",
		"incidentTime":            "14:30",
		"department":              "ICU",
		"location":                "Main Hospital",
		"category":                "Mortality",
		"subCategory":             "Unexpected death",
		"description":             "Patient experienced acute cardiac arrest following surgical procedure.",
		"actionTaken":             "CPR initiated immediately; resuscitation team responded. Pronounced dead at 15:05.",
		"openedDate":              "22/07/2026",
		"submittedTime":           "08:00",
		"handler":                 "Dr. John Doe",
		"manager":                 "Jane Smith",
		"specialty":               "Cardiology",
		"exactLocation":           "Bed 4, ICU Ward 2",
		"coding":                  "ICD-10-I46.9",
		"type":                    "Clinical Incident",
		"riskGrading":             "High",
		"result":                  "Fatal",
		"actualHarm":              "Severe / Death",
		"potentialHarm":           "Severe",
		"details":                 "Patient was undergoing routine post-op monitoring.",
		"patientInvolved":         true,
		"patientTold":             false,
		"familyTold":              true,
		"whatFamilyTold":          "Family was informed about cardiac complications and unsuccessful resuscitation efforts.",
		"incidentInvestigation":   "Internal review initiated by QA panel.",
		"reviewMeetingDate":       "25/07/2026",
		"qualityAssuranceLead":    "Dr. Alice Johnson",
		"docNotified":             true,
		"meetingDiscussionPoints": "Reviewed timeline of medication administration and monitoring telemetry logs.",
		"meetingActionPoints":     "Audit telemetry equipment calibration and update post-op cardiac monitoring protocol.",
		"levelOfInvestigation":    "Level 3",
	}
	jsonBody, err := json.Marshal(&payload)
	assert.NoError(t, err)

	r := gin.Default()
	r.POST("/api/v1/deathreport", a.deathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/deathreport", bytes.NewBuffer(jsonBody))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "The death has been reported", response["message"])
}

func TestUpdateDeathReport(t *testing.T) {
	db.TruncateTables(t, testPool)
	payload := &db.DeathReport{
		Ref:                     "DR-2026-001",
		ReportedDate:            "2026-07-29",
		IncidentDate:            "2026-07-28",
		IncidentTime:            "14:30",
		Department:              "Cardiology",
		Location:                "Building A",
		Category:                "Clinical Incident",
		SubCategory:             "Patient Care",
		Description:             "Test death report description for automated testing.",
		ActionTaken:             "Immediate review initiated.",
		OpenedDate:              "2026-07-29",
		SubmittedTime:           "09:00",
		Handler:                 "Dr. John Doe",
		Manager:                 "Jane Smith",
		Specialty:               "Internal Medicine",
		ExactLocation:           "Ward 3, Bed 12",
		Coding:                  "COD-101",
		Type:                    "Clinical Incident",
		RiskGrading:             "High",
		Result:                  "Under Review",
		ActualHarm:              "Severe",
		PotentialHarm:           "Critical",
		Details:                 "Additional test details regarding the event.",
		PatientInvolved:         true,
		PatientTold:             true,
		FamilyTold:              true,
		WhatFamilyTold:          "Family was informed by the attending physician.",
		IncidentInvestigation:   "Investigation ongoing by QA lead.",
		ReviewMeetingDate:       "2026-08-01",
		QualityAssuranceLead:    "Dr. Alice Johnson",
		DoctorNotified:          true,
		MeetingDiscussionPoints: "Discussed protocol adherence and timeline of events.",
		MeetingActionPoints:     "Update ward monitoring checklists.",
		LevelOfInvestigation:    "Comprehensive",
	}

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	err := insertDeathReport(payload, a, t)
	assert.NoError(t, err, "error inserting seed data")
	reqPayload, _ := json.Marshal(&payload)

	r := gin.Default()
	r.PUT("/api/v1/deathreport", mockAuthMiddleware("superadmin"), a.updateDeathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/v1/deathreport", bytes.NewBuffer(reqPayload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetDeathReports(t *testing.T) {
	db.TruncateTables(t, testPool)
	payload := &db.DeathReport{
		Ref:                     "DR-2026-001",
		ReportedDate:            "2026-07-29",
		IncidentDate:            "2026-07-28",
		IncidentTime:            "14:30",
		Department:              "Cardiology",
		Location:                "Building A",
		Category:                "Clinical Incident",
		SubCategory:             "Patient Care",
		Description:             "Test death report description for automated testing.",
		ActionTaken:             "Immediate review initiated.",
		OpenedDate:              "2026-07-29",
		SubmittedTime:           "09:00",
		Handler:                 "Dr. John Doe",
		Manager:                 "Jane Smith",
		Specialty:               "Internal Medicine",
		ExactLocation:           "Ward 3, Bed 12",
		Coding:                  "COD-101",
		Type:                    "Clinical Incident",
		RiskGrading:             "High",
		Result:                  "Under Review",
		ActualHarm:              "Severe",
		PotentialHarm:           "Critical",
		Details:                 "Additional test details regarding the event.",
		PatientInvolved:         true,
		PatientTold:             true,
		FamilyTold:              true,
		WhatFamilyTold:          "Family was informed by the attending physician.",
		IncidentInvestigation:   "Investigation ongoing by QA lead.",
		ReviewMeetingDate:       "2026-08-01",
		QualityAssuranceLead:    "Dr. Alice Johnson",
		DoctorNotified:          true,
		MeetingDiscussionPoints: "Discussed protocol adherence and timeline of events.",
		MeetingActionPoints:     "Update ward monitoring checklists.",
		LevelOfInvestigation:    "Comprehensive",
	}
	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	err := insertDeathReport(payload, a, t)
	assert.NoError(t, err)

	r := gin.Default()
	r.GET("/api/v1/deathReports", mockAuthMiddleware("superadmin"), a.getDeathReports)
	w := httptest.NewRecorder()
	dummyPayload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})
	req, _ := http.NewRequest("GET", "/api/v1/deathReports", bytes.NewBuffer(dummyPayload))
	r.ServeHTTP(w, req)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	t.Log(response)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchDeathReportsNoQuery(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	err := insertDeathReportWithPayload(a, t)
	assert.NoError(t, err)
	dummyPayload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/searchDeathReports", mockAuthMiddleware("superadmin"), a.searchDeathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchDeathReports", bytes.NewBuffer(dummyPayload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchDeathReports(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	err := insertDeathReportWithPayload(a, t)
	assert.NoError(t, err)
	dummyPayload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/searchDeathReports", mockAuthMiddleware("superadmin"), a.searchDeathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchDeathReports?searchQuery=cardiology", bytes.NewBuffer(dummyPayload))
	r.ServeHTTP(w, req)

	var response map[string]any
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	t.Log(response)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestSearchDeathReportsByDate(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}
	_ = insertDeathReportWithPayload(a, t)
	dummyPayload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})
	r := gin.Default()
	r.GET("/api/v1/deathreports", mockAuthMiddleware("superadmin"), a.getDeathReports)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/deathreports?dateFrom=2026-07-29&dateTo=2026-08-05", bytes.NewBuffer(dummyPayload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGlobalSearchWithDate(t *testing.T) {
	db.TruncateTables(t, testPool)

	a := &application{
		origins: "*",
		models:  db.NewModels(testPool),
	}

	_ = insertDeathReportWithPayload(a, t)
	dummyPayload, _ := json.Marshal(&map[string]any{
		"test": "test",
	})

	r := gin.Default()
	r.GET("/api/v1/searchDeathReport", mockAuthMiddleware("superadmin"), a.searchDeathReport)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/searchDeathReport?searchQuery=dr-2026&dateFrom=2026-07-29&dateTo=2026-08-05", bytes.NewBuffer(dummyPayload))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

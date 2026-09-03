package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type Property struct {
	ID               string `json:"id"`
	Owner            string `json:"owner"`
	Name             string `json:"name,omitempty"`
	Address          string `json:"address"`
	Location         string `json:"location,omitempty"`
	Type             string `json:"type"`
	Value            string `json:"value"`
	DeclaredValue    int64  `json:"declaredValue"`
	DeclaredRent     string `json:"declaredRent,omitempty"`
	Status           string `json:"status"`
	OwnershipHistory int    `json:"ownershipHistory"`
	PaymentHistory   string `json:"paymentHistory"`
	SurveyStatus     string `json:"surveyStatus"`
	NextDue          string `json:"nextDue,omitempty"`
	CreatedAt        string `json:"createdAt"`
	UpdatedAt        string `json:"updatedAt"`
}

type AuditCase struct {
	ID           string  `json:"id"`
	Property     string  `json:"property"`
	Owner        string  `json:"owner"`
	Auditor      string  `json:"auditor"`
	Priority     string  `json:"priority"`
	Status       string  `json:"status"`
	Started      string  `json:"started"`
	Due          string  `json:"due"`
	ResultStatus *string `json:"resultStatus"`
	ResultNotes  string  `json:"resultNotes"`
	ResultSentAt *string `json:"resultSentAt"`
}

type FlaggedCase struct {
	ID           string  `json:"id"`
	FilingID     string  `json:"filingId"`
	Property     string  `json:"property"`
	Taxpayer     string  `json:"taxpayer"`
	Reason       string  `json:"reason"`
	ReceivedAt   string  `json:"receivedAt"`
	Status       string  `json:"status"`
	Priority     string  `json:"priority"`
	ResultStatus *string `json:"resultStatus"`
	ResultNotes  string  `json:"resultNotes"`
	ResultSentAt *string `json:"resultSentAt"`
}

type SuccessfulFiling struct {
	ID          string `json:"id"`
	FilingID    string `json:"filingId"`
	Property    string `json:"property"`
	Taxpayer    string `json:"taxpayer"`
	ValidatedAt string `json:"validatedAt"`
	Status      string `json:"status"`
}

type Notice struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Property string `json:"property"`
	Amount   string `json:"amount"`
	DueDate  string `json:"dueDate"`
	Status   string `json:"status"`
	Type     string `json:"type"`
	Response string `json:"response,omitempty"`
}

type Payment struct {
	ID         string `json:"id"`
	PropertyID string `json:"propertyId"`
	Amount     string `json:"amount"`
	Method     string `json:"method"`
	Status     string `json:"status"`
	CreatedAt  string `json:"createdAt"`
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}

type Database struct {
	Properties        []Property         `json:"properties"`
	AuditCases        []AuditCase        `json:"auditCases"`
	FlaggedCases      []FlaggedCase      `json:"flaggedCases"`
	SuccessfulFilings []SuccessfulFiling `json:"successfulFilings"`
	Notices           []Notice           `json:"notices"`
	Payments          []Payment          `json:"payments"`
}

type Store struct {
	mu   sync.RWMutex
	path string
	db   Database
}

func Open(path string) (*Store, error) {
	s := &Store{path: path}
	if err := s.load(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) Snapshot() Database {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneDB(s.db)
}

func (s *Store) Authenticate(email, role string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	role = strings.TrimSpace(strings.ToLower(role))
	if email == "" {
		return User{}, errors.New("email is required")
	}
	if role == "" {
		return User{}, errors.New("role is required")
	}
	name := strings.Split(email, "@")[0]
	name = strings.ReplaceAll(name, ".", " ")
	name = strings.Title(name)
	return User{Email: email, Name: name, Role: role}, nil
}

func (s *Store) ListProperties() []Property {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := append([]Property(nil), s.db.Properties...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

func (s *Store) CreateProperty(p Property) (Property, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := today()
	p.ID = nextID("PROP", len(s.db.Properties)+1)
	if p.Status == "" {
		p.Status = "Pending"
	}
	if p.PaymentHistory == "" {
		p.PaymentHistory = "Good"
	}
	if p.SurveyStatus == "" {
		p.SurveyStatus = "Pending"
	}
	if p.OwnershipHistory == 0 {
		p.OwnershipHistory = 1
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	s.db.Properties = append([]Property{p}, s.db.Properties...)
	return p, s.saveLocked()
}

func (s *Store) UpdateProperty(id string, patch Property) (Property, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Properties {
		if s.db.Properties[i].ID == id {
			patch.ID = id
			patch.CreatedAt = s.db.Properties[i].CreatedAt
			patch.UpdatedAt = today()
			s.db.Properties[i] = patch
			return patch, s.saveLocked()
		}
	}
	return Property{}, os.ErrNotExist
}

func (s *Store) DeleteProperty(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Properties {
		if s.db.Properties[i].ID == id {
			s.db.Properties = append(s.db.Properties[:i], s.db.Properties[i+1:]...)
			return s.saveLocked()
		}
	}
	return os.ErrNotExist
}

func (s *Store) ListAuditCases() []AuditCase {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]AuditCase(nil), s.db.AuditCases...)
}

func (s *Store) CreateAuditCase(c AuditCase) (AuditCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = nextID("AUD-2026", len(s.db.AuditCases)+1)
	if c.Auditor == "" {
		c.Auditor = "Current User"
	}
	if c.Status == "" {
		c.Status = "Pending"
	}
	if c.Priority == "" {
		c.Priority = "Medium"
	}
	c.Started = today()
	s.db.AuditCases = append([]AuditCase{c}, s.db.AuditCases...)
	return c, s.saveLocked()
}

func (s *Store) UpdateAuditCase(id string, patch AuditCase) (AuditCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.AuditCases {
		if s.db.AuditCases[i].ID == id {
			patch.ID = id
			if patch.Started == "" {
				patch.Started = s.db.AuditCases[i].Started
			}
			s.db.AuditCases[i] = patch
			return s.db.AuditCases[i], s.saveLocked()
		}
	}
	return AuditCase{}, os.ErrNotExist
}

func (s *Store) DeleteAuditCase(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.AuditCases {
		if s.db.AuditCases[i].ID == id {
			s.db.AuditCases = append(s.db.AuditCases[:i], s.db.AuditCases[i+1:]...)
			return s.saveLocked()
		}
	}
	return os.ErrNotExist
}

func (s *Store) SendAuditResult(id, status, notes string) (Database, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.AuditCases {
		if s.db.AuditCases[i].ID != id {
			continue
		}
		c := s.db.AuditCases[i]
		now := today()
		c.ResultStatus = &status
		c.ResultNotes = notes
		c.ResultSentAt = &now
		s.db.AuditCases = append(s.db.AuditCases[:i], s.db.AuditCases[i+1:]...)
		if status == "Compliant" {
			s.db.SuccessfulFilings = append([]SuccessfulFiling{{
				ID:          "SUC-" + strings.TrimPrefix(c.ID, "AUD-"),
				FilingID:    filingID(c.ID),
				Property:    c.Property,
				Taxpayer:    c.Owner,
				ValidatedAt: now,
				Status:      "Validated",
			}}, s.db.SuccessfulFilings...)
		} else {
			s.db.FlaggedCases = append([]FlaggedCase{{
				ID:           "FLG-" + strings.TrimPrefix(c.ID, "AUD-"),
				FilingID:     filingID(c.ID),
				Property:     c.Property,
				Taxpayer:     c.Owner,
				Reason:       "Declared rent below benchmark",
				ReceivedAt:   now,
				Status:       "Pending Review",
				Priority:     c.Priority,
				ResultStatus: c.ResultStatus,
				ResultNotes:  c.ResultNotes,
				ResultSentAt: c.ResultSentAt,
			}}, s.db.FlaggedCases...)
		}
		return cloneDB(s.db), s.saveLocked()
	}
	return Database{}, os.ErrNotExist
}

func (s *Store) ListFlaggedCases() []FlaggedCase {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]FlaggedCase(nil), s.db.FlaggedCases...)
}

func (s *Store) UpdateFlaggedCase(id string, patch FlaggedCase) (FlaggedCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.FlaggedCases {
		if s.db.FlaggedCases[i].ID == id {
			patch.ID = id
			s.db.FlaggedCases[i] = patch
			return patch, s.saveLocked()
		}
	}
	return FlaggedCase{}, os.ErrNotExist
}

func (s *Store) SendFlaggedResult(id, status, notes string) (FlaggedCase, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.FlaggedCases {
		if s.db.FlaggedCases[i].ID == id {
			now := today()
			s.db.FlaggedCases[i].ResultStatus = &status
			s.db.FlaggedCases[i].ResultNotes = notes
			s.db.FlaggedCases[i].ResultSentAt = &now
			if status == "Compliant" {
				s.db.FlaggedCases[i].Status = "Resolved"
			}
			return s.db.FlaggedCases[i], s.saveLocked()
		}
	}
	return FlaggedCase{}, os.ErrNotExist
}

func (s *Store) ListSuccessfulFilings() []SuccessfulFiling {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]SuccessfulFiling(nil), s.db.SuccessfulFilings...)
}

func (s *Store) UpdateSuccessfulFiling(id string, patch SuccessfulFiling) (SuccessfulFiling, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.SuccessfulFilings {
		if s.db.SuccessfulFilings[i].ID == id {
			patch.ID = id
			s.db.SuccessfulFilings[i] = patch
			return patch, s.saveLocked()
		}
	}
	return SuccessfulFiling{}, os.ErrNotExist
}

func (s *Store) ListNotices() []Notice {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Notice(nil), s.db.Notices...)
}

func (s *Store) RespondNotice(id, response string) (Notice, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Notices {
		if s.db.Notices[i].ID == id {
			s.db.Notices[i].Status = "Responded"
			s.db.Notices[i].Response = response
			return s.db.Notices[i], s.saveLocked()
		}
	}
	return Notice{}, os.ErrNotExist
}

func (s *Store) ListPayments() []Payment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return append([]Payment(nil), s.db.Payments...)
}

func (s *Store) CreatePayment(p Payment) (Payment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p.ID = nextID("PAY-2026", len(s.db.Payments)+1)
	p.Status = "Paid"
	p.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	s.db.Payments = append([]Payment{p}, s.db.Payments...)
	return p, s.saveLocked()
}

func (s *Store) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		s.db = seed()
		return s.saveLocked()
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(b, &s.db); err != nil {
		return fmt.Errorf("parse database: %w", err)
	}
	return nil
}

func (s *Store) saveLocked() error {
	b, err := json.MarshalIndent(s.db, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func cloneDB(db Database) Database {
	db.Properties = append([]Property(nil), db.Properties...)
	db.AuditCases = append([]AuditCase(nil), db.AuditCases...)
	db.FlaggedCases = append([]FlaggedCase(nil), db.FlaggedCases...)
	db.SuccessfulFilings = append([]SuccessfulFiling(nil), db.SuccessfulFilings...)
	db.Notices = append([]Notice(nil), db.Notices...)
	db.Payments = append([]Payment(nil), db.Payments...)
	return db
}

func nextID(prefix string, n int) string {
	return fmt.Sprintf("%s-%03d", prefix, n)
}

func today() string {
	return time.Now().UTC().Format("2006-01-02")
}

func filingID(id string) string {
	digits := strings.NewReplacer("AUD", "", "FLG", "", "SUC", "", "-", "").Replace(id)
	if len(digits) > 5 {
		digits = digits[len(digits)-5:]
	}
	return fmt.Sprintf("FCT-IRS-%05s", digits)
}

func strptr(s string) *string {
	return &s
}

func seed() Database {
	return Database{
		Properties: []Property{
			{ID: "PROP-001", Owner: "Emeka Okonkwo", Name: "Commercial Complex A", Address: "Plot 42, Victoria Island", Location: "Victoria Island, Lagos", Type: "Commercial", Value: "N250,000,000", DeclaredValue: 250000000, DeclaredRent: "N12,000,000", Status: "Verified", OwnershipHistory: 1, PaymentHistory: "Good", SurveyStatus: "Verified", NextDue: "Mar 31, 2026", CreatedAt: "2026-01-10", UpdatedAt: "2026-01-10"},
			{ID: "PROP-002", Owner: "Adaobi Nnamdi", Name: "Residential Estate B", Address: "Block 7, Lekki Phase 2", Location: "Lekki Phase 1, Lagos", Type: "Residential", Value: "N80,000,000", DeclaredValue: 80000000, DeclaredRent: "N8,500,000", Status: "Pending", OwnershipHistory: 1, PaymentHistory: "Good", SurveyStatus: "Pending", NextDue: "Mar 31, 2026", CreatedAt: "2026-01-12", UpdatedAt: "2026-01-12"},
			{ID: "PROP-003", Owner: "Chidi Okafor", Name: "Mixed Use Development C", Address: "15 Admiralty Way, Lekki", Location: "Ikoyi, Lagos", Type: "Mixed Use", Value: "N180,000,000", DeclaredValue: 180000000, DeclaredRent: "N18,000,000", Status: "Verified", OwnershipHistory: 3, PaymentHistory: "Good", SurveyStatus: "Verified", NextDue: "Jun 30, 2026", CreatedAt: "2026-01-13", UpdatedAt: "2026-01-13"},
			{ID: "PROP-004", Owner: "Folake Adeyemi", Name: "Office Tower D", Address: "Plot 8, Banana Island", Location: "Admiralty Way, Lekki", Type: "Residential", Value: "N500,000,000", DeclaredValue: 500000000, DeclaredRent: "N22,000,000", Status: "Flagged", OwnershipHistory: 4, PaymentHistory: "Late", SurveyStatus: "Flagged", NextDue: "Jun 30, 2026", CreatedAt: "2026-01-14", UpdatedAt: "2026-01-14"},
			{ID: "PROP-005", Owner: "Ibrahim Bello", Name: "Land Parcel E", Address: "Block 3, Ikoyi", Location: "Epe, Lagos", Type: "Commercial", Value: "N350,000,000", DeclaredValue: 350000000, DeclaredRent: "N1,200,000", Status: "Verified", OwnershipHistory: 2, PaymentHistory: "Good", SurveyStatus: "Verified", NextDue: "Pending", CreatedAt: "2026-01-15", UpdatedAt: "2026-01-15"},
		},
		AuditCases: []AuditCase{
			{ID: "AUD-2026-001", Property: "Plot 42, Victoria Island", Owner: "Emeka Okonkwo", Auditor: "John Smith", Priority: "High", Status: "In Progress", Started: "2026-01-10", Due: "2026-01-25"},
			{ID: "AUD-2026-002", Property: "Block 7, Lekki Phase 2", Owner: "Adaobi Nnamdi", Auditor: "Sarah Johnson", Priority: "Medium", Status: "Pending", Started: "2026-01-12", Due: "2026-01-30"},
		},
		FlaggedCases: []FlaggedCase{
			{ID: "FLG-2026-011", FilingID: "FCT-IRS-00211", Property: "Plot 18, Wuse II", Taxpayer: "Nwosu Holdings", Reason: "Declared rent below benchmark", ReceivedAt: "2026-01-18", Status: "Pending Review", Priority: "High"},
			{ID: "FLG-2026-012", FilingID: "FCT-IRS-00225", Property: "Block 4, Maitama", Taxpayer: "Ayo Martins", Reason: "Declared rent below benchmark", ReceivedAt: "2026-01-18", Status: "Pending Review", Priority: "Medium"},
			{ID: "FLG-2026-014", FilingID: "FCT-IRS-00237", Property: "Plot 9, Asokoro", Taxpayer: "Dara Okafor", Reason: "Declared rent below benchmark", ReceivedAt: "2026-01-20", Status: "In Progress", Priority: "Critical", ResultStatus: strptr("Non-Compliant"), ResultNotes: "Under-declaration confirmed. Notice issued.", ResultSentAt: strptr("2026-01-21")},
		},
		SuccessfulFilings: []SuccessfulFiling{
			{ID: "SUC-2026-101", FilingID: "FCT-IRS-00311", Property: "Plot 12, Jabi", Taxpayer: "Kola Ibrahim", ValidatedAt: "2026-01-18", Status: "Validated"},
			{ID: "SUC-2026-102", FilingID: "FCT-IRS-00318", Property: "Block 5, Garki", Taxpayer: "Laila Musa", ValidatedAt: "2026-01-18", Status: "Validated"},
		},
		Notices: []Notice{
			{ID: "NT-2026-001", Title: "Q1 2026 Assessment Notice", Property: "Commercial Complex A", Amount: "N2,500,000", DueDate: "Mar 31, 2026", Status: "Pending", Type: "Assessment"},
			{ID: "NT-2026-002", Title: "Property Valuation Update", Property: "Residential Estate B", Amount: "N1,800,000", DueDate: "Feb 28, 2026", Status: "Responded", Type: "Valuation"},
			{ID: "NT-2026-003", Title: "Annual Compliance Review", Property: "Office Tower D", Amount: "N4,200,000", DueDate: "Jun 30, 2026", Status: "Resolved", Type: "Compliance"},
		},
		Payments: []Payment{},
	}
}

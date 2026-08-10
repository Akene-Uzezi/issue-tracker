# Incident Tracker

RESTful API for managing workplace incidents and safety reports. Built with Go, Gin, and PostgreSQL.

## Prerequisites

- Docker and Docker Compose
- Go 1.26+

## Setup

```bash
git clone <repository-url>
cd incident-tracker-backend
cp .env.example .env
go mod download
```

Start with Docker:

```bash
docker compose up -d
```

API available at `http://localhost:3002`.

## Development

```bash
air
```

Or run directly:

```bash
go run ./cmd/
```

API available at `http://localhost:3001`.

## Tests

```bash
go test -v -tags=test ./...
```

Or run the helper script:

```bash
./scripts/runtests.sh
```

Format code:

```bash
go fmt ./...
```

Run linter:

```bash
go vet ./...
```

## Scripts

```bash
# Access PostgreSQL shell
./scripts/login.sh

# Reset database (drop and recreate tables)
./scripts/resetdb.sh

# Recreate tables without dropping data
./scripts/createtable.sh
```

## Default Credentials

A superadmin user is seeded by default:

- Email: `admin@example.com`
- Password: The password is stored as a bcrypt hash in `tables.sql`. Use the database or reset it via code to set a known password.

**Note:** New users registered via `/api/v1/auth/register` are assigned a default password of `redeemershealthvillage` if none is provided. This is separate from the pre-seeded superadmin.

**Ports:** The API listens on `localhost:3001` when run directly and `localhost:3002` when using Docker Compose.

## Linked Frontends

This API is shared by two Next.js frontends. Both authenticate against it and use the same JWT, so a user session is valid across both apps.

| Frontend | Default Dev URL | Dashboard Route |
|----------|-----------------|-----------------|
| `incident-tracker-frontend` | http://localhost:3000 | `/dashboard` |
| `death-report` | http://localhost:3001 | `/dashboard` |

Both frontend origins must be listed in `allowedOrigins` (see `.env.example`) so CORS permits their requests. The example configuration allows `http://localhost:3000` and `http://localhost:3001`.

## Quick Usage

Login to get a token:

```bash
TOKEN=$(curl -s -X POST http://localhost:3002/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"yourpassword"}' | jq -r '.token')
```

Register a new user (requires superadmin token):

```bash
curl -X POST http://localhost:3002/api/v1/auth/register \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","name":"New User","password":"password123","role":"admin","department":"IT"}'
```

Report an incident (no auth required):

```bash
curl -X POST http://localhost:3002/api/v1/incidents \
  -H "Content-Type: application/json" \
  -d '{
    "principalName": "John Doe",
    "principalGender": "Male",
    "principalDob": "1990-01-15",
    "principalType": "patient",
    "patientId": "P12345",
    "patientWardDept": "Ward A",
    "peopleInvolved": "Nurse Smith",
    "dateOfIncident": "2026-06-09",
    "timeOfIncident": "14:00",
    "locationOfIncident": "Ward A, Room 3",
    "incidentWardDept": "Ward A",
    "witnesses": "Dr. Brown",
    "witnessType": "Staff",
    "witnessWardDept": "Ward A",
    "witnessJobTitle": "Doctor",
    "witenssPhone": "555-0100",
    "isNearMiss": false,
    "causeGroup": "Fall",
    "causes": "Wet floor",
    "prescribingDoctor": "Dr. Brown",
    "treatmentReceived": "First Aid",
    "equipmentInvolved": "No",
    "equipmentSentForRepair": false,
    "equipmentWithdrawn": false,
    "equipmentRetained": false,
    "isMedicalDevice": "No",
    "reporterName": "Jane Reporter",
    "reporterDesignation": "Nurse",
    "signature": true,
    "reporterInfo": "jane@example.com",
    "date": "2026-06-09",
    "severityLevel": "minor"
  }'
```

List incidents (requires auth):

```bash
curl http://localhost:3002/api/v1/incidents -H "Authorization: Bearer $TOKEN"
```

Add comment to incident (requires manager or admin):

```bash
curl -X POST http://localhost:3002/api/v1/incidents/comments \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"incidentId": 1, "userId": 2, "comment": "Follow up needed"}'
```

Update incident status (requires auth; reporter/supervisor/manager roles forbidden):

```bash
curl -X PATCH http://localhost:3002/api/v1/incidents/1/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"resolved"}'
```

Get user info (requires superadmin role):

```bash
curl "http://localhost:3002/api/v1/user?email=test@example.com" -H "Authorization: Bearer $TOKEN"
```

Get comments for incident (requires admin or manager):

```bash
curl "http://localhost:3002/api/v1/incidents/comments?incidentId=1" -H "Authorization: Bearer $TOKEN"
```

Get incident management logs (requires admin or manager role):

```bash
curl "http://localhost:3002/api/v1/incidents/1/managementlogs" -H "Authorization: Bearer $TOKEN"
```

Report a death (no auth required):

```bash
curl -X POST http://localhost:3002/api/v1/deathreport \
  -H "Content-Type: application/json" \
  -d '{
    "ref": "DR-001",
    "reportedDate": "2026-06-09",
    "incidentDate": "2026-06-09",
    "incidentTime": "14:00",
    "department": "IT",
    "location": "Ward A",
    "category": "Category",
    "subCategory": "SubCategory",
    "description": "Description",
    "actionTaken": "Action taken",
    "openedDate": "2026-06-09",
    "submittedTime": "14:00",
    "handler": "Handler",
    "manager": "Manager",
    "specialty": "Specialty",
    "exactLocation": "Exact Location",
    "coding": "Coding",
    "type": "Type",
    "riskGrading": "High",
    "result": "Result",
    "actualHarm": "Actual Harm",
    "potentialHarm": "Potential Harm",
    "details": "Details",
    "patientInvolved": true,
    "patientTold": true,
    "familyTold": true,
    "whatFamilyTold": "What family told",
    "incidentInvestigation": "Investigation",
    "reviewMeetingDate": "2026-06-09",
    "qualityAssuranceLead": "QA Lead",
    "doctorNotified": true,
    "meetingDiscussionPoints": "Discussion points",
    "meetingActionPoints": "Action points",
    "levelOfInvestigation": "Level"
  }'
```

Update a death report (no auth required):

```bash
curl -X PUT http://localhost:3002/api/v1/deathreport \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "ref": "DR-001",
    "reportedDate": "2026-06-09",
    "incidentDate": "2026-06-09",
    "incidentTime": "14:00",
    "department": "IT",
    "location": "Ward A",
    "category": "Category",
    "subCategory": "SubCategory",
    "description": "Updated description",
    "actionTaken": "Updated action",
    "openedDate": "2026-06-09",
    "submittedTime": "14:00",
    "handler": "Handler",
    "manager": "Manager",
    "specialty": "Specialty",
    "exactLocation": "Exact Location",
    "coding": "Coding",
    "type": "Type",
    "riskGrading": "High",
    "result": "Result",
    "actualHarm": "Actual Harm",
    "potentialHarm": "Potential Harm",
    "details": "Details",
    "patientInvolved": true,
    "patientTold": true,
    "familyTold": true,
    "whatFamilyTold": "What family told",
    "incidentInvestigation": "Investigation",
    "reviewMeetingDate": "2026-06-09",
    "qualityAssuranceLead": "QA Lead",
    "doctorNotified": true,
    "meetingDiscussionPoints": "Discussion points",
    "meetingActionPoints": "Action points",
    "levelOfInvestigation": "Level"
  }'
```

Get all death reports (no auth required):

```bash
curl "http://localhost:3002/api/v1/deathreports?page=1&limit=10"
```

Search death reports (no auth required):

```bash
curl "http://localhost:3002/api/v1/searchDeathReport?searchQuery=DR-001"
```

## Role Permissions

| Role | Permissions |
|------|-------------|
| superadmin | All endpoints including user management (register, update, disable, enable, reset password, get user), report incidents, view all incidents, update any incident status, submit incident management reports, update incident management reports, add comments, view comments |
| admin | Report incidents, view all incidents, update any incident status, submit incident management reports, update incident management reports, add comments, view comments |
| supervisor | Report incidents, view own department incidents (matched via `incident_ward_dept`, `patient_ward_dept`, or `staff_place_of_work`) |
| manager | Report incidents, view all incidents, view incident management reports and logs, add comments, view comments, submit incident management reports, update incident management reports |
| reporter | Report incidents via public endpoint only, view own department incidents |

## Docker Commands

```bash
# Start all services (API at localhost:3002)
docker compose up -d

# Stop services
docker compose down

# Remove volumes (fresh database)
docker compose down -v

# View logs
docker compose logs -f
```

For full API documentation, request/response schemas, and role permissions, see [SYSTEM_DESIGN.md](SYSTEM_DESIGN.md).

For architecture, layering, and design decisions, see [ARCHITECTURE.md](ARCHITECTURE.md).

For database schema details, see `tables.sql`.

For user API routes, roles, inputs, and outputs, see [users.md](users.md).

# Cyber Threat Detection System

A real-time cybersecurity monitoring system that detects threats like credential stuffing, privilege escalation, account takeover, data exfiltration, and insider threats.

## Quick Start

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)

### Setup & Run

1. **Clone the repository**
```bash
git clone <repository-url>
cd cyberLogsThreat
```

2. **Start all services**
```bash
make build
make up
```

That's it! The application will be running at:
- **Frontend**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Elasticsearch**: http://localhost:9200

### Default Login Credentials

- **Username**: `admin`
- **Password**: `adminpassword`

## Available Commands
```bash
make up        # Start all services
make down      # Stop all services
make restart   # Restart all services
make logs      # View logs from all services
make build     # Rebuild and start all services
make clean     # Remove all containers and data
make install   # Install dependencies
make status    # Check service status
```

## How to Use

### 1. Login
- Go to http://localhost:3000
- Login with default credentials or register a new account

### 2. Upload Logs
- Click on "Upload Logs" card
- Select a CSV file with log data
- Click "Upload"

### 3. View Logs
- Click on "View Logs" to see all uploaded logs
- Use search filters to find specific logs

### 4. Analyze Threats
- Go to "Threats" page
- Click "Start Analysis" button to trigger threat detection
- The system will analyze all logs and detect security threats
- Analysis results appear automatically after completion

### 5. View Threats
- See all detected security threats on the Threats page
- Each threat shows:
  - **Severity Level**: Critical, High, Medium, Low
  - **Threat Type**: Credential Stuffing, Privilege Escalation, etc.
  - **Affected User**: User ID involved in the threat
  - **Timestamp**: When the threat was detected
- Use search filters to find specific threats by user, type, or severity

### 6. Dashboard
- Main page shows real-time analytics and charts
- View threat distribution, activity timeline, and severity breakdown
- Monitor system statistics: total logs, threats detected, unique users, threat rate

## Threat Analysis Workflow

1. **Upload Logs** → Import your security logs via CSV
2. **Start Analysis** → Click "Start Analysis" on Threats page
3. **Detection** → System analyzes patterns and detects threats
4. **Review Results** → View detected threats with severity ratings
5. **Take Action** → Investigate high-severity threats

**Required columns**: timestamp, userId, ipAddress, action  
**Optional columns**: fileName, databaseQuery

## Architecture
```
├── api-gateway-service      # API Gateway with JWT authentication
├── log-ingester-service     # Log management and storage
├── threat-analyzer-service  # Threat detection engine
├── cyber_threat_frontend    # Next.js frontend
└── docker-compose.yml       # Docker orchestration
```

## Detected Threats

The system automatically detects:

1. **Credential Stuffing** - Multiple failed logins followed by success
2. **Privilege Escalation** - Failed login + dangerous database queries
3. **Account Takeover** - Same user from different IPs within minutes
4. **Data Exfiltration** - Rapid access to multiple restricted files
5. **Insider Threat** - File access during off-hours (2 AM - 5 AM)

## Technology Stack

**Backend:**
- Go (Golang)
- Elasticsearch
- JWT Authentication

**Frontend:**
- Next.js 14
- React
- TypeScript
- Tailwind CSS
- Recharts

**Infrastructure:**
- Docker
- Docker Compose

## API Endpoints

### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/refresh` - Refresh JWT token
- `GET /api/auth/users` - Get all users (admin only)

### Logs
- `GET /api/logs` - Get all logs (paginated)
- `GET /api/logs/search` - Search logs with filters
- `POST /api/logs/upload` - Upload CSV file
- `GET /api/logs/{id}` - Get specific log

### Threats
- `GET /api/threats` - Get all threats (paginated)
- `GET /api/threats/search` - Search threats with filters
- `POST /api/threats/analyze` - Trigger threat analysis
- `GET /api/threats/{id}` - Get specific threat

## Troubleshooting

**Services won't start:**
```bash
make clean
make build
```

**Port already in use:**
```bash
# Stop conflicting services or change ports in docker-compose.yml
docker-compose down
```

**Can't login:**
```bash
# Reset to default admin user
make build
```

**View service logs:**
```bash
make logs
```
package constants

const (
	APIKeyHeader = "X-API-Key"
	
	ActionLoginSuccess = "loginSuccess"
	ActionLoginFailed  = "loginFailed"
	ActionFileAccess   = "fileAccess"
	ActionDBQuery      = "databaseQuery"
	ActionNetworkRequest  = "networkRequest"
	
	ThreatCredentialStuffing   = "Credential Stuffing"
	ThreatPrivilegeEscalation  = "Privilege Escalation"
	ThreatAccountTakeover      = "Account Takeover"
	ThreatDataExfiltration     = "Data Exfiltration"
	ThreatInsiderThreat        = "Insider Threat"
	
	SeverityLow      = "Low"
	SeverityMedium   = "Medium"
	SeverityHigh     = "High"
	SeverityCritical = "Critical"
	
	IndexLogs    = "logs"
	IndexThreats = "threats"

	UserName = "admin"
	UserPassword = "adminpassword"
)
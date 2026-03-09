connection "upguard" {
  plugin = "upguard"

  # API key from your UpGuard CyberRisk Account Settings → API keys
  # Required permissions depend on the tables you query:
  # - VendorRisk: for vendor-related tables
  # - BreachRisk: for organization domain/IP tables
  # - Platform: for available risks and organization info
  # See https://cyber-risk.upguard.com/api/docs for details
  # api_key = "YOUR_API_KEY"
}

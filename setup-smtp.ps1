# SMTP Setup Helper Script for Dossier (PowerShell)

Write-Host "🔧 Dossier SMTP Configuration Helper" -ForegroundColor Cyan
Write-Host "=====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Choose your SMTP configuration:"
Write-Host "1) Gmail"
Write-Host "2) Outlook/Hotmail"
Write-Host "3) Yahoo Mail"
Write-Host "4) MailHog (Local Testing - No real emails)"
Write-Host "5) Custom SMTP Server"
Write-Host "6) View current configuration"
Write-Host ""

$choice = Read-Host "Enter your choice (1-6)"

switch ($choice) {
    "1" {
        Write-Host ""
        Write-Host "📧 Gmail SMTP Configuration" -ForegroundColor Yellow
        Write-Host "Note: You need to generate an App Password at https://myaccount.google.com/apppasswords" -ForegroundColor Yellow
        Write-Host ""
        $gmail = Read-Host "Enter your Gmail address"
        $password = Read-Host "Enter your App Password (16 characters)" -AsSecureString
        $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))

        $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable

# Server Configuration
PORT=8080

# OpenAI Configuration (Optional - works with mock summaries without it)
# Get your API key from https://platform.openai.com/api-keys
# OPENAI_API_KEY=

# Local LLM Configuration (Free alternative to OpenAI)
USE_LOCAL_LLM=true
OLLAMA_URL=http://ollama:11434

# SMTP Configuration for Email Delivery
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=$gmail
SMTP_PASSWORD=$passwordPlain
SMTP_FROM_EMAIL=$gmail
SMTP_FROM_NAME=Dossier
"@
        $envContent | Out-File -FilePath ".env" -Encoding utf8
        Write-Host "✅ Gmail configuration saved!" -ForegroundColor Green
    }

    "2" {
        Write-Host ""
        Write-Host "📧 Outlook/Hotmail SMTP Configuration" -ForegroundColor Yellow
        Write-Host ""
        $outlook = Read-Host "Enter your Outlook/Hotmail address"
        $password = Read-Host "Enter your password" -AsSecureString
        $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))

        $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable

# Server Configuration
PORT=8080

# OpenAI Configuration (Optional - works with mock summaries without it)
# Get your API key from https://platform.openai.com/api-keys
# OPENAI_API_KEY=

# Local LLM Configuration (Free alternative to OpenAI)
USE_LOCAL_LLM=true
OLLAMA_URL=http://ollama:11434

# SMTP Configuration for Email Delivery
SMTP_HOST=smtp-mail.outlook.com
SMTP_PORT=587
SMTP_USERNAME=$outlook
SMTP_PASSWORD=$passwordPlain
SMTP_FROM_EMAIL=$outlook
SMTP_FROM_NAME=Dossier
"@
        $envContent | Out-File -FilePath ".env" -Encoding utf8
        Write-Host "✅ Outlook configuration saved!" -ForegroundColor Green
    }

    "3" {
        Write-Host ""
        Write-Host "📧 Yahoo Mail SMTP Configuration" -ForegroundColor Yellow
        Write-Host "Note: You may need to enable 'Less secure app access' in Yahoo settings" -ForegroundColor Yellow
        Write-Host ""
        $yahoo = Read-Host "Enter your Yahoo email address"
        $password = Read-Host "Enter your password" -AsSecureString
        $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))

        $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable

# Server Configuration
PORT=8080

# OpenAI Configuration (Optional - works with mock summaries without it)
# Get your API key from https://platform.openai.com/api-keys
# OPENAI_API_KEY=

# Local LLM Configuration (Free alternative to OpenAI)
USE_LOCAL_LLM=true
OLLAMA_URL=http://ollama:11434

# SMTP Configuration for Email Delivery
SMTP_HOST=smtp.mail.yahoo.com
SMTP_PORT=587
SMTP_USERNAME=$yahoo
SMTP_PASSWORD=$passwordPlain
SMTP_FROM_EMAIL=$yahoo
SMTP_FROM_NAME=Dossier
"@
        $envContent | Out-File -FilePath ".env" -Encoding utf8
        Write-Host "✅ Yahoo configuration saved!" -ForegroundColor Green
    }

    "4" {
        Write-Host ""
        Write-Host "🧪 MailHog Local Testing Configuration" -ForegroundColor Yellow
        Write-Host "This will catch emails locally without sending real emails." -ForegroundColor Yellow
        Write-Host "Remember to uncomment the MailHog service in docker-compose.yml" -ForegroundColor Yellow

        $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable

# Server Configuration
PORT=8080

# OpenAI Configuration (Optional - works with mock summaries without it)
# Get your API key from https://platform.openai.com/api-keys
# OPENAI_API_KEY=

# Local LLM Configuration (Free alternative to OpenAI)
USE_LOCAL_LLM=true
OLLAMA_URL=http://ollama:11434

# SMTP Configuration for Email Delivery
SMTP_HOST=mailhog
SMTP_PORT=1025
SMTP_USERNAME=
SMTP_PASSWORD=
SMTP_FROM_EMAIL=dossier@localhost
SMTP_FROM_NAME=Dossier
"@
        $envContent | Out-File -FilePath ".env" -Encoding utf8
        Write-Host "✅ MailHog configuration saved!" -ForegroundColor Green
        Write-Host "📝 Next steps:" -ForegroundColor Cyan
        Write-Host "   1. Uncomment the MailHog service in docker-compose.yml"
        Write-Host "   2. Run: docker-compose up mailhog -d"
        Write-Host "   3. View emails at http://localhost:8025"
    }

    "5" {
        Write-Host ""
        Write-Host "⚙️ Custom SMTP Configuration" -ForegroundColor Yellow
        Write-Host ""
        $host = Read-Host "Enter SMTP host"
        $port = Read-Host "Enter SMTP port"
        $username = Read-Host "Enter SMTP username"
        $password = Read-Host "Enter SMTP password" -AsSecureString
        $passwordPlain = [Runtime.InteropServices.Marshal]::PtrToStringAuto([Runtime.InteropServices.Marshal]::SecureStringToBSTR($password))
        $fromEmail = Read-Host "Enter from email"
        $fromName = Read-Host "Enter from name"

        $envContent = @"
# Database Configuration
DATABASE_URL=postgres://postgres:postgres@localhost:5432/dossier?sslmode=disable

# Server Configuration
PORT=8080

# OpenAI Configuration (Optional - works with mock summaries without it)
# Get your API key from https://platform.openai.com/api-keys
# OPENAI_API_KEY=

# Local LLM Configuration (Free alternative to OpenAI)
USE_LOCAL_LLM=true
OLLAMA_URL=http://ollama:11434

# SMTP Configuration for Email Delivery
SMTP_HOST=$host
SMTP_PORT=$port
SMTP_USERNAME=$username
SMTP_PASSWORD=$passwordPlain
SMTP_FROM_EMAIL=$fromEmail
SMTP_FROM_NAME=$fromName
"@
        $envContent | Out-File -FilePath ".env" -Encoding utf8
        Write-Host "✅ Custom SMTP configuration saved!" -ForegroundColor Green
    }

    "6" {
        Write-Host ""
        Write-Host "📋 Current SMTP Configuration:" -ForegroundColor Cyan
        Write-Host "=============================="
        if (Test-Path ".env") {
            Get-Content ".env" | Where-Object { $_ -match "SMTP_" }
        } else {
            Write-Host "No .env file found!" -ForegroundColor Red
        }
    }

    default {
        Write-Host "Invalid choice. Please run the script again." -ForegroundColor Red
    }
}

if ($choice -ne "6") {
    Write-Host ""
    Write-Host "🔄 Restart the backend to apply changes:" -ForegroundColor Cyan
    Write-Host "   docker-compose restart backend"
    Write-Host ""
    Write-Host "🧪 Test your configuration:" -ForegroundColor Cyan
    Write-Host "   Go to http://localhost:5173 and click 'Test Email'"
}
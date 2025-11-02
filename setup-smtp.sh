#!/bin/bash
# SMTP Setup Helper Script for Dossier

echo "🔧 Dossier SMTP Configuration Helper"
echo "====================================="
echo ""
echo "Choose your SMTP configuration:"
echo "1) Gmail"
echo "2) Outlook/Hotmail"  
echo "3) Yahoo Mail"
echo "4) MailHog (Local Testing - No real emails)"
echo "5) Custom SMTP Server"
echo "6) View current configuration"
echo ""
read -p "Enter your choice (1-6): " choice

case $choice in
    1)
        echo ""
        echo "📧 Gmail SMTP Configuration"
        echo "Note: You need to generate an App Password at https://myaccount.google.com/apppasswords"
        echo ""
        read -p "Enter your Gmail address: " gmail
        read -s -p "Enter your App Password (16 characters): " password
        echo ""

        cat > .env << EOF
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
SMTP_PASSWORD=$password
SMTP_FROM_EMAIL=$gmail
SMTP_FROM_NAME=Dossier
EOF
        echo "✅ Gmail configuration saved!"
        ;;

    2)
        echo ""
        echo "📧 Outlook/Hotmail SMTP Configuration"
        echo ""
        read -p "Enter your Outlook/Hotmail address: " outlook
        read -s -p "Enter your password: " password
        echo ""

        cat > .env << EOF
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
SMTP_PASSWORD=$password
SMTP_FROM_EMAIL=$outlook
SMTP_FROM_NAME=Dossier
EOF
        echo "✅ Outlook configuration saved!"
        ;;

    3)
        echo ""
        echo "📧 Yahoo Mail SMTP Configuration"
        echo "Note: You may need to enable 'Less secure app access' in Yahoo settings"
        echo ""
        read -p "Enter your Yahoo email address: " yahoo
        read -s -p "Enter your password: " password
        echo ""

        cat > .env << EOF
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
SMTP_PASSWORD=$password
SMTP_FROM_EMAIL=$yahoo
SMTP_FROM_NAME=Dossier
EOF
        echo "✅ Yahoo configuration saved!"
        ;;

    4)
        echo ""
        echo "🧪 MailHog Local Testing Configuration"
        echo "This will catch emails locally without sending real emails."
        echo "Remember to uncomment the MailHog service in docker-compose.yml"

        cat > .env << EOF
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
EOF
        echo "✅ MailHog configuration saved!"
        echo "📝 Next steps:"
        echo "   1. Uncomment the MailHog service in docker-compose.yml"
        echo "   2. Run: docker-compose up mailhog -d"
        echo "   3. View emails at http://localhost:8025"
        ;;

    5)
        echo ""
        echo "⚙️ Custom SMTP Configuration"
        echo ""
        read -p "Enter SMTP host: " host
        read -p "Enter SMTP port: " port
        read -p "Enter SMTP username: " username
        read -s -p "Enter SMTP password: " password
        echo ""
        read -p "Enter from email: " from_email
        read -p "Enter from name: " from_name

        cat > .env << EOF
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
SMTP_PASSWORD=$password
SMTP_FROM_EMAIL=$from_email
SMTP_FROM_NAME=$from_name
EOF
        echo "✅ Custom SMTP configuration saved!"
        ;;

    6)
        echo ""
        echo "📋 Current SMTP Configuration:"
        echo "=============================="
        if [ -f .env ]; then
            grep "SMTP_" .env
        else
            echo "No .env file found!"
        fi
        ;;

    *)
        echo "Invalid choice. Please run the script again."
        ;;
esac

if [ "$choice" != "6" ]; then
    echo ""
    echo "🔄 Restart the backend to apply changes:"
    echo "   docker-compose restart backend"
    echo ""
    echo "🧪 Test your configuration:"
    echo "   Go to http://localhost:5173 and click 'Test Email'"
fi
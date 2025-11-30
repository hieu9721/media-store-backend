#!/bin/bash
echo "Generating JWT Secret..."
SECRET=$(openssl rand -hex 32)
echo ""
echo "Add this to your .env file:"
echo "JWT_SECRET=$SECRET"
echo ""
echo "Keep this secret safe and never commit it to version control!"

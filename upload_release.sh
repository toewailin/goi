#!/bin/bash

# Variables
GITHUB_USER="your-username"               # GitHub username
GITHUB_REPO="your-repo-name"              # Repository name
GITHUB_TOKEN="your-github-token"          # GitHub Personal Access Token (PAT)
VERSION="v1.0.2"                          # Release version (tag)
BINARY_PATH="build/goi-linux-amd64"       # Path to your compiled binary
RELEASE_NAME="GoI CLI Release v1.0.2"      # Name of the release
RELEASE_NOTES="First release of GoI CLI tool"  # Release notes

# 1. Create a new release on GitHub
CREATE_RELEASE=$(curl -s -X POST https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/releases \
  -H "Authorization: token $GITHUB_TOKEN" \
  -d @- << EOF
{
  "tag_name": "$VERSION",
  "target_commitish": "main",  # The branch you're releasing from, e.g., "main"
  "name": "$RELEASE_NAME",
  "body": "$RELEASE_NOTES",
  "draft": false,
  "prerelease": false
}
EOF
)

# Get the release ID from the response
UPLOAD_URL=$(echo $CREATE_RELEASE | jq -r '.upload_url' | sed -e "s/{?name,label}//")

# 2. Upload the compiled binary to the release as an asset
UPLOAD_RESPONSE=$(curl -s -X POST "$UPLOAD_URL?name=$(basename $BINARY_PATH)" \
  -H "Authorization: token $GITHUB_TOKEN" \
  -H "Content-Type: application/octet-stream" \
  --data-binary @"$BINARY_PATH")

# Check for successful upload
UPLOAD_SUCCESS=$(echo $UPLOAD_RESPONSE | jq -r '.browser_download_url')

if [[ "$UPLOAD_SUCCESS" != "null" ]]; then
  echo "Successfully uploaded binary to GitHub release: $UPLOAD_SUCCESS"
else
  echo "Failed to upload binary. Response: $UPLOAD_RESPONSE"
  exit 1
fi

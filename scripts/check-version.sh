#!/bin/bash

# Release automation script demonstrating VERSION file usage
# This script shows how the VERSION file serves as single source of truth

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Read current version from VERSION file
if [ ! -f "VERSION" ]; then
    echo -e "${RED}ERROR: VERSION file not found${NC}"
    exit 1
fi

CURRENT_VERSION=$(cat VERSION)
echo -e "${GREEN}Current version: v${CURRENT_VERSION}${NC}"

# Validate version format (semantic versioning)
if ! [[ $CURRENT_VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}ERROR: Invalid version format in VERSION file. Expected: MAJOR.MINOR.PATCH${NC}"
    exit 1
fi

# Function to check if version exists in CHANGELOG
check_changelog() {
    if ! grep -q "## \[v${CURRENT_VERSION}\]" CHANGELOG.md; then
        echo -e "${YELLOW}WARNING: Version v${CURRENT_VERSION} not found in CHANGELOG.md${NC}"
        return 1
    fi
    return 0
}

# Function to check if version exists in README
check_readme() {
    if ! grep -q "v${CURRENT_VERSION}" README.md; then
        echo -e "${YELLOW}WARNING: Version v${CURRENT_VERSION} not found in README.md${NC}"
        return 1
    fi
    return 0
}

# Function to check if git tag exists
check_git_tag() {
    if git tag | grep -q "^v${CURRENT_VERSION}$"; then
        echo -e "${YELLOW}WARNING: Git tag v${CURRENT_VERSION} already exists${NC}"
        return 1
    fi
    return 0
}

echo "Checking version consistency..."

# Check if version is documented
check_changelog
CHANGELOG_OK=$?

check_readme  
README_OK=$?

check_git_tag
TAG_OK=$?

if [ $CHANGELOG_OK -eq 0 ] && [ $README_OK -eq 0 ] && [ $TAG_OK -eq 0 ]; then
    echo -e "${GREEN}✓ Version v${CURRENT_VERSION} is ready for release${NC}"
    
    # Ask for confirmation
    read -p "Do you want to create git tag v${CURRENT_VERSION}? (y/N): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Create git tag
        git tag "v${CURRENT_VERSION}" -m "Release v${CURRENT_VERSION}"
        echo -e "${GREEN}✓ Created git tag v${CURRENT_VERSION}${NC}"
        
        echo "To push the tag to remote:"
        echo "  git push origin v${CURRENT_VERSION}"
    fi
else
    echo -e "${RED}✗ Version inconsistencies found. Please update documentation before release.${NC}"
    
    if [ $CHANGELOG_OK -ne 0 ]; then
        echo "  - Add v${CURRENT_VERSION} entry to CHANGELOG.md"
    fi
    
    if [ $README_OK -ne 0 ]; then
        echo "  - Update version references in README.md"
    fi
    
    if [ $TAG_OK -ne 0 ]; then
        echo "  - Git tag v${CURRENT_VERSION} already exists"
    fi
fi

# Show example of other automation uses
echo
echo "Other automation uses of VERSION file:"
echo "  # Read version in scripts:"
echo "  VERSION=\$(cat VERSION)"
echo
echo "  # Update go.mod:"
echo "  go mod edit -module=github.com/mustanish/common-utils/v\$(cat VERSION | cut -d. -f1)"
echo
echo "  # Generate release notes:"
echo "  gh release create v\$(cat VERSION) --title \"Release v\$(cat VERSION)\" --notes-file CHANGELOG.md"
echo
echo "  # Docker build:"
echo "  docker build -t common-utils:v\$(cat VERSION) ."
#!/bin/bash
# Component Audit: Verifies Flowbite Default Component Usage
# 
# ASSUMPTION: All UI components come from flowbite-react
# This script verifies that assumption is true.

echo "🔍 Component Audit: Verifying Flowbite Default Component Usage"
echo ""
echo "ASSUMPTION: All UI components come from flowbite-react"
echo ""

VIOLATIONS=0
FLOWBITE_COUNT=0

# Check for Flowbite imports
echo "📊 Checking component files..."
echo ""

# Find all component files
find src/components -name "*.tsx" -o -name "*.ts" | while read file; do
    # Check if file imports from flowbite-react
    if grep -q "from ['\"]flowbite-react['\"]" "$file"; then
        FLOWBITE_IMPORTS=$(grep -o "import {[^}]*} from ['\"]flowbite-react['\"]" "$file" | sed 's/import {\([^}]*\)}.*/\1/' | tr ',' '\n' | sed 's/^[[:space:]]*//' | sed 's/[[:space:]]*$//')
        if [ ! -z "$FLOWBITE_IMPORTS" ]; then
            echo "✅ $file"
            echo "$FLOWBITE_IMPORTS" | while read component; do
                echo "   Flowbite: $component"
                FLOWBITE_COUNT=$((FLOWBITE_COUNT + 1))
            done
        fi
    fi
    
    # Check for violations (other UI libraries)
    if grep -q "from ['\"]react-bootstrap['\"]" "$file" || grep -q 'from ["'\'']react-bootstrap["'\'']' "$file"; then
        echo "❌ $file: Uses react-bootstrap instead of flowbite-react"
        VIOLATIONS=$((VIOLATIONS + 1))
    fi
    
    if grep -q "from ['\"]antd['\"]" "$file" || grep -q 'from ["'\'']antd["'\'']' "$file"; then
        echo "❌ $file: Uses antd instead of flowbite-react"
        VIOLATIONS=$((VIOLATIONS + 1))
    fi
    
    if grep -q "from ['\"]@mui" "$file" || grep -q 'from ["'\'']@mui' "$file"; then
        echo "❌ $file: Uses Material-UI instead of flowbite-react"
        VIOLATIONS=$((VIOLATIONS + 1))
    fi
done

echo ""
echo "📈 Summary:"
echo "   Flowbite components found: $FLOWBITE_COUNT"
echo "   Violations: $VIOLATIONS"
echo ""

if [ $VIOLATIONS -gt 0 ]; then
    echo "❌ ASSERTION FAILED: Component usage violates Flowbite default assumption"
    echo "   Fix violations before proceeding."
    exit 1
else
    echo "✅ ASSERTION PASSED: All components use Flowbite defaults"
    exit 0
fi

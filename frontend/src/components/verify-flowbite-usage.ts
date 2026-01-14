/**
 * Component Audit: Verifies Flowbite Default Component Usage
 * 
 * ASSUMPTION: All UI components come from flowbite-react
 * This script verifies that assumption is true.
 * 
 * Run with: npx tsx src/components/verify-flowbite-usage.ts
 */

import { readFileSync, readdirSync, statSync } from 'fs';
import { join } from 'path';

// Flowbite React components that should be used
const FLOWBITE_COMPONENTS = [
  'Button',
  'Card',
  'Checkbox',
  'Label',
  'Datepicker',
  'Table',
  'Badge',
  'Spinner',
  'Alert',
  'TextInput',
  'Select',
  'Modal',
  'Dropdown',
  'Tooltip',
  'Progress',
] as const;

// Components directory
const COMPONENTS_DIR = join(process.cwd(), 'src/components');

// Custom components that are allowed (app-specific logic)
const ALLOWED_CUSTOM_COMPONENTS = [
  'ProgressIndicator', // Custom progress bar (could use Flowbite Progress, but custom is simpler)
] as const;

interface ComponentUsage {
  file: string;
  flowbiteImports: string[];
  customComponents: string[];
  violations: string[];
}

function findComponentFiles(dir: string, fileList: string[] = []): string[] {
  const files = readdirSync(dir);
  
  for (const file of files) {
    const filePath = join(dir, file);
    const stat = statSync(filePath);
    
    if (stat.isDirectory()) {
      findComponentFiles(filePath, fileList);
    } else if (file.endsWith('.tsx') || file.endsWith('.ts')) {
      fileList.push(filePath);
    }
  }
  
  return fileList;
}

function analyzeFile(filePath: string): ComponentUsage {
  const content = readFileSync(filePath, 'utf-8');
  const relativePath = filePath.replace(process.cwd() + '/', '');
  
  // Find Flowbite imports
  const flowbiteImportRegex = /from\s+['"]flowbite-react['"]/;
  const flowbiteMatch = content.match(/import\s+{([^}]+)}\s+from\s+['"]flowbite-react['"]/);
  
  const flowbiteImports: string[] = [];
  if (flowbiteMatch) {
    flowbiteImports.push(...flowbiteMatch[1].split(',').map(s => s.trim()));
  }
  
  // Find custom component definitions (export function ComponentName)
  const customComponentRegex = /export\s+(?:function|const)\s+([A-Z][a-zA-Z0-9]*)/g;
  const customComponents: string[] = [];
  let match;
  while ((match = customComponentRegex.exec(content)) !== null) {
    const componentName = match[1];
    // Skip if it's a Flowbite component or allowed custom
    if (!FLOWBITE_COMPONENTS.includes(componentName as any) && 
        !ALLOWED_CUSTOM_COMPONENTS.includes(componentName as any)) {
      customComponents.push(componentName);
    }
  }
  
  // Check for violations
  const violations: string[] = [];
  
  // Check if using non-Flowbite UI libraries
  if (content.includes("from 'react-bootstrap'") || 
      content.includes('from "react-bootstrap"')) {
    violations.push('Uses react-bootstrap instead of flowbite-react');
  }
  
  if (content.includes("from 'antd'") || content.includes('from "antd"')) {
    violations.push('Uses antd instead of flowbite-react');
  }
  
  if (content.includes("from 'mui'") || content.includes('from "@mui"')) {
    violations.push('Uses Material-UI instead of flowbite-react');
  }
  
  // Check for custom replacements of Flowbite components
  if (customComponents.some(c => FLOWBITE_COMPONENTS.includes(c as any))) {
    violations.push(`Custom component replaces Flowbite component: ${customComponents.filter(c => FLOWBITE_COMPONENTS.includes(c as any)).join(', ')}`);
  }
  
  return {
    file: relativePath,
    flowbiteImports,
    customComponents,
    violations,
  };
}

function main() {
  console.log('🔍 Component Audit: Verifying Flowbite Default Component Usage\n');
  console.log('ASSUMPTION: All UI components come from flowbite-react\n');
  
  const componentFiles = findComponentFiles(COMPONENTS_DIR);
  const results: ComponentUsage[] = [];
  
  for (const file of componentFiles) {
    const analysis = analyzeFile(file);
    if (analysis.flowbiteImports.length > 0 || analysis.customComponents.length > 0 || analysis.violations.length > 0) {
      results.push(analysis);
    }
  }
  
  // Report results
  let totalViolations = 0;
  let totalFlowbiteComponents = 0;
  let totalCustomComponents = 0;
  
  console.log('📊 Results:\n');
  
  for (const result of results) {
    if (result.flowbiteImports.length > 0) {
      console.log(`✅ ${result.file}`);
      console.log(`   Flowbite: ${result.flowbiteImports.join(', ')}`);
      totalFlowbiteComponents += result.flowbiteImports.length;
    }
    
    if (result.customComponents.length > 0) {
      console.log(`   Custom: ${result.customComponents.join(', ')}`);
      totalCustomComponents += result.customComponents.length;
    }
    
    if (result.violations.length > 0) {
      console.log(`   ❌ VIOLATIONS:`);
      result.violations.forEach(v => console.log(`      - ${v}`));
      totalViolations += result.violations.length;
    }
    
    console.log('');
  }
  
  // Summary
  console.log('📈 Summary:');
  console.log(`   Flowbite components used: ${totalFlowbiteComponents}`);
  console.log(`   Custom components: ${totalCustomComponents}`);
  console.log(`   Violations: ${totalViolations}\n`);
  
  // Assertion
  if (totalViolations > 0) {
    console.log('❌ ASSERTION FAILED: Component usage violates Flowbite default assumption');
    console.log('   Fix violations before proceeding.\n');
    process.exit(1);
  } else {
    console.log('✅ ASSERTION PASSED: All components use Flowbite defaults\n');
    process.exit(0);
  }
}

if (require.main === module) {
  main();
}

export { analyzeFile, findComponentFiles, FLOWBITE_COMPONENTS, ALLOWED_CUSTOM_COMPONENTS };

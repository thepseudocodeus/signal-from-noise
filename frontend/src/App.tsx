import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom';
import { QueryState, PRODUCTION_REQUESTS, YearRange, DataCategory, DateRange } from './types/query';
import { ProductionRequestSelectionPage } from './components/pages/ProductionRequestSelectionPage';
import { SplashStep } from './components/steps/SplashStep';
import { ProductionRequestStep } from './components/steps/ProductionRequestStep';
import { CategoryStep } from './components/steps/CategoryStep';
import { DateRangeStep } from './components/steps/DateRangeStep';
import { FileSelectionStep } from './components/steps/FileSelectionStep';

type Step = 'splash' | 'production-request' | 'category' | 'date-range' | 'file-selection';

function AppContent() {
  const navigate = useNavigate();
  const [currentStep, setCurrentStep] = useState<Step>('splash');
  const [selectedProductionRequestNumber, setSelectedProductionRequestNumber] = useState<number | null>(null);
  const [query, setQuery] = useState<QueryState>({
    productionRequest: null,
    categories: [],
    yearRange: {
      startYear: null,
      endYear: null
    },
    dateRange: null
  });

  const handleProductionRequestStart = (requestNumber: number) => {
    setSelectedProductionRequestNumber(requestNumber);
    // Navigate to main application flow
    navigate('/app');
    setCurrentStep('splash');
  };

  const handleSplashStart = () => {
    setCurrentStep('production-request');
  };

  const handleProductionRequestNext = () => {
    setCurrentStep('category');
  };

  const handleCategoryNext = () => {
    setCurrentStep('date-range');
  };

  const handleDateRangeNext = () => {
    setCurrentStep('file-selection');
  };

  const handleFileSelectionNext = (selectedFileIDs: number[]) => {
    // Zip creation is handled in FileSelectionStep
    // This could navigate to a success page if needed
  };

  const handleNewSearch = () => {
    navigate('/');
    setQuery({
      productionRequest: null,
      categories: [],
      yearRange: { startYear: null, endYear: null },
      dateRange: null
    });
    setSelectedProductionRequestNumber(null);
  };

  const handleBack = () => {
    const stepOrder: Step[] = ['splash', 'production-request', 'category', 'date-range', 'file-selection'];
    const currentIndex = stepOrder.indexOf(currentStep);
    if (currentIndex > 0) {
      setCurrentStep(stepOrder[currentIndex - 1]);
    }
  };

  const handleQueryTypeSelect = (request: typeof query.productionRequest) => {
    setQuery({ ...query, productionRequest: request });
  };

  const handleDateRangeChange = (dateRange: DateRange) => {
    setQuery({ ...query, dateRange });
  };

  const handleCategoryToggle = (category: DataCategory) => {
    const categories = query.categories.includes(category)
      ? query.categories.filter(c => c !== category)
      : [...query.categories, category];
    setQuery({ ...query, categories });
  };

  return (
    <Routes>
      {/* First Page: Production Request Selection */}
      <Route
        path="/"
        element={
          <ProductionRequestSelectionPage
            onStart={handleProductionRequestStart}
          />
        }
      />

      {/* Main Application Flow */}
      <Route
        path="/app"
        element={
          <div className="min-h-screen bg-gray-50">
            <div className="container mx-auto py-8">
              {currentStep === 'splash' && (
                <SplashStep onNext={handleSplashStart} />
              )}

              {currentStep === 'production-request' && (
                <ProductionRequestStep
                  requests={PRODUCTION_REQUESTS}
                  selected={query.productionRequest}
                  onSelect={handleQueryTypeSelect}
                  onNext={handleProductionRequestNext}
                  onBack={handleBack}
                />
              )}

              {currentStep === 'category' && (
                <CategoryStep
                  categories={['Email', 'Claims', 'Other'] as DataCategory[]}
                  selected={query.categories}
                  onToggle={handleCategoryToggle}
                  onNext={handleCategoryNext}
                  onBack={handleBack}
                />
              )}

              {currentStep === 'date-range' && (
                <DateRangeStep
                  dateRange={query.dateRange || { start: null, end: null }}
                  onDateRangeChange={handleDateRangeChange}
                  onNext={handleDateRangeNext}
                  onBack={handleBack}
                />
              )}

              {currentStep === 'file-selection' && (
                <FileSelectionStep
                  query={query}
                  onNext={handleFileSelectionNext}
                  onBack={handleBack}
                />
              )}
            </div>
          </div>
        }
      />
    </Routes>
  );
}

function App() {
  return (
    <BrowserRouter>
      <AppContent />
    </BrowserRouter>
  );
}

export default App;

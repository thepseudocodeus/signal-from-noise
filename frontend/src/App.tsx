import { useState, useEffect } from 'react';
import { BrowserRouter, Routes, Route, useNavigate } from 'react-router-dom';
import { QueryState, PRODUCTION_REQUESTS, YearRange } from './types/query';
import { ProductionRequestSelectionPage } from './components/pages/ProductionRequestSelectionPage';
import { SplashStep } from './components/steps/SplashStep';
import { ProductionRequestStep } from './components/steps/ProductionRequestStep';
import { YearRangeStep } from './components/steps/YearRangeStep';
import { ResultsStep } from './components/steps/ResultsStep';

type Step = 'splash' | 'production-request' | 'year-range' | 'results';

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
    }
  });

  const [fileCount, setFileCount] = useState<number>(0);
  const [fileSize, setFileSize] = useState<string>('0 bytes');

  useEffect(() => {
    if (currentStep === 'results') {
      setFileCount(1234);
      setFileSize('45.2 MB');
    }
  }, [currentStep]);

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
    setCurrentStep('year-range');
  };

  const handleYearRangeNext = () => {
    setCurrentStep('results');
  };

  const handleNewSearch = () => {
    navigate('/');
    setQuery({
      productionRequest: null,
      categories: [],
      yearRange: { startYear: null, endYear: null }
    });
    setSelectedProductionRequestNumber(null);
  };

  const handleBack = () => {
    const stepOrder: Step[] = ['splash', 'production-request', 'year-range', 'results'];
    const currentIndex = stepOrder.indexOf(currentStep);
    if (currentIndex > 0) {
      setCurrentStep(stepOrder[currentIndex - 1]);
    }
  };

  const handleQueryTypeSelect = (request: typeof query.productionRequest) => {
    setQuery({ ...query, productionRequest: request });
  };

  const handleDateRangeChange = (yearRange: YearRange) => {
    setQuery({ ...query, yearRange });
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

              {currentStep === 'year-range' && (
                <YearRangeStep
                  yearRange={query.yearRange}
                  onYearRangeChange={handleDateRangeChange}
                  onNext={handleYearRangeNext}
                  onBack={handleBack}
                />
              )}

              {currentStep === 'results' && (
                <ResultsStep
                  query={query}
                  fileCount={fileCount}
                  fileSize={fileSize}
                  onNewSearch={handleNewSearch}
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

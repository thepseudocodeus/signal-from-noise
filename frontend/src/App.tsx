import { useState, useEffect } from 'react';
import { QueryState, PRODUCTION_REQUESTS, YearRange } from './types/query';
import { SplashStep } from './components/steps/SplashStep';
import { ProductionRequestStep } from './components/steps/ProductionRequestStep';
import { YearRangeStep } from './components/steps/YearRangeStep';
import { ResultsStep } from './components/steps/ResultsStep';
import { GetEmailFileCount } from '../wailsjs/go/main/App';

type Step = 'splash' | 'production-request' | 'year-range' | 'results';

function App() {
  const [currentStep, setCurrentStep] = useState<Step>('splash');
  const [query, setQuery] = useState<QueryState>({
    productionRequest: null,
    categories: [],
    yearRange: {
      startYear: null,
      endYear: null
    }
  });

  // Results state (will come from Wails backend)
  const [fileCount, setFileCount] = useState<number>(0);
  const [fileSize, setFileSize] = useState<string>('0 bytes');

  // TODO: Connect to Wails backend to get real data
  useEffect(() => {
    if (currentStep === 'results') {
      // Simulate query execution
      // In real implementation:
      // const executeQuery = async () => {
      //   const count = await ExecuteQuery(query);
      //   setFileCount(count);
      // };
      // executeQuery();

      // Placeholder for demo
      setFileCount(1234);
      setFileSize('45.2 MB');
    }
  }, [currentStep]);

  const handleSplashStart = () => {
    setCurrentStep('production-request');
  };

  const handleProductionRequestNext = () => {
    setCurrentStep('year-range');
  };

  const handleYearRangeNext = () => {
    setCurrentStep('results');
  };

  const handleBack = () => {
    const stepOrder: Step[] = ['splash', 'production-request', 'year-range', 'results'];
    const currentIndex = stepOrder.indexOf(currentStep);
    if (currentIndex > 0) {
      setCurrentStep(stepOrder[currentIndex - 1]);
    }
  };

  const handleNewSearch = () => {
    setQuery({
      productionRequest: null,
      categories: [],
      yearRange: { startYear: null, endYear: null }
    });
    setCurrentStep('splash');
  };

  const handleProductionRequestSelect = (request: typeof query.productionRequest) => {
    setQuery({ ...query, productionRequest: request });
  };

  const handleYearRangeChange = (yearRange: YearRange) => {
    setQuery({ ...query, yearRange });
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto py-8">
        {currentStep === 'splash' && (
          <SplashStep onStart={handleSplashStart} />
        )}

        {currentStep === 'production-request' && (
          <ProductionRequestStep
            requests={PRODUCTION_REQUESTS}
            selected={query.productionRequest}
            onSelect={handleProductionRequestSelect}
            onNext={handleProductionRequestNext}
            onBack={handleBack}
          />
        )}

        {currentStep === 'year-range' && (
          <YearRangeStep
            yearRange={query.yearRange}
            onYearRangeChange={handleYearRangeChange}
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
  );
}

export default App;

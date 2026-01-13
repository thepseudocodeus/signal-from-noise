import { Button, Card } from 'flowbite-react';
import { ProgressIndicator } from '../shared/ProgressIndicator';

interface WelcomeStepProps {
  onNext: () => void;
}

export function WelcomeStep({ onNext }: WelcomeStepProps) {
  return (
    <div className="max-w-2xl mx-auto px-4 py-8">
      <ProgressIndicator current={0} total={5} />

      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Welcome to Signal from Noise
        </h1>
        <p className="text-xl text-gray-600 mb-6">
          Query your email data lake with ease
        </p>
        <p className="text-gray-500">
          We'll guide you through a few simple steps to find exactly what you need
        </p>
      </div>

      <Card className="mb-6">
        <div className="text-center">
          <h3 className="text-lg font-semibold text-gray-900 mb-2">
            What you can do:
          </h3>
          <ul className="text-left text-gray-600 space-y-2 max-w-md mx-auto">
            <li>• Search emails by date range</li>
            <li>• Filter by keywords</li>
            <li>• Organize by categories</li>
            <li>• Export results</li>
          </ul>
        </div>
      </Card>

      <div className="flex justify-center">
        <Button
          size="xl"
          onClick={onNext}
          className="px-8 py-3"
        >
          Get Started
        </Button>
      </div>
    </div>
  );
}

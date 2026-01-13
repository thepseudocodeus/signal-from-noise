import { Button, Card } from 'flowbite-react';
import { useState, useEffect } from 'react';
import { GetDataLakeStatus, GetEmailFileCount } from '../../../wailsjs/go/main/App';

interface SplashStepProps {
  onStart: () => void;
}

export function SplashStep({ onStart }: SplashStepProps) {
  const [status, setStatus] = useState<string>('checking...');
  const [fileCount, setFileCount] = useState<number>(0);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const checkStatus = async () => {
      try {
        setLoading(true);
        const statusResult = await GetDataLakeStatus();
        setStatus(statusResult);

        const count = await GetEmailFileCount();
        setFileCount(count);
      } catch (error) {
        console.error('Error checking data lake status:', error);
        setStatus('unavailable');
        setFileCount(0);
      } finally {
        setLoading(false);
      }
    };

    checkStatus();
  }, []);

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold text-gray-900 mb-4">
          Signal from Noise
        </h1>
        <p className="text-xl text-gray-600">
          Query your data lake with precision
        </p>
      </div>

      {/* Status Card */}
      <Card className="mb-6">
        <div className="flex justify-between items-center">
          <div>
            <h3 className="text-lg font-semibold text-gray-900 mb-2">Data Lake Status</h3>
            {loading ? (
              <p className="text-sm text-gray-600">Loading...</p>
            ) : (
              <>
                <p className="text-sm text-gray-600 mb-1">
                  Status: <span className={`font-medium ${status === 'available' ? 'text-green-600' : 'text-red-600'}`}>
                    {status}
                  </span>
                </p>
                <p className="text-sm text-gray-600">
                  Email Files: <span className="font-medium text-blue-600">{fileCount.toLocaleString()}</span>
                </p>
              </>
            )}
          </div>
        </div>
      </Card>

      {/* Search Capabilities (Greyed Out Preview) */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8 opacity-60">
        <Card className="bg-gray-50">
          <h3 className="font-semibold text-gray-700 mb-2">Production Requests</h3>
          <p className="text-sm text-gray-500">Select from 20 production requests</p>
        </Card>
        <Card className="bg-gray-50">
          <h3 className="font-semibold text-gray-700 mb-2">Data Categories</h3>
          <p className="text-sm text-gray-500">Claims, Email, Other</p>
        </Card>
        <Card className="bg-gray-50">
          <h3 className="font-semibold text-gray-700 mb-2">Year Range</h3>
          <p className="text-sm text-gray-500">Filter by date range</p>
        </Card>
      </div>

      <div className="flex justify-center">
        <Button
          size="xl"
          onClick={onStart}
          className="px-8 py-3"
        >
          Start Search
        </Button>
      </div>
    </div>
  );
}

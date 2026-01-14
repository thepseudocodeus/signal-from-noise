import { useState, useEffect } from 'react';
import { Button, Table, Checkbox, Badge, Spinner, Alert } from 'flowbite-react';
import { ProgressIndicator } from '../shared/ProgressIndicator';
import { QueryState } from '../../types/query';

interface FileInfo {
  id: number;
  path: string;
  directory: string;
  category: string;
  date: string;
  size: number;
  privileged: boolean;
  duplicate_hash: string;
  file_name: string;
}

interface FileSearchResult {
  files: FileInfo[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

interface FileSelectionStepProps {
  query: QueryState;
  onNext: (selectedFileIDs: number[]) => void;
  onBack: () => void;
}

export function FileSelectionStep({ query, onNext, onBack }: FileSelectionStepProps) {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [selectedFileIDs, setSelectedFileIDs] = useState<Set<number>>(new Set());
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [totalCount, setTotalCount] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [creatingZip, setCreatingZip] = useState(false);
  const [zipResult, setZipResult] = useState<string | null>(null);

  // Load files when component mounts or filters change
  useEffect(() => {
    loadFiles();
  }, [currentPage, query]);

  const loadFiles = async () => {
    setLoading(true);
    setError(null);

    try {
      // Import Wails bindings
      const { SearchFiles } = await import('../../../wailsjs/go/main/App');

      // Prepare date range
      const dateStart = query.dateRange?.start?.toISOString() || '';
      const dateEnd = query.dateRange?.end?.toISOString() || '';

      // Prepare categories
      const categories = query.categories || [];

      // Call backend
      const result = await SearchFiles({
        production_request_id: query.productionRequest?.id || '',
        date_start: dateStart,
        date_end: dateEnd,
        categories: categories,
        exclude_privileged: true, // Always exclude privileged by default
        page: currentPage,
        page_size: 50,
      });

      setFiles(result.files);
      setTotalPages(result.total_pages);
      setTotalCount(result.total_count);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load files');
      console.error('Error loading files:', err);
    } finally {
      setLoading(false);
    }
  };

  const toggleFileSelection = (fileId: number) => {
    const newSelected = new Set(selectedFileIDs);
    if (newSelected.has(fileId)) {
      newSelected.delete(fileId);
    } else {
      newSelected.add(fileId);
    }
    setSelectedFileIDs(newSelected);
  };

  const selectAll = () => {
    const allIDs = new Set(files.map(f => f.id));
    setSelectedFileIDs(allIDs);
  };

  const deselectAll = () => {
    setSelectedFileIDs(new Set());
  };

  const removeFile = (fileId: number) => {
    const newSelected = new Set(selectedFileIDs);
    newSelected.delete(fileId);
    setSelectedFileIDs(newSelected);
  };

  const formatFileSize = (bytes: number): string => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
    if (bytes < 1024 * 1024 * 1024) return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
    return (bytes / (1024 * 1024 * 1024)).toFixed(2) + ' GB';
  };

  const formatDate = (dateString: string): string => {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString();
    } catch {
      return dateString;
    }
  };

  const handleCreateZip = async () => {
    if (selectedFileIDs.size === 0) {
      setError('Please select at least one file');
      return;
    }

    setCreatingZip(true);
    setError(null);

    try {
      const { CreateZip } = await import('../../../wailsjs/go/main/App');

      const result = await CreateZip({
        production_request_id: query.productionRequest?.id || '',
        file_ids: Array.from(selectedFileIDs),
      });

      if (result.success) {
        setZipResult(result.zip_path);
      } else {
        setError(result.message || 'Failed to create zip file');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create zip file');
      console.error('Error creating zip:', err);
    } finally {
      setCreatingZip(false);
    }
  };

  const totalSelectedSize = files
    .filter(f => selectedFileIDs.has(f.id))
    .reduce((sum, f) => sum + f.size, 0);

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <ProgressIndicator current={4} total={5} />

      <div className="mb-8">
        <h2 className="text-3xl font-bold text-gray-900 mb-3">
          Select Files
        </h2>
        <p className="text-gray-600">
          Review and select files to include in your production request
        </p>
      </div>

      {error && (
        <Alert color="failure" className="mb-6">
          {error}
        </Alert>
      )}

      {zipResult && (
        <Alert color="success" className="mb-6">
          <div>
            <p className="font-semibold">Zip file created successfully!</p>
            <p className="text-sm mt-1">Location: {zipResult}</p>
          </div>
        </Alert>
      )}

      {/* Selection Summary */}
      <div className="mb-6 p-4 bg-gray-50 rounded-lg">
        <div className="flex justify-between items-center">
          <div>
            <span className="text-gray-700 font-medium">
              {selectedFileIDs.size} of {totalCount} files selected
            </span>
            {selectedFileIDs.size > 0 && (
              <span className="ml-4 text-gray-600">
                Total size: {formatFileSize(totalSelectedSize)}
              </span>
            )}
          </div>
          <div className="space-x-2">
            <Button size="sm" color="gray" onClick={selectAll}>
              Select All
            </Button>
            <Button size="sm" color="gray" onClick={deselectAll}>
              Deselect All
            </Button>
          </div>
        </div>
      </div>

      {/* Files Table */}
      {loading ? (
        <div className="flex justify-center items-center py-12">
          <Spinner size="xl" />
          <span className="ml-4 text-gray-600">Loading files...</span>
        </div>
      ) : files.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-600">No files found matching your criteria.</p>
        </div>
      ) : (
        <>
          <div className="overflow-x-auto mb-6">
            <Table hoverable>
              <Table.Head>
                <Table.HeadCell className="w-12">
                  <Checkbox
                    checked={files.length > 0 && files.every(f => selectedFileIDs.has(f.id))}
                    onChange={(e) => {
                      if (e.target.checked) {
                        selectAll();
                      } else {
                        deselectAll();
                      }
                    }}
                  />
                </Table.HeadCell>
                <Table.HeadCell>File Name</Table.HeadCell>
                <Table.HeadCell>Directory</Table.HeadCell>
                <Table.HeadCell>Category</Table.HeadCell>
                <Table.HeadCell>Date</Table.HeadCell>
                <Table.HeadCell>Size</Table.HeadCell>
                <Table.HeadCell>Actions</Table.HeadCell>
              </Table.Head>
              <Table.Body className="divide-y">
                {files.map((file) => (
                  <Table.Row key={file.id} className="bg-white">
                    <Table.Cell>
                      <Checkbox
                        checked={selectedFileIDs.has(file.id)}
                        onChange={() => toggleFileSelection(file.id)}
                      />
                    </Table.Cell>
                    <Table.Cell className="font-medium text-gray-900">
                      {file.file_name}
                    </Table.Cell>
                    <Table.Cell className="text-gray-600">
                      {file.directory}
                    </Table.Cell>
                    <Table.Cell>
                      <Badge color={
                        file.category === 'email' ? 'blue' :
                        file.category === 'claim' ? 'green' : 'gray'
                      }>
                        {file.category}
                      </Badge>
                    </Table.Cell>
                    <Table.Cell className="text-gray-600">
                      {formatDate(file.date)}
                    </Table.Cell>
                    <Table.Cell className="text-gray-600">
                      {formatFileSize(file.size)}
                    </Table.Cell>
                    <Table.Cell>
                      {selectedFileIDs.has(file.id) && (
                        <Button
                          size="xs"
                          color="failure"
                          onClick={() => removeFile(file.id)}
                        >
                          Remove
                        </Button>
                      )}
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          </div>

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-between items-center mb-6">
              <Button
                color="gray"
                disabled={currentPage === 1}
                onClick={() => setCurrentPage(p => Math.max(1, p - 1))}
              >
                Previous
              </Button>
              <span className="text-gray-600">
                Page {currentPage} of {totalPages}
              </span>
              <Button
                color="gray"
                disabled={currentPage >= totalPages}
                onClick={() => setCurrentPage(p => Math.min(totalPages, p + 1))}
              >
                Next
              </Button>
            </div>
          )}
        </>
      )}

      {/* Actions */}
      <div className="flex justify-between">
        <Button color="gray" onClick={onBack}>
          Back
        </Button>
        <div className="space-x-2">
          <Button
            onClick={handleCreateZip}
            disabled={selectedFileIDs.size === 0 || creatingZip}
            className="px-6"
          >
            {creatingZip ? (
              <>
                <Spinner size="sm" className="mr-2" />
                Creating Zip...
              </>
            ) : (
              'Create Zip File'
            )}
          </Button>
        </div>
      </div>
    </div>
  );
}

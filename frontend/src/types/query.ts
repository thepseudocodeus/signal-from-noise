export interface ProductionRequest {
  id: string;
  title: string;
  description: string;
}

export type DataCategory = "Claims" | "Email" | "Other";

export interface YearRange {
  startYear: number | null;
  endYear: number | null;
}

export interface QueryState {
  productionRequest: ProductionRequest | null;
  categories: DataCategory[];
  yearRange: YearRange;
}

// Demo production requests (will expand to 20)
export const PRODUCTION_REQUESTS: ProductionRequest[] = [
  {
    id: "PR-001",
    title: "Production Request #1",
    description: "Email communications timeline"
  },
  {
    id: "PR-002",
    title: "Production Request #2",
    description: "Claims-related documentation"
  },
  {
    id: "PR-003",
    title: "PR-003",
    description: "General document search"
  },
  {
    id: "explore",
    title: "Explore",
    description: "Conduct your own custom search"
  }
];

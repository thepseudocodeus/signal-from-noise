export interface ProductionRequest {
  id: string;
  title: string;
  description: string;
  estimatedTime: number;
}

// [ ] TODO: Make category the proper type to fulfill consumer contract
export interface Category {
  id: string;
  title: string;
  description: string;
}


export type DataCategory = "Claims" | "Email" | "Other";

export interface YearRange {
  startYear: number | null;
  endYear: number | null;
}

export interface DateRange {
  start: Date | null;
  end: Date | null;
}

export interface QueryState {
  productionRequest: ProductionRequest | null;
  categories: DataCategory[];
  yearRange: YearRange;
}

// Demo production requests (will expand to 20)
// [ ] TODO: Logic to estimate time of request - believe this is an ML feature inferred from usage data
export const PRODUCTION_REQUESTS: ProductionRequest[] = [
  {
    id: "PR-001",
    title: "Production Request #1",
    description: "[Description 1]",
    estimatedTime: 5,
  },
  {
    id: "PR-002",
    title: "Production Request #2",
    description: "[Description 2]",
    estimatedTime: 10,
  },
  {
    id: "PR-003",
    title: "Production Request #3",
    description: "[Description 3]",
    estimatedTime: 15,
  },
  {
    id: "PR-004",
    title: "Production Request #4",
    description: "[Description 4]",
    estimatedTime: 20,
  },
  {
    id: "PR-005",
    title: "Production Request #5",
    description: "[Description 5]",
    estimatedTime: 5,
  },
  {
    id: "PR-006",
    title: "Production Request #6",
    description: "[Description 6]",
    estimatedTime: 10,
  },
  {
    id: "PR-007",
    title: "Production Request #7",
    description: "[Description 7]",
    estimatedTime: 15,
  },
  {
    id: "PR-008",
    title: "Production Request #8",
    description: "[Description 8]",
    estimatedTime: 20,
  },
  {
    id: "PR-009",
    title: "Production Request #9",
    description: "[Description 9]",
    estimatedTime: 5,
  },
  {
    id: "PR-010",
    title: "Production Request #10",
    description: "[Description 10]",
    estimatedTime: 10,
  },
  {
    id: "PR-011",
    title: "Production Request #11",
    description: "[Description 11]",
    estimatedTime: 15,
  },
  {
    id: "PR-012",
    title: "Production Request #12",
    description: "[Description 12]",
    estimatedTime: 20,
  },
  {
    id: "PR-013",
    title: "Production Request #13",
    description: "[Description 13]",
    estimatedTime: 5,
  },
  {
    id: "PR-014",
    title: "Production Request #14",
    description: "[Description 14]",
    estimatedTime: 10,
  },
  {
    id: "PR-015",
    title: "Production Request #15",
    description: "[Description 15]",
    estimatedTime: 15,
  },
  {
    id: "PR-016",
    title: "Production Request #16",
    description: "[Description 16]",
    estimatedTime: 20,
  },
  {
    id: "PR-017",
    title: "Production Request #17",
    description: "[Description 17]",
    estimatedTime: 5,
  },
  {
    id: "PR-018",
    title: "Production Request #18",
    description: "[Description 18]",
    estimatedTime: 10,
  },
  {
    id: "PR-019",
    title: "Production Request #19",
    description: "[Description 19]",
    estimatedTime: 15,
  },
  {
    id: "PR-020",
    title: "Production Request #20",
    description: "[Description 20]",
    estimatedTime: 20,
  },
];

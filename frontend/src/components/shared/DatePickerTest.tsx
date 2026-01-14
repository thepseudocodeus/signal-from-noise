/**
 * DatePickerTest - Minimal test to verify Flowbite Datepicker works
 * 
 * Expert Approach: Test with absolute minimum to verify assumption
 * If this works, the issue is in DateRangeSelector
 * If this doesn't work, the issue is in setup/configuration
 */

import { Datepicker } from "flowbite-react";
import { useState } from "react";

export function DatePickerTest() {
  const [date, setDate] = useState<Date | null>(null);

  return (
    <div className="p-4">
      <h3 className="mb-4 text-lg font-bold">DatePicker Test (Minimal)</h3>
      <Datepicker
        value={date}
        onChange={(d) => setDate(d)}
      />
      {date && <p className="mt-2">Selected: {date.toLocaleDateString()}</p>}
    </div>
  );
}

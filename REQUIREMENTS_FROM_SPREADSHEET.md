# Requirements from SearchEngineIdeas.xlsx
## Comprehensive Search Process

**Date:** 2026-01-13
**Source:** User brainstorming spreadsheet

---

## Process Flow (Steps)

### Step 1: Splash/Dashboard Page
- **Action:** Show all search components on the page
- **Requirements:** Minimal dashboard with sections greyed out but visible
- **Constraints:** Cannot select anything other than where to start
- **Note:** Overview of all search capabilities before starting

### Step 2: Trigger a Search
- **Action:** Display instructions
- **Requirements:** Great UX-UI
- **Constraints:** Simple, intuitive

### Step 3: Show Initial State
- **Action:** Show no files and 0 size
- **Note:** Display empty state

### Steps 4-5: Backend Preparation
- **Step 4:** Prepare Julia to handle data (start Julia processor)
- **Step 5:** Prepare Polars to handle data (start Python Polars lazy dataframe)
- **Requirements:** Ensure interop and dependencies loaded for UV in environment
- **Constraints:** Interop between Julia and Python must be seamless, must be quick and not overload memory/CPU/drive

### Step 6: Start Step-by-Step Sequence
- **Action:** Show user simple instructions
- **Requirements:** Like Typeform - display what to do right now only
- **Note:** Begin interactive step-by-step process

### Step 7: Set Search Goal
- **Action:** Pick a production request
- **Requirements:** Select one of production request items from dropdown menu
- **Constraints:** Can only select 1 item; "Explore" selection allows to conduct own search
- **Options:** Production Request #1 through #20, or "Explore"

### Step 8: Set Default Exclusions
- **Action:** Remove common items not needed
- **Default:** Privilege documents (attorney-client-privilege)
- **Note:** Automatic exclusion, can be customized

### Step 9: Select Data Category
- **Action:** Remove not selected categories
- **Options:** Claims, Email, Other (can select any number, only 3 total)
- **Note:** Multiple selection allowed

### Step 10: Display Progress
- **Action:** Display number of files and size so far
- **Note:** Update after each filtering step

### Step 11: Start Year
- **Action:** Ask user for start year
- **Note:** Remove any data before year

### Step 12: End Year
- **Action:** Ask user for end year
- **Note:** Remove any data after year

### Step 13: Update Display
- **Action:** Update display of files and size so far
- **Note:** Show progress after date filtering

### Step 14: Select Relevant Claim (if applicable)
- **Action:** Ask user to select relevant claim
- **Constraints:** Skip if claims not selected in Step 9
- **Note:** Dynamic list from data lake

### Step 15: Select Data Kind
- **Action:** Ask user to select data kind
- **Options:** Application, Archive, Document, Image, Executable, Folder, Video, Audio, Presentation, Text, Other

### Step 16: Update Display
- **Action:** Update display of files and size so far
- **Note:** Show progress after data kind filtering

### Step 17: Select Relevant EEO Area
- **Action:** Ask user to select relevant EEO area
- **Options:** Protected class areas, EEO laws, relevant arguments

---

## Selection Options (From Spreadsheet)

### Production Requests (Step 7)
- Production Request #1 through #20
- Explore (allows custom search)

### Data Categories (Step 9)
- Claim
- Email
- Other

### Protected Class Areas (Step 17)
- Race & Color
- Religion
- Sex/Gender
- National Origin
- Age (40+)
- Disability
- Genetic Information
- Veteran Status
- Other Class
- None

### EEO Laws (Step 17)
- Title VII of the Civil Rights Act of 1964
- Age Discrimination in Employment Act (ADEA)
- Americans with Disabilities Act (ADA)
- Equal Pay Act (EPA)
- Genetic Information Nondiscrimination Act (GINA)
- US Telework Enhancement Act of 2010
- Other
- None

### Relevant Arguments (Step 17)
- Discrimination
- Adverse Action
- Harassment
- Reasonable Accommodation
- Retaliation
- Disparate Treatment
- Disparate Impact
- Affirmative Action
- EEO Statement/Tagline

### Claims (Step 14 - Dynamic)
- None
- [Insert Information From Data Lake On Claims]
- Make Default: TRUE/FALSE/None

### Data Kinds (Step 15)
- Application
- Archive
- Document
- Image
- Executable
- Folder
- Video
- Audio
- Presentation
- Text
- Other Kind

### Additional Filters (Future/Extended)
- File Types (Excel, Word, PDF, PNG, JPG, ZIP, PowerPoint, etc.)
- File Sizes (< 100 KB through >= 10 GB)
- Time of Day (0:00:00 through 23:45:00)
- Quarters (Q1, Q2, Q3, Q4)
- Months (January through December)
- Weekdays (Sunday through Saturday)
- Years (2000 through 2026)
- Tags (WithAttorney, WithBoss, WithColleague, Medical, WorkRelated, etc.)
- Sentiment (Positive, Negative, Neutral, Unknown)
- Folder Sizes
- And many more...

---

## Key Differences from Initial Implementation

1. **Splash/Dashboard First:** Need initial overview page showing all capabilities
2. **Production Requests:** 20 specific requests + "Explore" option
3. **Year Selection:** Years instead of full date pickers (simpler for large datasets)
4. **Conditional Steps:** Claims step only if Claims category selected
5. **Progress Updates:** Show file count and size after each filtering step
6. **More Complex Filtering:** EEO areas, laws, arguments, data kinds, etc.
7. **Default Exclusions:** Automatic privilege document exclusion

---

## Implementation Priority (For Demo)

### Must Have (Demo):
1. Splash/Dashboard page
2. Production Request selection (20 requests + Explore)
3. Data Category selection
4. Year range selection (start/end year)
5. Progress display (file count/size)
6. Basic filtering flow

### Nice to Have (Demo):
7. Data Kind selection
8. EEO Area selection
9. Claims selection (if Claims category selected)
10. Default exclusions

### Future (Post-Demo):
11. All advanced filters (file types, sizes, times, tags, sentiment, etc.)
12. Julia/Python backend integration
13. Full data lake processing

---

## UI/UX Requirements

- **Typeform-Inspired:** One question at a time, clear instructions
- **Progress Tracking:** Show file count and size updates
- **Conditional Logic:** Skip steps based on previous selections
- **Simple First:** Start with year selection (not full date pickers)
- **Dashboard Overview:** Initial splash showing all capabilities

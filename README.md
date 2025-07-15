# Resume Analyzer

A powerful CLI tool that processes PDF resumes using AWS Textract for OCR and AWS Bedrock for intelligent summarization and analysis.

## Features

- **PDF Processing**: Converts multi-page PDFs to text using AWS Textract
- **Intelligent Summarization**: Uses AWS Bedrock (Claude) to extract key information
- **Consolidated Analysis**: Creates CSV files with applicant comparisons
- **Technical Skills Assessment**: Specifically identifies key technical skills
- **Seniority Evaluation**: Assesses experience levels based on multiple factors

## Prerequisites

- Go 1.16 or higher (for local development)
- Docker (for containerized usage)
- AWS CLI configured with appropriate credentials
- AWS Textract and Bedrock access

## Installation

### Option 1: Local Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd resume-analyzer
   ```

2. **Install dependencies:**
   ```bash
   make deps
   ```

3. **Build the application:**
   ```bash
   make build
   ```

### Option 2: Docker Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd resume-analyzer
   ```

2. **Build the Docker image:**
   ```bash
   docker build -t resume-analyzer .
   ```

   Or use Docker Compose:
   ```bash
   docker-compose build
   ```

## Configuration

Create a `.resume-analyzer.yaml` file in your home directory or project root:

```yaml
bedrock_model_id: anthropic.claude-3-5-sonnet-20240620-v1:0
anthropic_version: bedrock-2023-05-31
```

## Usage

### Quick Start (Complete Workflow)

#### Local Usage
Run the entire pipeline with one command:

```bash
make all-steps
```

#### Docker Usage
```bash
# Using Docker directly
docker run --rm -v $(pwd)/input_pdfs:/app/input_pdfs:ro \
  -v $(pwd)/output_txts:/app/output_txts \
  -v $(pwd)/output_summaries:/app/output_summaries \
  -v $(pwd)/output_consolidated:/app/output_consolidated \
  resume-analyzer convert-pdfs -i input_pdfs -o output_txts && \
docker run --rm -v $(pwd)/output_txts:/app/output_txts:ro \
  -v $(pwd)/output_summaries:/app/output_summaries \
  resume-analyzer summarize -i output_txts -o output_summaries && \
docker run --rm -v $(pwd)/output_summaries:/app/output_summaries:ro \
  -v $(pwd)/output_consolidated:/app/output_consolidated \
  resume-analyzer consolidate -i output_summaries -o output_consolidated/consolidated_table_$(date +%Y%m%d_%H%M%S).csv

# Using Docker Compose (easier)
docker-compose run --rm resume-analyzer convert-pdfs -i input_pdfs -o output_txts && \
docker-compose run --rm resume-analyzer summarize -i output_txts -o output_summaries && \
docker-compose run --rm resume-analyzer consolidate -i output_summaries -o output_consolidated/consolidated_table_$(date +%Y%m%d_%H%M%S).csv
```

This will:
1. Convert PDFs to text (`input_pdfs` → `output_txts`)
2. Generate summaries (`output_txts` → `output_summaries`)
3. Create consolidated CSV (`output_summaries` → `consolidated_table_YYYYMMDD_HHMMSS.csv`)

### Organizing Multiple Resume Sets with Subfolders

For better organization when processing multiple sets of resumes (e.g., different job positions, departments, or recruitment campaigns), use the subfolder feature:

#### Recommended Directory Structure
```
resume-analyzer/
├── input_pdfs/
│   ├── engineering/          # Engineering positions
│   │   ├── senior-dev.pdf
│   │   ├── frontend-dev.pdf
│   │   └── backend-dev.pdf
│   ├── marketing/            # Marketing positions
│   │   ├── content-writer.pdf
│   │   └── social-media.pdf
│   └── sales/                # Sales positions
│       ├── account-exec.pdf
│       └── sales-manager.pdf
├── output_txts/
│   ├── engineering/          # Auto-created when using SUBFOLDER
│   ├── marketing/
│   └── sales/
├── output_summaries/
│   ├── engineering/          # Auto-created when using SUBFOLDER
│   ├── marketing/
│   └── sales/
└── output_consolidated/
    ├── engineering/          # Auto-created when using SUBFOLDER
    ├── marketing/
    └── sales/
```

#### Processing Specific Resume Sets

```bash
# Process engineering resumes
make SUBFOLDER=engineering all-steps

# Process marketing resumes
make SUBFOLDER=marketing all-steps

# Process sales resumes
make SUBFOLDER=sales all-steps
```

This approach provides:
- **Organized workflow**: Keep different resume sets separate
- **Parallel processing**: Process multiple sets simultaneously
- **Easy comparison**: Compare results across different positions
- **Clean outputs**: Each subfolder gets its own output directories

### Individual Commands

#### 1. Convert PDFs to Text

**Local:**
```bash
# Process all PDFs in input_pdfs/
make convert-pdfs

# Process PDFs in a specific subfolder
make SUBFOLDER=engineering convert-pdfs

# Direct command
./bin/resume-analyzer convert-pdfs -i input_pdfs -o output_txts
```

**Docker:**
```bash
# Using Docker directly
docker run --rm -v $(pwd)/input_pdfs:/app/input_pdfs:ro \
  -v $(pwd)/output_txts:/app/output_txts \
  resume-analyzer convert-pdfs -i input_pdfs -o output_txts

# Using Docker Compose
docker-compose run --rm resume-analyzer convert-pdfs -i input_pdfs -o output_txts
```

#### 2. Generate Summaries

**Local:**
```bash
# Generate summaries for all processed texts
make summarize

# Generate summaries for a specific subfolder
make SUBFOLDER=engineering summarize

# Direct command
./bin/resume-analyzer summarize -i output_txts -o output_summaries
```

**Docker:**
```bash
# Using Docker directly
docker run --rm -v $(pwd)/output_txts:/app/output_txts:ro \
  -v $(pwd)/output_summaries:/app/output_summaries \
  resume-analyzer summarize -i output_txts -o output_summaries

# Using Docker Compose
docker-compose run --rm resume-analyzer summarize -i output_txts -o output_summaries
```

#### 3. Create Consolidated CSV

**Local:**
```bash
# Create consolidated CSV for all summaries
make consolidate

# Create consolidated CSV for a specific subfolder
make SUBFOLDER=engineering consolidate

# Direct command
./bin/resume-analyzer consolidate -i output_summaries -o consolidated_table.csv
```

**Docker:**
```bash
# Using Docker directly
docker run --rm -v $(pwd)/output_summaries:/app/output_summaries:ro \
  -v $(pwd)/output_consolidated:/app/output_consolidated \
  resume-analyzer consolidate -i output_summaries -o output_consolidated/consolidated_table.csv

# Using Docker Compose
docker-compose run --rm resume-analyzer consolidate -i output_summaries -o output_consolidated/consolidated_table.csv
```

#### 4. Query All Resumes

**Local:**
```bash
# Query all processed resumes
./bin/resume-analyzer query -p "Who has the most experience with Python?" -i output_txts

# Query resumes from a specific subfolder
./bin/resume-analyzer query -p "Who has the most experience with Python?" -i output_txts/engineering

# Query with output to file (optional)
./bin/resume-analyzer query -p "Compare the technical skills of all candidates" -i output_txts -o query_response.txt
```

**Docker:**
```bash
# Using Docker directly
docker run --rm -v $(pwd)/output_txts:/app/output_txts:ro \
  -v $(pwd):/app/query_output \
  resume-analyzer query -p "Who has the most experience with Python?" -i output_txts -o query_output/query_response.txt

# Using Docker Compose
docker-compose run --rm resume-analyzer query -p "Who has the most experience with Python?" -i output_txts -o query_response.txt
```

### Directory Structure

#### Basic Structure
```
resume-analyzer/
├── input_pdfs/           # Place your PDF resumes here
├── output_txts/          # Extracted text files
├── output_summaries/     # AI-generated summaries
├── output_consolidated/  # Consolidated CSV files
├── query_response.txt    # Query responses (optional)
└── bin/                  # Built executable
```

#### Advanced Structure with Subfolders
```
resume-analyzer/
├── input_pdfs/
│   ├── engineering/      # Engineering positions
│   ├── marketing/        # Marketing positions
│   └── sales/           # Sales positions
├── output_txts/
│   ├── engineering/      # Auto-created when using SUBFOLDER
│   ├── marketing/
│   └── sales/
├── output_summaries/
│   ├── engineering/      # Auto-created when using SUBFOLDER
│   ├── marketing/
│   └── sales/
├── output_consolidated/
│   ├── engineering/      # Auto-created when using SUBFOLDER
│   ├── marketing/
│   └── sales/
└── bin/
```

## Technical Skills Detected

The tool specifically identifies these technical skills:

- **Frontend**: TypeScript, JavaScript, React, Vue, Angular, Next.js
- **Backend**: Python, Golang
- **AI/ML**: AI, LLM, Machine Learning
- **Cloud**: AWS, GCP, Azure, Alibaba Cloud
- **DevOps**: Terraform, CI/CD, Docker, Kubernetes

## Output Format

### Summary Files
Each resume gets a detailed summary with:
- Name, current role, company
- Years of experience
- Seniority assessment (Junior/Mid/Senior/Lead/Manager/Director/VP/C-Level)
- Technical skills breakdown
- Status and achievements

### Consolidated CSV
A CSV file with columns:
Applicant,Role,Seniority,Status,Current Position,Current Company,Years of Exp,CV Link,Skillset,Remarks

The CSV format makes it easy to:
- Import into spreadsheet applications (Excel, Google Sheets)
- Process with data analysis tools
- Filter and sort applicant data
- Generate reports and visualizations

### Query Responses
The query command allows you to ask custom questions about all resumes at once. Examples:

- **Skill Comparison**: "Who has the most experience with React and TypeScript?"
- **Experience Analysis**: "Which candidates have more than 5 years of experience?"
- **Role Matching**: "Find candidates suitable for a Senior Backend Developer role"
- **Company Analysis**: "Which candidates have worked at FAANG companies?"
- **Technical Assessment**: "Compare the cloud computing skills of all candidates"

The query combines all resume texts into a single prompt, allowing Bedrock to provide comprehensive analysis across all candidates.

## Makefile Commands

### Basic Commands
```bash
make build          # Build the application
make clean          # Clean build artifacts
make clean-outputs  # Clean all output directories
make convert-pdfs   # Convert PDFs to text
make summarize      # Generate summaries
make consolidate    # Create consolidated CSV
make query          # Show query command examples
make all-steps      # Run complete workflow
make help           # Show all available commands
```

### Subfolder Commands
```bash
# Process specific subfolder
make SUBFOLDER=engineering convert-pdfs
make SUBFOLDER=engineering summarize
make SUBFOLDER=engineering consolidate
make SUBFOLDER=engineering all-steps

# Clean specific subfolder outputs
make SUBFOLDER=engineering clean-outputs

# Examples for different departments
make SUBFOLDER=marketing all-steps
make SUBFOLDER=sales all-steps
make SUBFOLDER=hr all-steps
```

## Development

### Running Tests
```bash
make test
make test-coverage
```

### Code Quality
```bash
make fmt    # Format code
make lint   # Lint code (requires golangci-lint)
```

### Cross-Platform Build
```bash
make build-all  # Build for Linux, macOS, Windows
```

## Troubleshooting

### AWS Configuration
- Ensure AWS credentials are properly configured
- Verify access to Textract and Bedrock services
- Check region settings (default: ap-southeast-1)

### Docker Issues
- **"Permission denied"**: Ensure the output directories have proper permissions
  ```bash
  mkdir -p output_txts output_summaries output_consolidated
  chmod 755 output_txts output_summaries output_consolidated
  ```
- **"AWS credentials not found"**: Mount your AWS credentials or use environment variables
  ```bash
  # Mount AWS credentials
  docker run --rm -v ~/.aws:/home/appuser/.aws:ro resume-analyzer --help
  
  # Or use environment variables
  docker run --rm -e AWS_ACCESS_KEY_ID=your_key -e AWS_SECRET_ACCESS_KEY=your_secret resume-analyzer --help
  ```
- **"Volume mount issues"**: Ensure the directories exist before running Docker commands
  ```bash
  mkdir -p input_pdfs output_txts output_summaries output_consolidated
  ```

### PDF Issues
- Ensure PDFs are readable and not corrupted
- Multi-page PDFs are automatically split and processed
- Maximum 5MB per page for Textract processing

### Common Errors
- **"Request has unsupported document format"**: PDF may be corrupted or too large
- **"Failed to load AWS config"**: Check AWS credentials and region
- **"No content in response"**: Bedrock API may be unavailable

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

[Add your license information here] 
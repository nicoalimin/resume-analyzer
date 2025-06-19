# Resume Analyzer

A powerful CLI tool that processes PDF resumes using AWS Textract for OCR and AWS Bedrock for intelligent summarization and analysis.

## Features

- **PDF Processing**: Converts multi-page PDFs to text using AWS Textract
- **Intelligent Summarization**: Uses AWS Bedrock (Claude) to extract key information
- **Consolidated Analysis**: Creates CSV files with applicant comparisons
- **Technical Skills Assessment**: Specifically identifies key technical skills
- **Seniority Evaluation**: Assesses experience levels based on multiple factors

## Prerequisites

- Go 1.16 or higher
- AWS CLI configured with appropriate credentials
- AWS Textract and Bedrock access

## Installation

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

## Configuration

Create a `.resume-analyzer.yaml` file in your home directory or project root:

```yaml
bedrock_model_id: anthropic.claude-3-5-sonnet-20240620-v1:0
anthropic_version: bedrock-2023-05-31
```

## Usage

### Quick Start (Complete Workflow)

Run the entire pipeline with one command:

```bash
make all-steps
```

This will:
1. Convert PDFs to text (`input_pdfs` → `output_txts`)
2. Generate summaries (`output_txts` → `output_summaries`)
3. Create consolidated CSV (`output_summaries` → `consolidated_table_YYYYMMDD_HHMMSS.csv`)

### Individual Commands

#### 1. Convert PDFs to Text
```bash
make convert-pdfs
# or
./bin/resume-analyzer convert-pdfs -i input_pdfs -o output_txts
```

#### 2. Generate Summaries
```bash
make summarize
# or
./bin/resume-analyzer summarize -i output_txts -o output_summaries
```

#### 3. Create Consolidated CSV
```bash
make consolidate
# or
./bin/resume-analyzer consolidate -i output_summaries -o consolidated_table.csv
```

#### 4. Query All Resumes
```bash
# Query with output to console (default)
./bin/resume-analyzer query -p "Who has the most experience with Python?" -i output_txts

# Query with output to file (optional)
./bin/resume-analyzer query -p "Compare the technical skills of all candidates" -i output_txts -o query_response.txt
```

### Directory Structure

```
resume-analyzer/
├── input_pdfs/           # Place your PDF resumes here
├── output_txts/          # Extracted text files
├── output_summaries/     # AI-generated summaries
├── consolidated_table_*.csv  # Final CSV file
├── query_response.txt    # Query responses (optional)
└── bin/                  # Built executable
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
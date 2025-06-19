package prompts

// GetSummaryPrompt returns the enhanced summarization prompt for resume analysis
func GetSummaryPrompt(text string) string {
	return `Please provide a comprehensive summary of the following resume. Focus on extracting key information for recruitment purposes:

**Key Information to Extract:**
1. **Name**: Full name of the applicant
2. **Current Role/Position**: Current job title
3. **Current Company**: Current employer
4. **Years of Experience**: Total years of professional experience
5. **Seniority Level**: Assess as Junior/Mid/Senior/Lead/Manager/Director/VP/C-Level based on:
   - Years of experience
   - Scope of responsibilities
   - Team size managed
   - Technical complexity handled
   - Leadership indicators

**Technical Skills Assessment:**
Please specifically identify and highlight these skills if present:
- **Frontend**: TypeScript, JavaScript, React, Vue, Angular, Next.js
- **Backend**: Python, Golang
- **AI/ML**: AI, LLM, Machine Learning
- **Cloud**: AWS, GCP, Azure, Alibaba Cloud
- **DevOps**: Terraform, CI/CD, Docker, Kubernetes

**Additional Information:**
- **Status**: Active/Passive/Open to opportunities
- **Key Achievements**: Notable accomplishments
- **Education**: Relevant education background
- **Remarks**: Any special notes or observations

**Resume Content:**
` + text + `

Please provide a structured summary that captures all the above information clearly.`
}

// GetExtractionPrompt returns the prompt for extracting structured information from summaries
func GetExtractionPrompt(summary string) string {
	return `Extract the following information from this resume summary and return ONLY a JSON object with these exact keys (use "N/A" if not found):

{
  "name": "Full Name",
  "role": "Job Role/Title",
  "seniority": "Junior/Mid/Senior/Lead/Manager/Director/VP/C-Level",
  "status": "Active/Passive/Open to opportunities",
  "current_position": "Current Job Title",
  "current_company": "Current Company Name",
  "years_of_exp": "X years",
  "cv_link": "N/A",
  "skillset": "Key skills separated by commas",
  "remarks": "Brief notes or observations"
}

Resume Summary:
` + summary + `

JSON:`
}

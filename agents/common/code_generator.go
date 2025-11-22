package common

import (
	"context"
	"fmt"
	"os"

	"github.com/natalie/go-flags-eval/tools"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/genai"
)

// CodeGeneratorConfig configures a code generation agent
type CodeGeneratorConfig struct {
	APIKey      string
	ModelName   string
	Instruction string
	OutputDir   string
}

// CodeGeneratorAgent creates Go code files based on instructions
type CodeGeneratorAgent struct {
	agent   *llmagent.Agent
	context context.Context
}

// NewCodeGeneratorAgent creates a new code generation agent
func NewCodeGeneratorAgent(ctx context.Context, cfg CodeGeneratorConfig) (*CodeGeneratorAgent, error) {
	// Create model
	model, err := gemini.NewModel(ctx, cfg.ModelName, &genai.ClientConfig{
		APIKey: cfg.APIKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create model: %w", err)
	}

	// Create tools
	writeFileTool, err := tools.NewFileWriteTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create write file tool: %w", err)
	}

	readFileTool, err := tools.NewFileReadTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create read file tool: %w", err)
	}

	listFilesTool, err := tools.NewListFilesTool()
	if err != nil {
		return nil, fmt.Errorf("failed to create list files tool: %w", err)
	}

	// Create agent
	agentInstance, err := llmagent.New(llmagent.Config{
		Name:        "code_generator",
		Model:       model,
		Description: "Generates Go code based on instructions",
		Instructions: fmt.Sprintf(`You are a code generation assistant. %s
Output directory: %s

When generating code:
1. Create well-structured, idiomatic Go code
2. Include proper package declarations
3. Add comments for exported functions
4. Use write_file to save code to files
5. Follow Go best practices`, cfg.Instruction, cfg.OutputDir),
		Tools: []tool.Tool{writeFileTool, readFileTool, listFilesTool},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output dir: %w", err)
	}

	return &CodeGeneratorAgent{
		agent:   agentInstance,
		context: ctx,
	}, nil
}

// Generate runs the code generation agent with a specific task
func (a *CodeGeneratorAgent) Generate(task string) error {
	// TODO: Implement agent execution when ADK API is clearer
	// For now this is a placeholder showing the structure
	_ = task
	return fmt.Errorf("not implemented: agent execution pending ADK API clarification")
}

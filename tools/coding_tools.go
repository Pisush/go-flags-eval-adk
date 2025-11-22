package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// File reading tool
type FileReadInput struct {
	Path string `json:"path" jsonschema:"description=Path to the file to read"`
}

type FileReadOutput struct {
	Path    string `json:"path"`
	Content string `json:"content"`
	Size    int    `json:"size"`
}

func NewFileReadTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "read_file",
			Description: "Reads the contents of a file from disk",
		},
		func(ctx tool.Context, input FileReadInput) (FileReadOutput, error) {
			content, err := os.ReadFile(input.Path)
			if err != nil {
				return FileReadOutput{}, fmt.Errorf("failed to read file: %w", err)
			}

			return FileReadOutput{
				Path:    input.Path,
				Content: string(content),
				Size:    len(content),
			}, nil
		},
	)
}

// File writing tool
type FileWriteInput struct {
	Path    string `json:"path" jsonschema:"description=Path to the file to write"`
	Content string `json:"content" jsonschema:"description=Content to write to the file"`
}

type FileWriteOutput struct {
	Path         string `json:"path"`
	BytesWritten int    `json:"bytes_written"`
	Success      bool   `json:"success"`
}

func NewFileWriteTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "write_file",
			Description: "Writes content to a file on disk",
		},
		func(ctx tool.Context, input FileWriteInput) (FileWriteOutput, error) {
			err := os.WriteFile(input.Path, []byte(input.Content), 0644)
			if err != nil {
				return FileWriteOutput{}, fmt.Errorf("failed to write file: %w", err)
			}

			return FileWriteOutput{
				Path:         input.Path,
				BytesWritten: len(input.Content),
				Success:      true,
			}, nil
		},
	)
}

// Grep search tool
type GrepInput struct {
	Pattern string `json:"pattern" jsonschema:"description=Pattern to search for"`
	Path    string `json:"path" jsonschema:"description=Path to search in (file or directory)"`
}

type GrepMatch struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Content string `json:"content"`
}

type GrepOutput struct {
	Pattern    string      `json:"pattern"`
	Matches    []GrepMatch `json:"matches"`
	MatchCount int         `json:"match_count"`
}

func NewGrepTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "grep",
			Description: "Search for a pattern in files within a directory or specific file",
		},
		func(ctx tool.Context, input GrepInput) (GrepOutput, error) {
			var matches []GrepMatch

			err := filepath.Walk(input.Path, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return err
				}

				content, err := os.ReadFile(path)
				if err != nil {
					return nil // Skip files we can't read
				}

				lines := strings.Split(string(content), "\n")
				for lineNum, line := range lines {
					if strings.Contains(line, input.Pattern) {
						matches = append(matches, GrepMatch{
							File:    path,
							Line:    lineNum + 1,
							Content: strings.TrimSpace(line),
						})
					}
				}

				return nil
			})

			if err != nil {
				return GrepOutput{}, fmt.Errorf("grep failed: %w", err)
			}

			return GrepOutput{
				Pattern:    input.Pattern,
				Matches:    matches,
				MatchCount: len(matches),
			}, nil
		},
	)
}

// List files tool
type ListFilesInput struct {
	Path string `json:"path" jsonschema:"description=Directory path to list"`
}

type FileEntry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

type ListFilesOutput struct {
	Path  string      `json:"path"`
	Files []FileEntry `json:"files"`
	Count int         `json:"count"`
}

func NewListFilesTool() (tool.Tool, error) {
	return functiontool.New(
		functiontool.Config{
			Name:        "list_files",
			Description: "Lists all files and directories in a given path",
		},
		func(ctx tool.Context, input ListFilesInput) (ListFilesOutput, error) {
			entries, err := os.ReadDir(input.Path)
			if err != nil {
				return ListFilesOutput{}, fmt.Errorf("failed to list files: %w", err)
			}

			var files []FileEntry
			for _, entry := range entries {
				files = append(files, FileEntry{
					Name:  entry.Name(),
					IsDir: entry.IsDir(),
				})
			}

			return ListFilesOutput{
				Path:  input.Path,
				Files: files,
				Count: len(files),
			}, nil
		},
	)
}
